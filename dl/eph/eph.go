package eph

import (
	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	Catalog(c *Catalog, d *Domain, at string) error
}

//
func (el *EphBeginDomain) Catalog(c *Catalog, p *Domain, at string) (err error) {
	name := lang.Underscore(el.Name)
	if kid := c.GetDomain(name); len(kid.at) > 0 {
		err = errutil.New("domain", p.name, " at", p.at, "redeclared", at)
	} else {
		// initialize domain:
		kid.originalName = el.Name
		kid.at = at
		kid.inflect = p.inflect
		kid.deps.add(p) // we are dependent on the parent domain
		// add any explicit dependencies too
		for _, req := range el.Requires {
			name := lang.Underscore(req)
			kid.deps.add(c.GetDomain(name))
		}
		c.processing.Push(kid)
	}
	return
}

// pop the most recent domain
func (el *EphEndDomain) Catalog(c *Catalog, p *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	if name := lang.Underscore(el.Name); name != p.name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", p.name)
	} else {
		c.processing.Pop()
	}
	return
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// eph_plural: plural, singular, domain, path
// note: i can actually add things to the dbs and resolve the domain order later.
func (el *EphPlural) Catalog(c *Catalog, d *Domain, at string) (err error) {

	// next

	d.inflect.AddPluralExact(el.Singular, el.Plural, true)
	return
}
