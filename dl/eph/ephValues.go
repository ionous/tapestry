package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

func (c *Catalog) WriteValues(w Writer) error {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		for _, v := range n.values {
			if value, e := encodeLiteral(v.value); e != nil {
				err = errutil.Append(err, e)
			} else if e := w.Write(mdl_value, d.name, n.name, v.field, value, v.at); e != nil {
				err = e
				break
			}
		}
		return
	})
}

// name of a noun to assembly info
func (op *EphValues) Phase() Phase { return ValuePhase }

// Noun  string           `if:"label=noun,type=text"`
// Field string           `if:"label=has,type=text"`
// Value literal.LiteralValue `if:"label=value"`
func (op *EphValues) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if noun, e := getClosestNoun(d, op.Noun); e != nil {
		err = e
	} else if field, ok := UniformString(op.Field); !ok {
		err = InvalidString(op.Field)
	} else if value := op.Value; value == nil {
		err = errutil.New("null value", op.Noun, op.Field)
	} else {
		if e := noun.AddLiteralValue(field, value, at); e != nil {
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
