package weave

import (
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// fix: remove this function
func (d *Domain) Pluralize(singular string) (ret string, err error) {
	if d.catalog != nil {
		ret = d.catalog.run.PluralOf(singular)
	} else {
		ret = singular // for tests
	}
	return
}

// fix: remove this function
func (d *Domain) Singularize(plural string) (ret string, err error) {
	if d.catalog != nil {
		ret = d.catalog.run.SingularOf(plural)
	} else {
		ret = plural // for tests
	}
	return
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
//
// tbd: consider appending the origin (at) to store the location of each definition?
func (cat *Catalog) AssertPlural(opSingular, opPlural string) error {
	return cat.Schedule(assert.PluralPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if plural, ok := UniformString(opPlural); !ok {
			err = InvalidString(opPlural)
		} else if singular, ok := UniformString(opSingular); !ok {
			err = InvalidString(opSingular)
		} else {
			duplicated := false
			if e := cat.qx.FindPluralDefinitions(plural, func(domain, one, at string) (err error) {
				why := Redefined
				if one == singular {
					why = Duplicated
					duplicated = true
				}
				key := MakeKey("plurals", plural)
				e := domainError{Domain: d.name, Err: newConflict(key, why,
					Definition{key, at, one},
					singular,
				)}
				if one == singular {
					LogWarning(e)
				} else {
					err = e
				}
				return err
			}); e != nil {
				err = e
			} else if !duplicated {
				if e := cat.writer.Write(mdl.Plural, d.name, plural, singular, at); e != nil {
					err = errutil.New("error writing plurals", e)
				}
			}
		}
		return
	})
}
