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
	artifacts  DomainArtifacts
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
	log.Println("Warning:", e) // for now good enough
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

func (c *Catalog) ResolveDomains() (ret ResolvedDomains, err error) {
	out := make([]*Domain, 0, len(c.domains))
	deps := make([]int, 0, len(c.domains))
	// walk all domains in the map
	for n, d := range c.domains {
		if len(d.at) == 0 {
			err = errutil.Append(err, errutil.New("domain never declared", d.name))
		} else if dep, e := c.getDependentDomains(n); e != nil {
			err = errutil.Append(err, e)
		} else {
			out = append(out, d) // add the depth of the tree
			deps = append(deps, len(dep.Ancestors(true)))
		}
	}
	if err == nil {
		sort.Sort(&nameDeps{out, deps})
		ret = out
	}
	return
}

// used by assembler to check that domains with multiple parents don't contain conflicting information.
// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
func (c *Catalog) checkRivals(d *Domain) (err error) {
	// if this succeeds, then our dependencies have been at least resolved; and,
	// since a child domain by definition has a greater depth than a parent dependency
	// we also know that its parent domains have been processed.
	if res, e := d.GetDependencies(); e != nil {
		err = e
	} else if parents := res.Ancestors(false); len(parents) > 1 {
		def := make(Definitions) // start with nothing and merge in to check for artifacts
		for _, p := range parents {
			// get the artifacts built from the named domain p
			if pdef, ok := c.artifacts[p]; ok {
				if e := def.Merge(pdef, c.Warn); e != nil {
					err = DomainError{p, e}
					break
				}
			}
		}
	}
	return
}

// for each domain in the passed list, output its full ancestry tree ( or just its parents )
func (c *Catalog) WriteDomains(ds ResolvedDomains, fullTree bool) (err error) {
	for _, d := range ds {
		if deps, e := d.GetDependencies(); e != nil {
			err = errutil.Append(err, e)
		} else if e := c.Write(mdl_domain, d.name, strings.Join(deps.Ancestors(fullTree), ",")); e != nil {
			err = errutil.Append(err, errutil.New("domain", d.name, "couldn't write", e))
		}
	}
	return
}

// private helper to sort domains by least to most dependencies
type nameDeps struct {
	domains []*Domain
	deps    []int
}

func (n *nameDeps) Len() int {
	return len(n.deps)
}
func (n *nameDeps) Less(i, j int) (okay bool) {
	if adep, bdep := n.deps[i], n.deps[j]; adep < bdep {
		okay = true
	} else if adep == bdep {
		if a, b := n.domains[i], n.domains[j]; a.name < b.name {
			okay = true
		}
	}
	return
}
func (n *nameDeps) Swap(i, j int) {
	n.domains[i], n.domains[j] = n.domains[j], n.domains[i]
	n.deps[i], n.deps[j] = n.deps[j], n.deps[i]
}
