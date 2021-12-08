package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

// name of a noun to assembly info
func (el *EphValues) Phase() Phase { return ValuePhase }

// Noun  string           `if:"label=noun,type=text"`
// Field string           `if:"label=has,type=text"`
// Value literal.LiteralValue `if:"label=value"`
func (el *EphValues) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if noun, e := getClosestNoun(d, el.Noun); e != nil {
		err = e
	} else if field, ok := UniformString(el.Field); !ok {
		err = InvalidString(el.Field)
	} else if value := el.Value; value == nil {
		err = errutil.New("null value", el.Noun, el.Field)
	} else {
		if e = noun.AddLiteralValue(field, value, at); e != nil {
			var conflict *Conflict
			if errors.As(e, &conflict) && conflict.Reason == Duplicated {
				LogWarning(e)
			} else {
				err = e
			}
		}
	}
	return
}
