package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"github.com/ionous/errutil"
)

// CopyValue: create a new value from a snapshot of the passed value
// panics on error because it assumes all values should be copyable
func CopyValue(val Value) (ret Value) {
	if val == nil {
		panic("failed to copy nil value")
	}
	switch a := val.Affinity(); a {
	// because we dont have a 'value.set()', the values of primitives are immutable.
	// therefore we dont have to copy them, which saves us from having to mange their subtypes.
	// ( ex. copy of an int, should still probably be an int under the hood. )
	case affine.Bool, affine.Number, affine.Text:
		ret = val

	case affine.Record:
		ret = copyRecord(val)

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

// duplicate the passed slice of floats
// ( b/c golang's built in copy doesnt allocate )
func copyFloats(src []float64) []float64 {
	out := make([]float64, len(src))
	copy(out, src)
	return out
}

// duplicate the passed slice of strings.
func copyStrings(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

// duplicates all of the passed records
// panics on error because it assumes all records are copyable.
func copyRecords(src []*Record) []*Record {
	out := make([]*Record, len(src))
	for i, el := range src {
		out[i] = copyRecordValues(el)
	}
	return out
}

func copyRecord(v Value) (ret Value) {
	if rec, ok := v.Record(); !ok {
		ret = v
	} else {
		out := copyRecordValues(rec)
		ret = RecordOf(out)
	}
	return
}

// panics on error because it assumes all records are copyable.
func copyRecordValues(src *Record) (ret *Record) {
	if src != nil {
		values := make([]Value, len(src.values))
		for i, v := range src.values {
			// if a value hasn't been accessed, it might still be at its default value
			// skip copying it, the new record will also have the default value for that field
			if v != nil {
				values[i] = CopyValue(v)
			}
		}
		ret = &Record{kind: src.kind, values: values}
	}
	return
}
