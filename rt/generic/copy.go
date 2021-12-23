package generic

import (
	"git.sr.ht/~ionous/iffy/affine"
	"github.com/ionous/errutil"
)

// CopyValue: create a new value from a snapshot of the passed value
// panics on error because it assumes all values should be copyable
func CopyValue(val Value) (ret Value) {
	if val == nil {
		panic("failed to copy nil value")
	}
	switch a := val.Affinity(); a {
	// because we dont have a value.set the values of primitives are immutable
	// so we dont have to actually copy them, which saves us from having to mange their subtypes
	// ( ex. copy of an int, should still probably be an int under the hood. )
	case affine.Bool, affine.Number, affine.Text:
		ret = val

	case affine.Record:
		vs := copyRecord(val.Record())
		ret = RecordOf(vs)

	case affine.NumList:
		vs := copyFloats(val.Floats())
		ret = FloatsFrom(vs, val.Type())

	case affine.TextList:
		vs := copyStrings(val.Strings())
		ret = StringsFrom(vs, val.Type())

	case affine.RecordList:
		vs := copyRecords(val.Records())
		ret = RecordsFrom(vs, val.Type())

	default:
		panic(errutil.Sprint("failed to dupe value of %s:%v(%T)", a, val, val))
	}
	return
}

// duplicates all of the passed records
// panics on error because it assumes all records are copyable.
func copyRecords(src []*Record) []*Record {
	out := make([]*Record, len(src))
	for i, el := range src {
		cpy := copyRecord(el)
		out[i] = cpy
	}
	return out
}

// assumes in value is a record.
// panics on error because it assumes all records are copyable.
func copyRecord(v *Record) (ret *Record) {
	cnt := v.kind.NumField()
	values := make([]Value, cnt)
	for i := 0; i < cnt; i++ {
		if el, e := v.GetIndexedField(i); e != nil {
			panic(e)
		} else {
			values[i] = CopyValue(el)
		}
	}
	return &Record{kind: v.kind, values: values}
}
