package eph

import (
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains         map[string]*Domain
	Stack           DomainStack
	resolvedDomains cachedTable
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) EnsureDomain(n, at string) (ret *Domain) {
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
	return
}

// walk the domains and run the commands remaining in their queues
func (c *Catalog) AssembleCatalog(phaseActions PhaseActions) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds { // list of dependencies
			d := deps.Leaf().(*Domain) // panics if it fails
			if e := d.Assemble(phaseActions); e != nil {
				err = e
				break
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
