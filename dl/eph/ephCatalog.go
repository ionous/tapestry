package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains    map[string]*Domain
	processing DomainStack
	artifacts  DomainArtifacts
	plurals    PluralTable
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
func (c *Catalog) CheckConflict(name, cat, at, key, value string) (err error) {
	if c.artifacts == nil {
		c.artifacts = make(DomainArtifacts)
	}
	return c.artifacts.CheckConflict(name, (*catDependencyFinder)(c), cat, at, key, value)
}

func (c *Catalog) GetDependentDomains(n string) (ret Dependents, err error) {
	if d, ok := c.domains[n]; !ok {
		err = errutil.New("unknown domains have no dependencies")
	} else {
		ret, err = d.deps.Resolve(d.name, (*catDependencyFinder)(c))
	}
	return
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) EnsureDomain(n string) (ret *Domain) {
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n /*, finder: c*/}
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
	} else {
		phase := ephAt.Eph.Phase()
		if phase == DomainPhase {
			err = ephAt.Eph.Assemble(c, d, ephAt.At)
		} else {
			d.phases[phase] = append(d.phases[phase], ephAt)
		}
	}
	return
}

// work out the hierarchy of all the domains, and return them in a list.
// the list has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
func (c *Catalog) ResolveDomains() (ret DependencyTable, err error) {
	out := make([]Dependents, 0, len(c.domains))
	// walk all domains in the map
	for n, d := range c.domains {
		if len(d.at) == 0 {
			err = errutil.Append(err, errutil.New("domain never declared", d.name))
		} else if dep, e := c.GetDependentDomains(n); e != nil {
			err = errutil.Append(err, e)
		} else {
			out = append(out, dep)
		}
	}
	if err == nil {
		ret = out
		ret.SortTable()
	}
	return
}

// for each domain, determine the kinds that it defined
func (c *Catalog) ResolveKinds(ds DependencyTable) (ret ResolvedKinds, err error) {
	var out ResolvedKinds
	for _, n := range ds {
		if d, ok := c.GetDomain(n.Name()); !ok {
			err = errutil.Fmt("unknown domain %q", n.Name())
			break
		} else if e := d.kinds.ResolveKinds(&out); e != nil {
			err = e
			break
		}
	}
	return
}

// walk the domains and run the commands remaining in their queues
func (c *Catalog) ProcessDomains(phaseActions PhaseActions) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, res := range ds {
			if d, ok := c.GetDomain(res.Name()); !ok {
				err = errutil.New("unknown domain", res.Name())
				break
			} else if e := c.checkRivals(res); e != nil {
				err = e
				break
			} else {
				for phase, ephlist := range d.phases {
					for _, el := range ephlist {
						if e := el.Eph.Assemble(c, d, el.At); e != nil {
							err = errutil.Append(err, e)
						}
					}
					if err != nil {
						break
					} else if act, ok := phaseActions[Phase(phase)]; ok {
						if e := act(c, d); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
	}
	return
}

// used by assembler to check that domains with multiple parents don't contain conflicting information.
// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
func (c *Catalog) checkRivals(res Dependents) (err error) {
	if parents := res.Ancestors(false); len(parents) > 1 {
		def := make(Definitions) // start with nothing and merge in to check for artifacts
		for _, p := range parents {
			// get the artifacts built from the named domain p
			if pdef, ok := c.artifacts[p]; ok {
				if e := def.Merge(pdef); e != nil {
					err = DomainError{p, e}
					break
				}
			}
		}
	}
	return
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) GetRequirements(name string) (ret *Requires, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.deps, true
	}
	return
}
