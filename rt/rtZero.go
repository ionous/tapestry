package rt

import (
	"errors"
	"fmt"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
)

var (
	True  = BoolOf(true)
	False = BoolOf(false)
	Zero  = FloatOf(0.0)
	Empty = StringOf("")
)

const defaultType = "" // empty string

// returned when trying to generate the zero value of a field containing a record
// record fields are nil by default so that types are able to refer to themselves
// in go, types can use pointers for self references ( ex. struct LinkedList { Link *LinkedList } )
type NilRecord struct {
	Class string // name of the record type that the field wants
	Field int    // index of the field in its parent record
}

func IsNilRecord(e error) bool {
	var z NilRecord
	return errors.As(e, &z)
}

func (e NilRecord) NoPanic() {}

func (e NilRecord) Error() string {
	return fmt.Sprintf("%q contained a nil record (at %d)", e.Class, e.Field)
}

// ZeroValue generates a zero value for the specified affinity;
// asking for a zero record returns an error.
func ZeroValue(aff affine.Affinity) (ret Value, err error) {
	return ZeroField(aff, "", -1)
}

// ZeroField generates a zero value for the specified data ( a field )
// asking for a zero record returns a NilRecord error.
func ZeroField(aff affine.Affinity, cls string, idx int) (ret Value, err error) {
	switch aff {
	case affine.Bool:
		ret = BoolFrom(false, cls)

	case affine.Number:
		ret = FloatFrom(0, cls)

	case affine.NumList:
		ret = FloatsFrom(nil, cls)

	case affine.Text:
		ret = StringFrom("", cls)

	case affine.TextList:
		ret = StringsFrom(nil, cls)

	case affine.Record:
		if idx >= 0 && len(cls) > 0 {
			err = NilRecord{cls, idx}
		} else {
			err = errors.New("invalid record")
		}

	case affine.RecordList:
		ret = RecordsFrom(nil, cls)

	default:
		log.Panicf("unhandled affinity %s", aff)
	}
	return
}
