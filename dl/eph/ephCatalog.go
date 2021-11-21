package eph

import (
	"log"
	"sort"
	"strings"

	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	Writer     // not a huge fan of this here.... hrm...
	domains    map[string]*Domain
	processing DomainStack
	conflicts  DomainConflicts
}

// helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) GetDependencies(name string) (ret *Dependencies, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.deps, true
	}
	return
}

func (c *Catalog) Warn(e error) {
	log.Println(e) // for now good enough
}

func (c *Catalog) CheckConflicts(name, cat, at, key, value string) (err error) {
	if c.conflicts == nil {
		c.conflicts = make(DomainConflicts)
	}
	return c.conflicts.CheckConflicts(name, (*catDependencyFinder)(c), cat, at, key, value)
}

// primarily for testing: return list of all the domains that the passed uniformly named domain requires
func (c *Catalog) GetDependentDomains(n string) (ret []string, err error) {
	if dep, e := c.getDependentDomains(n); e != nil {
		err = e
	} else {
		ret = dep.Ancestors(true)
	}
	return
}

func (c *Catalog) getDependentDomains(n string) (ret ResolvedDependencies, err error) {
	if d, ok := c.domains[n]; !ok {
		err = errutil.New("unknown domains have no dependencies")
	} else {
		ret, err = d.deps.Resolve(d.name, (*catDependencyFinder)(c))
	}
	return
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) GetDomain(n string) (ret *Domain) {
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n}
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
		if phase == Domains {
			err = ephAt.Eph.Catalog(c, d, ephAt.At)
		} else {
			d.phases[phase] = append(d.phases[phase], ephAt)
		}
	}
	return
}

// Process
// func (c *Catalog) EphemeraAssemble() (err error) {
// 	// we need to handle the case(s) where one parent domain contains ephemera that conflicts with another parent domain:
// 	// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
// 	// i dont really have a good way of doing this.... just have to do it manually.
// 	for n, _ := range c.domains {
// 		if res, e := c.getDependentDomains(n); e != nil {
// 			err = e
// 		} else {
// 			if parents := res.Ancestors(false); len(parents) > 0 {
// 				def := make(Definitions) // start with nothing and merge in to check for conflicts
// 				for _, p := range parents {
// 					pdef := c.conflicts[p]
// 					if e := def.Merge(p, pdef); e != nil {
// 						err = e
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}
// 	// need to walk over domains --

// 	return
// }

// for each domain in the passed list, output its full ancestry tree ( or just its parents )
func (c *Catalog) WriteDomains(list ResolvedDomains, fullTree bool) (err error) {
	for _, n := range list {
		if deps, e := c.getDependentDomains(n); e != nil {
			err = errutil.Append(err, e)
		} else if e := c.Write(mdl_domain, n, strings.Join(deps.Ancestors(fullTree), ",")); e != nil {
			err = errutil.Append(err, errutil.New("domain", n, "couldn't write", e))
		}
	}
	return
}

type ResolvedDomains []string

func (c *Catalog) ResolveAllDomains() (ret ResolvedDomains, err error) {
	names := make([]string, 0, len(c.domains))
	deps := make([]int, 0, len(c.domains))
	// walk all domains in the map
	for n, d := range c.domains {
		if len(d.at) == 0 {
			err = errutil.Append(err, errutil.New("domain never declared", d.name))
		} else if dep, e := c.getDependentDomains(n); e != nil {
			err = errutil.Append(err, e)
		} else {
			names = append(names, n) // add the depth of the tree
			deps = append(deps, len(dep.Ancestors(true)))
		}
	}
	if err == nil {
		sort.Sort(&nameDeps{names, deps})
		ret = names
	}
	return
}

// private helper to sort domains by least to most dependencies
type nameDeps struct {
	names []string
	deps  []int
}

func (n *nameDeps) Len() int {
	return len(n.names)
}
func (n *nameDeps) Less(i, j int) (okay bool) {
	if adep, bdep := n.deps[i], n.deps[j]; adep < bdep {
		okay = true
	} else if adep == bdep {
		if a, b := n.names[i], n.names[j]; a > b {
			okay = true
		}
	}
	return
}
func (n *nameDeps) Swap(i, j int) {
	n.names[i], n.names[j] = n.names[j], n.names[i]
	n.deps[i], n.deps[j] = n.deps[j], n.deps[i]
}
