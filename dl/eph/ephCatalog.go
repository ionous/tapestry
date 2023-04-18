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
}

// fix? consider moving domain error to catalog processing internals ( and removing explicit external use )
type DomainError struct {
	Domain string
	Err    error
}

func (n DomainError) Error() string {
	return errutil.Sprintf("%v in domain %q", n.Err, n.Domain)
}
func (n DomainError) Unwrap() error {
	return n.Err
}

func (c *Catalog) AddEphemera(at string, ep Ephemera) (err error) {
	// fix: queue first, and then run?
	// ( would need an initial "top" domain to queue into i think )
	if phase := ep.Phase(); phase == assert.DomainPhase {
		err = ep.Assemble(c, nil, at)
	} else {
		if d, ok := c.processing.Top(); !ok {
			err = errutil.New("no top domain")
		} else {
			err = d.AddEphemera(at, ep)
		}
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
func (c *Catalog) AssembleCatalog(phaseActions PhaseActions) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
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
