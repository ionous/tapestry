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
	resolvedDomains DependencyTable
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

// used by ephemera during assembly to record some piece of information
// that would cause problems it were specified differently elsewhere.
// ex. some in game password specified as the word "secret" in one place, but "mongoose" somewhere else.
func (c *Catalog) AddDefinition(domain, cat, at, key, value string) (err error) {
	if d, ok := c.GetDomain(domain); !ok {
		err = errutil.New("unknown domain", domain)
	} else if ds, e := d.GetDependencies(); e != nil {
		err = e
	} else {
		key := CategoryKey{cat, key}
		err = CheckConflicts(ds.FullTree(), (*catArtifactFinder)(c), key, at, value)
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

// creates domains, suspends all other ephemera until the domains are resolved.
func (c *Catalog) AddEphemera(ephAt EphAt) (err error) {
	if d, ok := c.processing.Top(); !ok {
		err = errutil.New("no domain")
	} else if currPhase, phase := c.phase, ephAt.Eph.Phase(); currPhase > phase {
		err = errutil.New("unexpected phase")
	} else if phase == DomainPhase {
		err = ephAt.Eph.Assemble(c, d, ephAt.At)
	} else {
		d.phases[phase] = append(d.phases[phase], ephAt)
	}
	return
}

// work out the hierarchy of all the domains, and return them in a list.
// the list has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
func (c *Catalog) ResolveDomains() (ret DependencyTable, err error) {
	if len(c.resolvedDomains) != 0 {
		ret = c.resolvedDomains
	} else {
		m := TableMaker(len(c.domains))
		for _, d := range c.domains {
			m.ResolveDep(d) // accumulates any errors
		}
		if res, e := m.GetSortedTable(); e != nil {
			err = e
		} else {
			ret, c.resolvedDomains = res, res
		}
	}
	return
}

// walk the domains and run the commands remaining in their queues
func (c *Catalog) ProcessDomains(phaseActions PhaseActions) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds {
			if e := c.AssembleDomain(deps, phaseActions); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (c *Catalog) AssembleDomain(deps Dependencies, phaseActions PhaseActions) (err error) {
	n := deps.Leaf().Name()
	if d, ok := c.GetDomain(n); !ok {
		err = errutil.New("unknown domain", n)
	} else {
		c.processing = DomainStack{d} // so ephemera can add other ephemera
		if e := c.checkRivals(deps); e != nil {
			err = e
		} else {
			for w, ephlist := range d.phases {
				for _, el := range ephlist {
					if e := el.Eph.Assemble(c, d, el.At); e != nil {
						err = errutil.Append(err, e)
					}
				}
				if err != nil {
					break
				} else if act, ok := phaseActions[Phase(w)]; ok {
					if e := act(c, d); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}

// used by assembler to check that domains with multiple parents don't contain conflicting information.
// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
func (c *Catalog) checkRivals(res Dependencies) (err error) {
	if parents := res.Parents(); len(parents) > 1 {
		def := make(Artifacts) // start with nothing and merge in to check for artifacts
		for _, p := range parents {
			if d, ok := p.(*Domain); ok {
				if e := def.Merge(d.defs); e != nil {
					err = DomainError{d.name, e}
					break
				}
			}
		}
	}
	return
}

// private helper to make the catalog compatible with the ArtifactFinder ( for domains )
type catArtifactFinder Catalog

func (c *catArtifactFinder) GetArtifacts(name string) (ret *Artifacts, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.defs, true
	}
	return
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) FindDependency(name string) (ret Dependency, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = d, true
	}
	return
}
