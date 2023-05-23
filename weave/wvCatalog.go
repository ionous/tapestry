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
	writer          mdl.Modeler
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
	m, e := mdl.NewModeler(db)
	if e != nil {
		panic(e)
	}
	// fix: this needs cleanup
	// write should be called modeler
	// initial cursor should be set externally or passed in
	// what should be public for Catalog?
	// no panics on creation... etc.
	return &Catalog{
		macros:      make(macroReg),
		oneTime:     make(map[string]bool),
		autoCounter: make(Counters),
		cursor:      "x", // fix
		db:          tables.NewCache(db),
		domains:     make(map[string]*Domain),
		writer:      m,
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
	} else if d, e := c.addDomain(n, c.cursor, reqs...); e != nil {
		err = e
	} else {
		c.processing.Push(d)
	}
	return
}

func (c *Catalog) ensureDomain(n string) (ret *Domain) {
	// find or create the domain
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n, catalog: c}
		c.pendingDomains = append(c.pendingDomains, d)
		c.domains[n] = d
		ret = d
	}
	return
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) addDomain(n, at string, reqs ...string) (ret *Domain, err error) {
	d := c.ensureDomain(n)
	if d.currPhase > assert.DomainStart {
		err = errutil.New("can't add new dependencies to parent domains", d.name)
	} else {
		// domains are implicitly dependent on their parent domain
		if p, ok := c.processing.Top(); ok {
			reqs = append(reqs, p.name)
		}
		// probably asking for  trouble:
		// the tests have no top level domain (tapesrty) the way weave does
		// we still need them to wind up in the table eventually...
		if len(reqs) == 0 {
			err = c.writer.Domain(d.name, "", at)
		} else {
			for _, req := range reqs {
				if dep, ok := UniformString(req); !ok {
					err = errutil.New("invalid name", req)
					break
				} else if e := c.findDomainCycles(d.name, dep); e != nil {
					err = e
					break
				} else if e := c.writer.Domain(d.name, dep, at); e != nil {
					err = e
					break
				}
			}
		}
	}
	if err == nil {
		ret = d
	}
	return
}

// check before inserting a reference to avoid circularity;
// errors if one was detected.
func (c *Catalog) findDomainCycles(n, req string) (err error) {
	if n == req {
		err = errutil.Fmt("circular reference: %q can't depend on itself", n)
	} else {
		var exists bool
		if e := c.db.QueryRow(
			`select 1 
		from domain_tree 
		where base = ?1
		and uses = ?2
		and base != uses`, req, n).Scan(&exists); e != nil && e != sql.ErrNoRows {
			err = e
		} else if exists {
			err = errutil.Fmt("circular reference: %q requires %q", req, n)
		}
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

// tbd: maybe this could pull in the relevant domains;
// currently it happens after the domains have been activated
// and therefore compares everything to everything each time.
func findConflicts(db tables.Querier) (ret []conflict, err error) {
	if rows, e := db.Query(
		`select 'plural', a.domain, a.at, a.many, a.one
			from active_plurals as a 
			join active_plurals as b 
				using(many)
			where a.domain < b.domain 
			and a.one != b.one
		union all
		select 'opposite', a.domain, a.at, a.oneWord, a.otherWord 
			from active_rev as a 
			join active_rev as b 
				using(oneWord)
			where a.domain < b.domain 
			and a.otherWord != b.otherWord
		union all 
		select 'kind', a.domain, a.at, a.kind, ''
			from active_kinds as a 
			join active_kinds as b 
				using(kind)
			where a.domain < b.domain
		union all 
		select 'grammar', a.domain, a.at, a.name, ''
			from active_grammar as a 
			join active_grammar as b 
				using(name)
			where a.domain < b.domain 
			and a.prog != b.prog
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
		if deps, e := d.resolveKinds(); e != nil {
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
	if rows, e := c.db.Query(`select domain from domain_order`); e != nil {
		err = errutil.New("resolve domains", e)
	} else {
		var name string
		err = tables.ScanAll(rows, func() (err error) {
			if d, ok := c.GetDomain(name); !ok {
				err = errutil.Fmt("unexpected domain %q", name)
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
			if ks, e := d.resolveKinds(); e != nil {
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
