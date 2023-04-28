package eph

import (
	"database/sql"
	"log"

	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains         map[string]*Domain
	processing      DomainStack
	resolvedDomains cachedTable
	Errors          []error
	CommandBuilder
	cursor string
	writer Writer
	run    *qna.Runner
	qx     query.Queryx
}

func NewCatalog(db *sql.DB) *Catalog {
	var c Catalog
	qx, e := qdb.NewQueryx(db)
	if e != nil {
		panic(e)
	}
	c.cursor = "x"
	c.run = qna.NewRuntimeOptions(
		log.Writer(),
		qx,
		qna.DecodeNone("unsupported decoder"),
		qna.NewOptions())
	// set command builder
	c.q = WriterFun(c.writeEphemera)
	c.qx = qx
	c.writer = mdl.Writer(ExecWriter(db).Write)

	return &c
}

func (c *Catalog) Runtime() rt.Runtime {
	return c.run
}

// initializes and returns itself
func (c *Catalog) Weaver(db *sql.DB) *Catalog {
	return c
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

func (c *Catalog) writeEphemera(ep Ephemera) {
	var err error
	at := c.cursor
	if phase := ep.Phase(); phase == assert.DomainPhase {
		err = ep.Assemble(c, nil, at)
	} else {
		if d, ok := c.processing.Top(); !ok {
			err = errutil.New("no top domain")
		} else {
			err = d.QueueEphemera(at, ep)
		}
	}
	if err != nil {
		c.Errors = append(c.Errors, err)
	}
	return
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

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
		ret.AddRequirement(req)
	}
	// we are dependent on the parent domain too
	// ( adding it last keeps it closer to the right side of the parent list )
	if p, ok := c.processing.Top(); ok {
		ret.AddRequirement(p.name)
	}
	return
}

const TempSplit = assert.PluralPhase + 1

// walk the domains and run the commands remaining in their queues
func (c *Catalog) AssembleCatalog() (err error) {
	phaseActions := PhaseActions{
		assert.AncestryPhase: AncestryActions,
		assert.FieldPhase:    FieldActions,
		assert.NounPhase:     NounActions,
	}
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
					err = errutil.New("failed to write domain", d.name)
					return // FIX
				}
				if _, e := c.run.ActivateDomain(d.name); e != nil {
					err = errutil.New("error activating domain", d.name)
					return
				}
				// every phase inside this domain:
				for p := assert.Phase(0); p < TempSplit; p++ {
					switch p {
					case assert.PluralPhase:
						if e := c.qx.FindActiveConflicts(func(domain, key, value, at string) error {
							// fix: accumulate all conflicts
							return errutil.New("error checking plural conflicts")
						}); e != nil {
							err = e
							return
						}
					}
					if e := d.AssembleDomain(p, PhaseFlags{}); e != nil {
						err = errutil.New("error assembling phase", p, e)
						return // FIX
					}
				}
			}
		}

		// FIX: want to add extensions to the db ( fields of kinds, etc. ) as we go.
		// DONT want to walk across all domains first.
		// walks across all domains for each phase to support things like fields:
		// which exist per kind but which can be added to by multiple domains.
	Loop:
		for p := TempSplit; p < assert.NumPhases; p++ {
			act := phaseActions[p]
			for _, deps := range ds {
				d := deps.Leaf().(*Domain) // panics if it fails
				if e := d.AssembleDomain(p, act.Flags); e != nil {
					err = e
					break Loop
				} else if do := act.Do; do != nil {
					if e := do(d); e != nil {
						err = e
						break Loop
					}
				}
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

func (c *Catalog) WritePhase(p assert.Phase, w Writer) (err error) {
	// switch or map better?
	switch p {
	case assert.PluralPhase:
		if e := c.WriteOpposites(w); e != nil {
			err = e
		}
	case assert.AncestryPhase:
		err = c.WriteKinds(w)
	case assert.PropertyPhase:
	case assert.AspectPhase:
	case assert.FieldPhase:
		err = c.WriteFields(w)
	case assert.MacroPhase:
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

	case assert.RefPhase:
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
