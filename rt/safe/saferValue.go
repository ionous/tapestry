package safe

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

// used when converting values to fields that might require objects.
// if the target field (ex. a pattern local) requires text of a certain type
// and the incoming value is untyped, try to convert it.
func RectifyText(run rt.Runtime, val rt.Value, aff affine.Affinity, cls string) (ret rt.Value, err error) {
	ret = val // provisionally.
	// assigning to a field of typed text (which refers to an object?)
	if aff == affine.Text && len(cls) > 0 {
		// and the input is untyped?
		// (tbd: reject objects of incompatible type?)
		if val.Affinity() == affine.Text && val.Len() > 0 && len(val.Type()) == 0 {
			// is the target field of object type?
			if k, e := run.GetKindByName(cls); e != nil {
				err = e
			} else if k.Implements(kindsOf.Kind.String()) {
				ret, err = run.GetField(meta.ObjectId, val.String())
			}
		}
	}
	return
}

// The opposite of [Truthy].
func Falsy(v rt.Value) (ret bool) {
	return !Truthy(v)
}

// Determine the true/false implication of a value.
// Bool values simply return their value.
// Num values: are true when not exactly zero.
// Text values: are true whenever they contain content.
// List values: are true whenever the list is non-empty.
// ( note this is similar to python, and different than javascript. )
// Record values: are true whenever they have been initialized.
// ( only sub-records start uninitialized; record variables are always true. )
func Truthy(v rt.Value) (ret bool) {
	switch aff := v.Affinity(); aff {
	case affine.Bool:
		ret = v.Bool()

	case affine.Num:
		ret = v.Int() != 0

	case affine.Text:
		ret = v.String() != ""

	case affine.Record:
		ret = v.Record() != nil

	case affine.TextList, affine.NumList, affine.RecordList:
		ret = v.Len() > 0
	}
	return
}

// Attempt to coerce the passed value into the passed affinity.
// Bool and num values can become text ( "true" or "false", or the digits as text. )
// All values can become bool ( according to their truthiness. )
func ConvertValue(run rt.Runtime, val rt.Value, out affine.Affinity) (ret rt.Value, err error) {
	switch aff := val.Affinity(); {
	case aff == out:
		ret = val

	case aff == affine.Bool && out == affine.Text:
		ret = rt.StringOf(strconv.FormatBool(val.Bool()))

	case aff == affine.Num && out == affine.Text:
		ret = rt.StringOf(strconv.FormatFloat(val.Float(), 'g', -1, 64))

	case out == affine.Bool:
		ret = rt.BoolOf(Truthy(val))

	default:
		if e := Check(val, aff); e != nil {
			err = e
		} else {
			ret = val
		}
	}
	return
}
