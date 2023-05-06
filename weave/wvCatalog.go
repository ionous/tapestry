package weave

import (
	"database/sql"
	"log"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains         map[string]*Domain
	processing      DomainStack
	pendingDomains  []*Domain
	resolvedDomains cachedTable
	Errors          []error
	cursor          string
	writer          Writer
	run             rt.Runtime
	db              *tables.Cache

	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	macros      macroReg
	autoCounter Counters
	env         Environ
}

func NewCatalog(db *sql.DB) *Catalog {
	qx, e := qdb.NewQueryx(db)
	if e != nil {
		panic(e)
	}
	// fix: this needs cleanup
	// why are use "writer" when db is the only output?
	// -- why does that .Insert statement produce a type

	// qx should be separated from q ( maybe put into Mdl? )
	// initial cursor should be set externally or passed in
	// what should be public for Catalog?
	// or...
	return &Catalog{
		macros:      make(macroReg),
		oneTime:     make(map[string]bool),
		autoCounter: make(Counters),
		cursor:      "x", // fix
		db:          tables.NewCache(db),
		domains:     make(map[string]*Domain),
		writer:      mdl.Writer(ExecWriter(db).Write),
		run: qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			qna.DecodeNone("unsupported decoder"),
			qna.NewOptions()),
	}
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

func (c *Catalog) EndDomain() (err error) {
	if _, ok := c.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else {
		c.processing.Pop()
	}
	return
}

// calls to schedule() between begin/end domain write to this newly declared domain.
func (c *Catalog) BeginDomain(name string, requires []string) (err error) {
	if n, ok := UniformString(name); !ok {
		err = InvalidString(name)
	} else if reqs, e := UniformStrings(requires); e != nil {
		err = e // transform all the names first to determine any errors
	} else {
		d := c.ensureDomain(n, c.cursor)
		if d, ok := c.processing.Top(); ok {
			reqs = append(reqs, d.name) // domains implicitly depend on the current domain
		}
		if e := d.addDependencies(reqs); e != nil {
			err = e
		} else {
			c.processing.Push(d)
		}
	}
	return
}

// find or create the domain
func (c *Catalog) ensureDomain(n, at string) (ret *Domain) {
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{Requires: Requires{name: n, at: at}, catalog: c}
		c.domains[n] = d
		c.pendingDomains = append(c.pendingDomains, d)
		ret = d
	}
	return
}

type conflict struct {
	Category, Domain, At, Key, Value string
}

func (c conflict) Error() string {
	return errutil.Sprint("unexpected conflict in domain %q at %q for %s %q",
		c.Domain, c.At, c.Category, c.Value)
}

func (c *Catalog) findConflicts() (err error) {
	if res, e := findConflicts(c.db); e != nil {
		err = e
	} else {
		for _, e := range res {
			err = errutil.Append(err, e)
		}
	}
	return
}

func findConflicts(db tables.Querier) (ret []conflict, err error) {
	if rows, e := db.Query(
		`select 'plural', a.domain, a.at, a.many, a.one
			from active_plurals as a 
			join active_plurals as b 
			where a.domain != b.domain 
			and a.many == b.many 
			and a.one != b.one
		union all
		select 'opposite', a.domain, a.at, a.oneWord, a.otherWord 
			from active_rev as a 
			join active_rev as b 
			where a.domain != b.domain 
			and a.oneWord == b.oneWord 
			and a.otherWord != b.otherWord
			`); e != nil {
		err = errutil.New("find conflicts", e)
	} else {
		var cat, domain, key, value string
		var at sql.NullString
		if e := tables.ScanAll(rows, func() (_ error) {
			ret = append(ret, conflict{
				Category: cat,
				Domain:   domain,
				At:       at.String,
				Key:      key,
				Value:    value,
			})
			return
		}, &cat, &domain, &at, &key, &value); e != nil {
			err = e
		}
	}
	return
}

func (c *Catalog) postPhase(p assert.Phase, d *Domain) (err error) {
	switch p {
	// PostFields: with the current domain looping pattern
	// we have to hit all locals from all domains before writing
	// and macro is after locals...
	case assert.MacroPhase:
		if deps, e := d.ResolveDomainKinds(); e != nil {
			err = e
		} else {
			for _, dep := range deps {
				kind := dep.Leaf().(*ScopedKind)
				fields := [][]UniformField{
					kind.header.paramList,
					kind.header.resList,
					kind.pendingFields,
				}
				for _, list := range fields {
					for _, field := range list {
						if e := field.assembleField(kind); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
	case assert.NounPhase:
		_, err = d.ResolveDomainNouns()
	}
	return
}

func (c *Catalog) writePhase(p assert.Phase) (err error) {
	w := c.writer
	// switch or map better?
	switch p {
	case assert.FieldPhase:
		err = c.WriteFields(w)
	case assert.NounPhase:
		err = c.WriteNouns(w)
	case assert.ValuePhase:
		err = c.WriteValues(w)
	case assert.RelativePhase:
		if e := c.WriteRelations(w); e != nil {
			err = e
		} else {
			err = c.WritePairs(w)
		}
	case assert.PatternPhase:
		if e := c.WritePatterns(w); e != nil {
			err = e
		} else if e := c.WriteLocals(w); e != nil {
			err = e
		} else if e := c.WriteRules(w); e != nil {
			err = e
		}
	case assert.AliasPhase:
		err = c.WriteNames(w)

	case assert.DirectivePhase:
		// grammar
		err = c.WriteDirectives(w)

	case assert.PostDomain:
		// when to write these?
		err = c.WriteChecks(w)
	}
	return
}

// FIX -- its a goal to remove this function
// in theory the assembly should walk all domains in the proper order
// and assemble them one at a time, so nobody else needs to know the domain order.
//
// work out the hierarchy of all the domains, and return them in a list.
// the list has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
func (c *Catalog) ResolveDomains() (ret []*Domain, err error) {
	if rows, e := c.db.Query(`select domain from mdl_domain order by path`); e != nil {
		err = errutil.New("resolve domains", e)
	} else {
		var name string
		err = tables.ScanAll(rows, func() (err error) {
			if d, ok := c.GetDomain(name); !ok {
				err = errutil.New("no such domain", name)
			} else {
				ret = append(ret, d)
			}
			return
		}, &name)
	}
	return
}

// FIX -- its a goal to remove this function
func (c *Catalog) ResolveKinds() (ret DependencyTable, err error) {
	var out DependencyTable
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, d := range ds {
			if ks, e := d.ResolveDomainKinds(); e != nil {
				err = errutil.Append(err, e)
			} else {
				out = append(out, ks...)
			}
		}
	}
	if err == nil {
		ret = out
	}
	return
}

// FIX -- its a goal to remove this function
func (c *Catalog) ResolveNouns() (ret DependencyTable, err error) {
	// fix? is there anyway to make this more "automatically" resolve domains and kinds?
	var out DependencyTable
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, d := range ds {
			if ns, e := d.ResolveDomainNouns(); e != nil {
				err = e
				break
			} else {
				out = append(out, ns...)
			}
		}
	}
	if err == nil {
		ret = out
	}
	return
}

type partialWriter struct {
	w      Writer
	fields []interface{}
}

func (p *partialWriter) Write(q string, args ...interface{}) error {
	return p.w.Write(q, append(p.fields, args...)...)
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) FindDependency(name string) (ret Dependency, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = d, true
	}
	return
}
