package eph

import (
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	name, at string
	catalog  *Catalog
	phases   [NumPhases][]EphAt
	reqs     Requires // other domains this needs ( can have multiple direct parents )
	defs     Artifacts
	kinds    ScopedKinds
}

// implement the Dependency interface
func (d *Domain) Name() string                           { return d.name }
func (d *Domain) AddRequirement(name string)             { d.reqs.AddRequirement(name) }
func (d *Domain) GetDependencies() (Dependencies, error) { return d.reqs.GetDependencies() }
func (d *Domain) Resolve() (ret Dependencies, err error) {
	if len(d.at) == 0 {
		err = DomainError{d.name, errutil.New("never defined")}
	} else if d.catalog == nil {
		err = DomainError{d.name, errutil.New("no catalog")}
	} else {
		ret, err = d.reqs.Resolve(d, (*catDependencyFinder)(d.catalog))
	}
	return
}

// return the uniformly named domain ( if it exists )
func (d *Domain) GetKind(n string) (ret *ScopedKind, okay bool) {
	if k, ok := d.kinds[n]; ok {
		ret, okay = k, true
	} else if d.catalog != nil { // skip for tests.
		// if not in this domain, then maybe in a parent domain....
		// ( dont force resolve here, if its not resolved... then stop trying )
		if deps, e := d.reqs.GetDependencies(); e != nil {
			LogWarning(e)
		} else {
			list := deps.Ancestors()
			for i, cnt := 0, len(list); i < cnt; i++ {
				el := list[cnt-i-1].(*Domain)
				if k, ok := el.kinds[n]; ok {
					ret, okay = k, true
					break
				}
			}
		}
	}
	return
}

// return the uniformly named domain ( creating it in this domain if necessary )
func (d *Domain) EnsureKind(n, at string) (ret *ScopedKind) {
	if k, ok := d.GetKind(n); ok {
		ret = k
	} else {
		k = &ScopedKind{name: n, at: at, domain: d}
		if d.kinds == nil {
			d.kinds = map[string]*ScopedKind{n: k}
		} else {
			d.kinds[n] = k
		}
		ret = k
	}
	return
}

// // distill a tree of kinds into a set of names and their hierarchy
func (d *Domain) ResolveKinds() (ret DependencyTable, err error) {
	m := TableMaker(len(d.kinds))
	for n, k := range d.kinds {
		if res, ok := m.ResolveDep(k); ok && len(res.Parents()) > 1 {
			err = errutil.Append(err, errutil.New(n, "has more than one parent"))
		}
	}
	if dt, e := m.GetSortedTable(); e != nil {
		err = errutil.Append(err, e)
	} else {
		ret = dt
	}
	return
}

// EphBeginDomain
func (el *EphBeginDomain) Phase() Phase { return DomainPhase }

//
func (el *EphBeginDomain) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if n, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if kid, ok := c.GetDomain(n); ok {
		err = errutil.New("domain", n, " at", kid.at, "redeclared", kid.at)
	} else {
		kid := c.EnsureDomain(n, at)
		// add any explicit dependencies
		for _, req := range el.Requires {
			if sub, ok := UniformString(req); !ok {
				err = errutil.Append(err, InvalidString(req))
			} else {
				kid.AddRequirement(sub)
			}
		}
		if err == nil {
			// we are dependent on the parent domain too
			// ( adding it last keeps it closer to the right side of the parent list )
			kid.AddRequirement(d.name)
			c.processing.Push(kid)
		}
	}
	return
}

// EphEndDomain
func (el *EphEndDomain) Phase() Phase { return DomainPhase }

// pop the most recent domain
func (el *EphEndDomain) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	if n, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if n != d.name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", d.name)
	} else {
		c.processing.Pop()
	}
	return
}
