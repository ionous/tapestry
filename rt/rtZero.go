package rt

import (
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

// ZeroValue generates a zero value for the specified affinity
// Record values (and lists) are nil.
func ZeroValue(aff affine.Affinity, cls string) (ret Value, err error) {
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
		ret = RecordFrom(cls)

	case affine.RecordList:
		ret = RecordsFrom(nil, cls)

	default:
		log.Panicf("unhandled affinity %s", aff)
	}
	return
}
