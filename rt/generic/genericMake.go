package generic

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
	return RecordFrom(v, v.Type())
}
func RecordsOf(typeName string, vs []*Record) Value {
	return RecordsFrom(vs, typeName)
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
func RecordFrom(v *Record, subtype string) Value {
	return makeValue(affine.Record, subtype, v)
}

// changes a named object into just the name
// the returned value has affinity text, with a class of the object's kind.
func ObjectAsText(v Value) (ret Value) {
	if v == nil {
		ret = StringOf("") // fix: or "nothing"?
	} else {
		id, typeName := v.String(), v.Type()
		ret = makeValue(affine.Text, typeName, id)
	}
	return
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

func makeValue(a affine.Affinity, subtype string, i interface{}) (ret refValue) {
	// fix? we arent writing these to values in the db
	// so having them here feels a bit odd
	// -- especially because the subtype of text here becomes "string" and that's not a kind....
	// (ex. see matchTypes)
	// if len(subtype) == 0 {
	// 	t := r.TypeOf(i)
	// 	if t.Kind() == r.Ptr {
	// 		t = t.Elem()
	// 	}
	// 	if t.Kind() == r.Slice {
	// 		t = t.Elem()
	// 	}
	// 	subtype = t.String()
	// }
	return refValue{a: a, i: i, t: subtype}
}
