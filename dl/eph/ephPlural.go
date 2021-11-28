package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

// domain name to plural lookup
type PluralTable map[string]PluralPairs

// returns true if newly added
func (pd *PluralTable) AddPair(domain, plural, singular string) (okay bool) {
	if *pd == nil {
		*pd = make(map[string]PluralPairs)
	}
	pairs := (*pd)[domain] // this is a copy
	if pairs.AddPair(plural, singular) {
		(*pd)[domain] = pairs // so we have to write any changes back
		okay = true
	}
	return
}

func (pd PluralTable) FindSingular(names DependencyFinder, domain, plural string) (ret string, err error) {
	if s, ok := pd.findSingular(domain, plural); ok {
		ret = s
	} else if dep, ok := names.FindDependency(domain); !ok {
		err = errutil.New("unknown dependency", domain)
	} else if requires, e := dep.GetDependencies(); e != nil {
		err = e
	} else {
		search := requires.Ancestors()
		for {
			if cnt := len(search); cnt == 0 {
				break
			} else {
				dep, search = search[cnt-1], search[:cnt-1]
				if s, ok := pd.findSingular(dep.Name(), plural); ok {
					ret = s
					break
				}
			}
		}
	}
	return
}

func (pd PluralTable) findSingular(n, plural string) (ret string, okay bool) {
	if ps, ok := pd[n]; ok {
		ret, okay = ps.FindSingular(plural)
	}
	return
}

// while it'd probably be faster to do this while we assemble,
// keep this assembly separate from the writing produces nicer code and tests.
func (pd PluralTable) WritePlurals(w Writer) (err error) {
	for d, ps := range pd {
		for i, p := range ps.plural {
			s := ps.singular[i]
			if e := w.Write(mdl_plural, d, p, s); e != nil {
				err = errutil.Append(err, DomainError{d, e})
			}
		}
	}
	return
}

func (el *EphPlural) Phase() Phase { return PluralPhase }

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (el *EphPlural) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if many, ok := UniformString(el.Plural); !ok {
		err = InvalidString(el.Plural)
	} else if one, ok := UniformString(el.Singular); !ok {
		err = InvalidString(el.Singular)
	} else {
		var de DomainError
		var conflict *Conflict
		if e := c.AddDefinition(d.name, mdl_plural, at, many, one); e == nil {
			c.plurals.AddPair(d.name, many, one)
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
					c.plurals.AddPair(d.name, many, one)
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
