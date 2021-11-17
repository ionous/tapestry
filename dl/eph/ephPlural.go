package eph

import "errors"

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (el *EphPlural) Catalog(c *Catalog, d *Domain, at string) (err error) {
	var conflict *Conflict
	if e := c.CheckConflicts(d.name, mdl_plural, at, el.Plural, el.Singular); e == nil {
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
