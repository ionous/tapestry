package eph

import (
	"errors"

	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

// returns true if newly added
func (d *Domain) AddPlural(plural, singular string) (okay bool) {
	return d.pairs.AddPair(plural, singular)
}

// use the domain rules ( and hierarchy ) to turn the passed plural into its singular form
func (d *Domain) Singularize(plural string) (ret string, err error) {
	if explict, e := d.FindSingluar(plural); e != nil {
		err = e
	} else if len(explict) > 0 {
		ret = explict
	} else {
		ret = inflect.Singularize(plural)
	}
	return
}

func (d *Domain) FindSingluar(plural string) (ret string, err error) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if n, ok := scope.pairs.FindSingular(plural); ok {
			ret, err = n, Visited
		}
		return
	}); e != nil && e != Visited {
		err = e
	}
	return
}

// while it'd probably be faster to do this while we assemble,
// keep this assembly separate from the writing produces nicer code and tests.
func (c *Catalog) WritePlurals(w Writer) (err error) {
	if deps, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			d := dep.Leaf().(*Domain)
			for i, p := range d.pairs.plural {
				s := d.pairs.singular[i]
				defs := d.phases[PluralPhase].defs
				at := defs[p].at
				if e := w.Write(mdl_plural, d.name, p, s, at); e != nil {
					err = errutil.Append(err, DomainError{d.name, e})
				}
			}
		}
	}
	return
}

func (el *EphPlurals) Phase() Phase { return PluralPhase }

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (el *EphPlurals) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if many, ok := UniformString(el.Plural); !ok {
		err = InvalidString(el.Plural)
	} else if one, ok := UniformString(el.Singular); !ok {
		err = InvalidString(el.Singular)
	} else {
		var de DomainError
		var conflict *Conflict
		if e := d.AddDefinition(many, at, one); e == nil {
			d.AddPlural(many, one)
		} else if !errors.As(e, &de) || !errors.As(de.Err, &conflict) {
			err = e // some unknown error?
		} else {
			switch conflict.Reason {
			case Redefined:
				// redefined definitions are only a problem in the same domain.
				// ( ie. we allow subdomains to reset / override the plurals )
				if d.name == de.Domain {
					err = e
				} else {
					d.AddPlural(many, one)
					// FIX! see Domain.AddDefinition
					// the earlier "AddDefinition" doesnt actually add it because this is a redefinition
					// *but* we actually do want that information....
					defs := d.phases[d.currPhase]
					defs.AddDefinition(many, Definition{at: at, value: one})
					d.phases[d.currPhase] = defs
					//
					LogWarning(e) // even though its okay, let the user know.
				}
			case Duplicated:
				// duplicated definitions are all okay;
				// but if its in a derived domain: let the user know.
				if de.Domain != d.name {
					LogWarning(e)
				}
			default:
				err = e // some unknown conflict?
			}
		}
	}

	return
}
