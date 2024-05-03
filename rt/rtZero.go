package rt

var (
	True  = BoolOf(true)
	False = BoolOf(false)
	Zero  = FloatOf(0.0)
	Empty = StringOf("")
)

const defaultType = "" // empty string

// NewDefaultValue generates a zero value for the specified affinity
// Record values (and lists) are nil.
func NewDefaultValue(aff affine.Affinity, cls string) (ret Value, err error) {
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
		err = errutil.New("default value requested for unhandled affinity", aff)
	}
	return
}
