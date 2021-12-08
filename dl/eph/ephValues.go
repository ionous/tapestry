package eph

import (
	"errors"

	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteValues(w Writer) error {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		for _, v := range n.values {
			if value, e := cout.Encode(v.value.(jsn.Marshalee), literal.CompactEncoder); e != nil {
				err = errutil.Append(err, e)
				break
			} else if e := w.Write(mdl_val, d.name, n.name, v.field, value, v.at); e != nil {
				err = e
				break
			}
		}
		return
	})
}

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
