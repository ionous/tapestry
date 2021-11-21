package eph

import (
	"errors"
)

func (el *EphPlural) Phase() Phase { return Plurals }

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (el *EphPlural) Catalog(c *Catalog, d *Domain, at string) (err error) {
	if many, ok := UniformString(el.Plural); !ok {
		err = InvalidString(el.Plural)
	} else if one, ok := UniformString(el.Singular); !ok {
		err = InvalidString(el.Singular)
	} else {
		var conflict *Conflict
		if e := c.CheckConflicts(d.name, mdl_plural, at, many, one); e == nil {
			writePlural(c, d, many, one)
		} else if !errors.As(e, &conflict) {
			err = e // some unknown error?
		} else if conflict.Reason != Duplicated {
			// duplicated definitions are okay, but we dont need to store them.
			// redefined definitions are only a problem in this domain
			// ( we test for !redefined in case there's some unexpected error code. )
			if d.name == conflict.Domain || conflict.Reason != Redefined {
				err = e
			} else {
				writePlural(c, d, many, one)
				c.Warn(e) // redefined vs an earlier domain: let the user know.
			}
		}
	}
	return
}

func writePlural(c *Catalog, d *Domain, many, one string) (err error) {
	if e := c.Write(mdl_plural, d.name, many, one); e != nil {
		err = e
	} else {
		// d.inflect.AddPluralExact(one, many, true)
	}
	return
}
