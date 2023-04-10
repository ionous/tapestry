package eph

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteValues(w Writer) error {
	// FIX: nouns should be able to store EVALS too
	// example: an object with a counter in its description.
	return forEachNoun(c, func(n *ScopedNoun) (err error) {
		if rv := n.localRecord; rv.isValid() {
			for _, fv := range rv.rec.Fields {
				if value, e := marshalout(fv.Value); e != nil {
					err = errutil.Append(err, e)
				} else if e := w.Write(mdl.Value, n.domain.name, n.name, fv.Field, value, rv.at); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

// name of a noun to assembly info
func (op *EphValues) Phase() assert.Phase { return assert.ValuePhase }

// note: values are written per *noun* not per domain....
func (op *EphValues) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if noun, e := getClosestNoun(d, op.Noun); e != nil {
		err = e
	} else if rv, e := noun.recordValues(at); e != nil {
		err = e
	} else if field, ok := UniformString(op.Field); !ok {
		err = InvalidString(op.Field)
	} else if path, e := UniformStrings(op.Path); e != nil {
		err = e
	} else if value := op.Value; value == nil {
		err = errutil.New("null value", op.Noun, op.Field)
	} else {
		var conflict *Conflict
		if e := rv.writeValue(noun.name, at, field, path, value); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e)
		} else {
			err = e // might be nil
		}
	}
	return
}
