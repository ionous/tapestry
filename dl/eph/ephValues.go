package eph

import (
	"errors"

	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteValues(w Writer) error {
	return forEachNoun(c, func(k *ScopedKind, n *ScopedNoun) (err error) {
		for _, v := range n.values {
			// FIX: nouns should be able to store EVALS too
			// example: an object with a counter in its description.
			if value, e := marshalout(v.value); e != nil {
				err = errutil.Append(err, e)
			} else if e := w.Write(mdl.Value, n.domain.name, n.name, v.field, value, v.at); e != nil {
				err = e
				break
			}
		}
		return
	})
}

// name of a noun to assembly info
func (op *EphValues) Phase() Phase { return ValuePhase }

func (op *EphValues) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if noun, e := getClosestNoun(d, op.Noun); e != nil {
		err = e
	} else if field, ok := UniformString(op.Field); !ok {
		err = InvalidString(op.Field)
	} else if value := op.Value; value == nil {
		err = errutil.New("null value", op.Noun, op.Field)
	} else {
		// fix? should values be written per domain instead of per noun?
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
