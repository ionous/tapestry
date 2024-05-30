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

func Truthy(v rt.Value) (ret bool) {
	switch aff := v.Affinity(); aff {
	case affine.Bool:
		ret = v.Bool()

	case affine.Num:
		ret = v.Int() > 0

	case affine.Text:
		ret = v.String() != ""

	case affine.Record:
		ret = true

	case affine.TextList, affine.NumList, affine.RecordList:
		ret = v.Len() > 0
	}
	return
}

func ConvertValue(run rt.Runtime, val rt.Value, out affine.Affinity) (ret rt.Value, err error) {
	switch aff := val.Affinity(); {
	case aff == out:
		ret = val

	case out == affine.Text && aff == affine.Bool:
		ret = rt.StringOf(strconv.FormatBool(val.Bool()))

	case out == affine.Text && aff == affine.Num:
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
