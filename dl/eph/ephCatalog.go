package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
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
}

// initializes and returns itself
func (c *Catalog) Weaver() *Catalog {
	if c.q == nil {
		c.cursor = "x"
		c.q = WriterFun(c.writeEphemera)
	}
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

// walk the domains and run the commands remaining in their queues
func (c *Catalog) AssembleCatalog(w Writer, phaseActions PhaseActions) (err error) {
	// ds has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		// fix? create a dev/null writer instead of testing for nil?
		if w != nil {
			// wip: trying to unwind the writes so that we write one domain at a time
			for _, deps := range ds {
				d := deps.Leaf().(*Domain) // panics if it fails
				if e := d.WriteDomain(w); e != nil {
					err = e
					return // FIX
				}
			}
		}

		// FIX: want to add extensions to the db ( fields of kinds, etc. ) as we go.
		// DONT want to walk across all domains first.
		// walks across all domains for each phase to support things like fields:
		// which exist per kind but which can be added to by multiple domains.
	Loop:
		for w := assert.Phase(0); w < assert.NumPhases; w++ {

			act := phaseActions[w]
			for _, deps := range ds {
				d := deps.Leaf().(*Domain) // panics if it fails
				if e := d.AssembleDomain(w, act.Flags); e != nil {
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
		if err == nil && w != nil {
			for p := assert.Phase(0); p < assert.NumPhases; p++ {
				if e := c.WritePhase(p, w); e != nil {
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
		if e := c.WritePlurals(w); e != nil {
			err = e
		} else if e := c.WriteOpposites(w); e != nil {
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
