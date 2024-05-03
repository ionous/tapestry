package rt

import (
	"git.sr.ht/~ionous/tapestry/affine"
)

func BoolOf(v bool) Value {
	return BoolFrom(v, defaultType)
}
func StringOf(v string) Value {
	return StringFrom(v, defaultType)
}
func FloatOf(v float64) Value {
	return FloatFrom(v, defaultType)
}
func IntOf(v int) Value {
	return IntFrom(v, defaultType)
}
func RecordOf(v *Record) Value {
	return makeValue(affine.Record, v.Name(), v)
}
func StringsOf(vs []string) Value {
	return StringsFrom(vs, defaultType)
}
func FloatsOf(vs []float64) Value {
	return FloatsFrom(vs, defaultType)
}

func BoolFrom(v bool, subtype string) Value {
	return makeValue(affine.Bool, subtype, v)
}
func StringFrom(v string, subtype string) Value {
	return makeValue(affine.Text, subtype, v)
}
func FloatFrom(v float64, subtype string) Value {
	return makeValue(affine.Number, subtype, v)
}
func IntFrom(v int, subtype string) Value {
	return makeValue(affine.Number, subtype, v)
}

// returns a nil record of the specified type
func RecordFrom(subtype string) Value {
	var n *Record
	return makeValue(affine.Record, subtype, n)
}

func StringsFrom(vs []string, subtype string) (ret Value) {
	if a := affine.TextList; vs != nil {
		ret = makeValue(a, subtype, &vs)
	} else {
		ret = makeValue(a, subtype, new([]string))
	}
	return
}

func FloatsFrom(vs []float64, subtype string) (ret Value) {
	if a := affine.NumList; vs != nil {
		// note: this address is of the unique slice "vs" which shares memory with the slice passed
		// but has its own length. quite possibly this should be marked as "read-only"
		ret = makeValue(a, subtype, &vs)
	} else {
		ret = makeValue(a, subtype, new([]float64))
	}
	return
}

func RecordsFrom(vs []*Record, subtype string) (ret Value) {
	if a := affine.RecordList; vs != nil {
		ret = makeValue(a, subtype, &vs)
	} else {
		ret = makeValue(a, subtype, new([]*Record))
	}
	return
}

func makeValue(a affine.Affinity, subtype string, i any) variant {
	return variant{a: a, i: i, t: subtype}
}
