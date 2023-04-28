package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// note: values are written per *noun* not per domain....
func (ctx *Context) AssertNounValue(opNoun, opField string, opPath []string, opValue literal.LiteralValue) (err error) {
	d, at := ctx.d, ctx.at
	if noun, e := getClosestNoun(d, opNoun); e != nil {
		err = e
	} else if rv, e := noun.recordValues(at); e != nil {
		err = e
	} else if field, ok := UniformString(opField); !ok {
		err = InvalidString(opField)
	} else if path, e := UniformStrings(opPath); e != nil {
		err = e
	} else if value := opValue; value == nil {
		err = errutil.New("null value", opNoun, opField)
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
