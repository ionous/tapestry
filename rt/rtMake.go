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
	return makeVariant(affine.Record, v.Name(), v)
}
func StringsOf(vs []string) Value {
	return StringsFrom(vs, defaultType)
}
func FloatsOf(vs []float64) Value {
	return FloatsFrom(vs, defaultType)
}

func BoolFrom(v bool, subtype string) Value {
	return makeVariant(affine.Bool, subtype, v)
}
func StringFrom(v string, subtype string) Value {
	return makeVariant(affine.Text, subtype, v)
}
func FloatFrom(v float64, subtype string) Value {
	return makeVariant(affine.Num, subtype, v)
}
func IntFrom(v int, subtype string) Value {
	return makeVariant(affine.Num, subtype, v)
}

func StringsFrom(vs []string, subtype string) (ret Value) {
	if a := affine.TextList; vs != nil {
		ret = makeVariant(a, subtype, &vs)
	} else {
		ret = makeVariant(a, subtype, new([]string))
	}
	return
}

func FloatsFrom(vs []float64, subtype string) (ret Value) {
	if a := affine.NumList; vs != nil {
		// note: this address is of the unique slice "vs" which shares memory with the slice passed
		// but has its own length. quite possibly this should be marked as "read-only"
		ret = makeVariant(a, subtype, &vs)
	} else {
		ret = makeVariant(a, subtype, new([]float64))
	}
	return
}

func RecordsFrom(vs []*Record, subtype string) (ret Value) {
	if a := affine.RecordList; vs != nil {
		ret = makeVariant(a, subtype, &vs)
	} else {
		ret = makeVariant(a, subtype, new([]*Record))
	}
	return
}

func makeVariant(a affine.Affinity, subtype string, i any) variant {
	return variant{a: a, i: i, t: subtype}
}
