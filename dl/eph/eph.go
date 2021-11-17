package eph

import (
	"errors"

	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	Catalog(c *Catalog, d *Domain, at string) error
}

//
func (el *EphBeginDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	name := lang.Underscore(el.Name)
	if kid := c.GetDomain(name); len(kid.at) > 0 {
		err = errutil.New("domain", d.name, " at", d.at, "redeclared", at)
	} else {
		// initialize domain:
		kid.originalName = el.Name
		kid.at = at
		kid.inflect = d.inflect
		// add any explicit dependencies
		for _, req := range el.Requires {
			name := lang.Underscore(req)
			kid.deps.add(c.GetDomain(name))
		}
		// we are dependent on the parent domain too
		// ( adding it last keeps it closer to the right side of the parent list )
		kid.deps.add(d)
		c.processing.Push(kid)
	}
	return
}

// pop the most recent domain
func (el *EphEndDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	if name := lang.Underscore(el.Name); name != d.name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", d.name)
	} else {
		c.processing.Pop()
	}
	return
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (el *EphPlural) Catalog(c *Catalog, d *Domain, at string) (err error) {
	var conflict *Conflict
	if e := d.CheckConflicts(mdl_plural, at, el.Plural, el.Singular); e == nil {
		el.store(c, d)
	} else if !errors.As(e, &conflict) {
		err = e // some unknown error?
	} else if conflict.Reason != Duplicated {
		// duplicated definitions are okay, but we dont need to store them.
		// redefined definitions are only a problem in this domain
		// ( we test for !redefined in case there's some unexpected error code. )
		if d.name == conflict.Domain || conflict.Reason != Redefined {
			err = e
		} else {
			el.store(c, d)
			c.Warn(e) // redefined vs an earlier domain: let the user know.
		}
	}
	return
}

func (el *EphPlural) store(c *Catalog, d *Domain) (err error) {
	if e := c.Write(mdl_plural, d.name, el.Plural, el.Singular); e != nil {
		err = e
	} else {
		d.inflect.AddPluralExact(el.Singular, el.Plural, true)
	}
	return
}
