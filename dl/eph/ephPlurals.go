package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
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

func (op *EphPlurals) Phase() assert.Phase { return assert.PluralPhase }

func (op *EphPlurals) Weave(k assert.Assertions) (err error) {
	return k.AssertPlural(op.Singular, op.Plural)
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
//
// tbd: consider appending the origin (at) to store the location of each definition?
// alt: `on conflict (domain, many) where @one == one do nothing` ( or do set )
func (op *EphPlurals) Assemble(ctx *Context) (err error) {
	c, d, at := ctx.c, ctx.d, ctx.at
	if plural, ok := UniformString(op.Plural); !ok {
		err = InvalidString(op.Plural)
	} else if singular, ok := UniformString(op.Singular); !ok {
		err = InvalidString(op.Singular)
	} else {
		duplicated := false
		if e := c.qx.FindPluralDefinitions(plural, func(domain, one, at string) (err error) {
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
			if e := c.writer.Write(mdl.Plural, d.name, plural, singular, at); e != nil {
				err = errutil.New("error writing plurals", e)
			}
		}
	}
	return
}
