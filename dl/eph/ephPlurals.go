package eph

import (
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

// returns true if newly added
func (d *Domain) AddPlural(plural, singular string) (okay bool) {
	return d.plural.AddPair(plural, singular)
}

// use the domain rules ( and hierarchy ) to turn the passed singular word into its plural form.
// the way plurals are defined, there can be more than one plural word for a given singular word.
// in that case, attempts to pick one.
func (d *Domain) Pluralize(singular string) (ret string, err error) {
	// dont bother with one letter kinds ( ex. tests )
	if len(singular) < 2 {
		ret = singular
	} else if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if n, ok := scope.plural.FindPlural(singular); ok {
			ret, err = n, Visited
		}
		return
	}); e == nil { // not found
		ret = inflect.Pluralize(singular)
	} else if e != Visited {
		err = e
	}
	return
}

// use the domain rules ( and hierarchy ) to turn the passed plural into its singular form
// the way plurals are defined, there can be more than one plural word for a given singular word.
// in that case, attempts to pick one.
func (d *Domain) Singularize(plural string) (ret string, err error) {
	// dont bother with one letter kinds ( ex. tests )
	if len(plural) < 2 {
		ret = plural
	} else if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if n, ok := scope.plural.FindSingular(plural); ok {
			ret, err = n, Visited
		}
		return
	}); e == nil { // not found
		ret = inflect.Singularize(plural)
	} else if e != Visited {
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
			for i, p := range d.plural.plural {
				s := d.plural.singular[i]
				def := d.GetDefinition(MakeKey("plurals", p))
				if e := w.Write(mdl.Plural, d.name, p, s, def.at); e != nil {
					err = errutil.Append(err, DomainError{d.name, e})
				}
			}
		}
	}
	return
}

func (op *EphPlurals) Phase() Phase { return PluralPhase }

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (op *EphPlurals) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if many, ok := UniformString(op.Plural); !ok {
		err = InvalidString(op.Plural)
	} else if one, ok := UniformString(op.Singular); !ok {
		err = InvalidString(op.Singular)
	} else if ok, e := d.RefineDefinition(MakeKey("plurals", many), at, one); e != nil {
		err = e
	} else if ok {
		if !d.AddPlural(many, one) {
			err = errutil.New("couldnt add plurals", many, one)
		}
	}
	return
}
