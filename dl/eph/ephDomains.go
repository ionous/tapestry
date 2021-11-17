package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

type Domain struct {
	name, at string
	inflect  inflect.Ruleset
	eph      EphList
	deps     Dependencies
}

//
func (el *EphBeginDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	if kid := c.GetDomain(el.Name); len(kid.at) > 0 {
		err = errutil.New("domain", d.name, " at", d.at, "redeclared", at)
	} else {
		// initialize domain:
		kid.at = at
		kid.inflect = d.inflect
		// add any explicit dependencies
		for _, req := range el.Requires {
			d := c.GetDomain(req)
			kid.deps.AddDependency(d.name)
		}
		// we are dependent on the parent domain too
		// ( adding it last keeps it closer to the right side of the parent list )
		kid.deps.AddDependency(d.name)
		c.processing.Push(kid)
	}
	return
}

// pop the most recent domain
func (el *EphEndDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	kid := c.GetDomain(el.Name)
	if name := kid.name; name != d.name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", d.name)
	} else {
		c.processing.Pop()
	}
	return
}
