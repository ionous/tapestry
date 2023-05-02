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
	resolvedDomains cachedTable
	Errors          []error
	cursor          string
	writer          Writer
	run             rt.Runtime
	db              *sql.DB

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
		db:          db,
		writer:      mdl.Writer(ExecWriter(db).Write),
		run: qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			qna.DecodeNone("unsupported decoder"),
			qna.NewOptions()),
	}
}

func (c *Catalog) BeginDomain(name string, requires []string) (err error) {
	at := c.cursor
	if n, ok := UniformString(name); !ok {
		err = InvalidString(name)
	} else if reqs, e := UniformStrings(requires); e != nil {
		err = e // transform all the names first to determine any errors
	} else if d, e := c.EnsureDomain(n, at, reqs...); e != nil {
		err = e
	} else {
		c.processing.Push(d)
	}
	return
}

func (c *Catalog) EndDomain() (err error) {
	if _, ok := c.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else {
		c.processing.Pop()
	}
	return
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

// fix: this should be private; need to fix check first
// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) EnsureDomain(n, at string, reqs ...string) (ret *Domain, err error) {
	// find or create the domain
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{Requires: Requires{name: n, at: at}, catalog: c}
		if c.domains == nil {
			c.domains = map[string]*Domain{n: d}
		} else {
			c.domains[n] = d
		}
		ret = d
	}
	// add the passed requirements
	// ( it filters for uniqueness )
	for _, req := range reqs {
		if ret.name == req {
			err = errutil.New("domain depends on itself", ret.name)
			return // fix: early return
		} else {
			ret.AddRequirement(req)
		}
	}
	// we are dependent on the parent domain too
	// ( adding it last keeps it closer to the right side of the parent list )
	if p, ok := c.processing.Top(); ok {
		if ret.name == p.name {
			err = errutil.New("domain depends on itself", ret.name)
			return // fix: early return
		} else {
			ret.AddRequirement(p.name)
		}
	}
	return
}

var TempSplit = assert.PluralPhase + 1

// walk the domains and run the commands remaining in their queues
func (c *Catalog) AssembleCatalog() (err error) {
	// ds has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		// fix? create a dev/null writer instead of testing for nil?
		if c.writer != nil {
			// wip: trying to unwind the writes so that we write one domain at a time
			for _, deps := range ds {
				d := deps.Leaf().(*Domain) // panics if it fails
				if e := d.WriteDomain(c.writer); e != nil {
					err = e
					return // FIX
				}
				if _, e := c.run.ActivateDomain(d.name); e != nil {
					err = e
					return
				}
				// every phase inside this domain:
				for p := assert.Phase(0); p < TempSplit; p++ {
					switch p {
					// fix: this isnt supposed to be tied particularly to the plural phase.
					case assert.PluralPhase:
						if e := c.findConflicts(); e != nil {
							err = e
							return
						}
					}
					c.processing.Push(d)
					ctx := Weaver{d: d, phase: p, Runtime: c.run}
					if e := d.runPhase(&ctx); e != nil {
						err = e
						return // FIX
					}
					c.processing.Pop()
				}
			}
		}

		// FIX: want to add extensions to the db ( fields of kinds, etc. ) as we go.
		// DONT want to walk across all domains first.
		// walks across all domains for each phase to support things like fields:
		// which exist per kind but which can be added to by multiple domains.
	Loop:
		for p := TempSplit; p < assert.NumPhases; p++ {
			for _, deps := range ds {
				d := deps.Leaf().(*Domain) // panics if it fails
				c.processing.Push(d)
				ctx := Weaver{d: d, phase: p, Runtime: c.run}
				if e := d.runPhase(&ctx); e != nil {
					err = e
					break Loop
				} else if e := c.postPhase(p, d); e != nil {
					err = e
					break Loop
				}
				c.processing.Pop()
			}
		}
		if err == nil && c.writer != nil {
			for p := assert.Phase(0); p < assert.NumPhases; p++ {
				if e := c.WritePhase(p, c.writer); e != nil {
					err = e
					break
				}
			}
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

func findConflicts(db *sql.DB) (ret []conflict, err error) {
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
		err = e
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
	case assert.AncestryPhase:
		_, err = d.ResolveKinds()

	// PostFields: with the current domain looping pattern
	// we have to hit all locals from all domains before writing
	// and macro is after locals...
	case assert.MacroPhase:
		if deps, e := d.ResolveKinds(); e != nil {
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
		_, err = d.ResolveNouns()
	}
	return
}

func (c *Catalog) WritePhase(p assert.Phase, w Writer) (err error) {
	// switch or map better?
	switch p {
	case assert.AncestryPhase:
		err = c.WriteKinds(w)
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

// work out the hierarchy of all the domains, and return them in a list.
// the list has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
func (c *Catalog) ResolveDomains() (DependencyTable, error) {
	return c.resolvedDomains.resolve(func() (ret DependencyTable, err error) {
		m := TableMaker(len(c.domains))
		for _, d := range c.domains {
			m.ResolveDep(d) // accumulates any errors
		}
		return m.GetSortedTable()
	})
}

func (c *Catalog) ResolveKinds() (ret DependencyTable, err error) {
	var out DependencyTable
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range ds {
			d := dep.Leaf().(*Domain)
			if ks, e := d.ResolveKinds(); e != nil {
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

func (c *Catalog) ResolveNouns() (ret DependencyTable, err error) {
	// fix? is there anyway to make this more "automatically" resolve domains and kinds?
	var out DependencyTable
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range ds {
			d := dep.Leaf().(*Domain)
			if ns, e := d.ResolveNouns(); e != nil {
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
