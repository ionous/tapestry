package eph

import (
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	Catalog(c *Catalog, d *Domain, at string) error
}

// receive ephemera from the importer
type Catalog struct {
	domains Domains
}

//
func (c *Catalog) AddEphemera(pat EphAt) (err error) {
	if d, ok := c.domains.processing.Top(); !ok {
		err = errutil.New("no domain")
	} else if e := pat.Eph.Catalog(c, d, pat.At); e != nil {
		err = e
	}
	return
}

//
func (el *EphBeginDomain) Catalog(c *Catalog, p *Domain, at string) (err error) {
	if kid := c.domains.GetDomain(el.Name); len(kid.at) > 0 {
		err = errutil.New("domain", p.name, " at", p.at, "redeclared", at)
	} else {
		// initialize domain:
		kid.at = at
		kid.inflect = p.inflect
		kid.deps.add(p) // we are dependent on the parent domain
		// add all the other dependencies too
		for _, req := range el.Requires {
			kid.deps.add(c.domains.GetDomain(req))
		}
		c.domains.processing.Push(kid)
	}
	return
}

// pop the most recent domain
func (el *EphEndDomain) Catalog(c *Catalog, d *Domain, at string) (err error) {
	if d.name != el.Name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", d.name)
	} else {
		c.domains.processing.Pop()
	}
	return
}

// add to the plurals to the database and remember the plural for the current domain's set of rules
// eph_plural: plural, singular, domain, path
func (el *EphPlural) Catalog(c *Catalog, d *Domain, at string) (err error) {
	d.inflect.AddPluralExact(el.Singular, el.Plural, true)
	return
}
