package eph

import (
	"github.com/ionous/errutil"
)

type Domain struct {
	name, at string
	phases   [NumPhases][]EphAt
	deps     Dependencies
}

func (d *Domain) GetDependencies() (ResolvedDependencies, error) {
	return d.deps.GetDependencies()
}

func (el *EphBeginDomain) Phase() Phase { return Domains }

//
func (el *EphBeginDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	if n, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if kid := c.GetDomain(n); len(kid.at) > 0 {
		err = errutil.New("domain", n, " at", d.at, "redeclared", at)
	} else {
		// initialize domain:
		kid.at = at
		// kid.inflect = d.inflect
		// add any explicit dependencies
		for _, req := range el.Requires {
			if sub, ok := UniformString(req); !ok {
				err = errutil.Append(err, InvalidString(req))
			} else {
				d := c.GetDomain(sub)
				kid.deps.AddDependency(d.name)
			}
		}
		if err == nil {
			// we are dependent on the parent domain too
			// ( adding it last keeps it closer to the right side of the parent list )
			kid.deps.AddDependency(d.name)
			c.processing.Push(kid)
		}
	}
	return
}

func (el *EphEndDomain) Phase() Phase { return Domains }

// pop the most recent domain
func (el *EphEndDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
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
