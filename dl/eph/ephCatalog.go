package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains         map[string]*Domain
	processing      DomainStack
	plurals         PluralTable
	phase           Phase
	resolvedDomains cachedTable
}

// use the domain rules ( and hierarchy ) to turn the passed plural into its singular form
func (c *Catalog) Singularize(domain, plural string) (ret string, err error) {
	if explict, e := c.plurals.FindSingular((*catDependencyFinder)(c), domain, plural); e != nil {
		err = e
	} else if len(explict) > 0 {
		ret = explict
	} else {
		ret = inflect.Singularize(plural)
	}
	return
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
		d = &Domain{name: n, at: at, catalog: c}
		if c.domains == nil {
			c.domains = map[string]*Domain{n: d}
		} else {
			c.domains[n] = d
		}
		ret = d
	}
	return
}

func (c *Catalog) AddEphemera(ephs ...EphAt) (err error) {
Out:
	for _, ephAt := range ephs {
		if d, ok := c.processing.Top(); !ok {
			err = errutil.New("no domain")
			break
		} else if currPhase, phase := c.phase, ephAt.Eph.Phase(); currPhase > phase {
			err = errutil.New("unexpected phase")
			break
		} else if phase == DomainPhase {
			// fix: queue first, and then run?
			for _, ephAt := range ephs {
				if e := ephAt.Eph.Assemble(c, d, ephAt.At); e != nil {
					err = errutil.Append(err, e)
					break Out
				}
			}
		} else {
			// fix? consider storing sorted by phase? or storing linear in order and scanning by phase?
			// that way we dont need all the separate lists and we can append....
			els := d.phases[phase]
			els.eph = append(els.eph, ephAt)
			d.phases[phase] = els
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

//  traverse the domains and then kinds in a reasonable order
func (cat *Catalog) WriteFields(w Writer) (err error) {
	if ds, e := cat.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range ds {
			d := dep.Leaf().(*Domain)
			if ks, e := d.ResolveKinds(); e != nil {
				err = e
				break
			} else {
				for _, kep := range ks {
					k := kep.Leaf().(*ScopedKind)
					for _, f := range k.fields {
						f.Write(&partialFields{w: w, fields: []interface{}{d.Name(), k.Name()}})
					}
				}
			}
		}
	}
	return
}

type partialFields struct {
	w      Writer
	fields []interface{}
}

func (p *partialFields) Write(q string, args ...interface{}) error {
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
