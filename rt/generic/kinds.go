package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

// Kinds database ( primarily for generating default values )
// tbd: a better way of handling this might be to pass this into GetIndexedField... etc.
type Kinds interface {
	GetKindByName(n string) (*Kind, error)
}

// NewDefaultValue generates a zero value for the specified affinity;
// uses the passed Kinds to generate empty records when necessary.
func NewDefaultValue(ks Kinds, aff affine.Affinity, cls string) (ret Value, err error) {
	// return the default value for the
	switch aff {
	case affine.Bool:
		ret = BoolFrom(false, cls)

	case affine.Number:
		ret = FloatFrom(0, cls)

	case affine.NumList:
		ret = FloatsFrom(nil, cls)

	case affine.Text: //fix!
		var defaultValue string
		if len(cls) > 0 {
			if ak, e := ks.GetKindByName(cls); e == nil {
				if Base(ak) == kindsOf.Aspect.String() {
					trait := ak.Field(0) // get the default trait.
					defaultValue = trait.Name
				}
			}
		}
		ret = StringFrom(defaultValue, cls)

	case affine.TextList:
		ret = StringsFrom(nil, cls)

	case affine.Record:
		if k, e := ks.GetKindByName(cls); e != nil {
			err = errutil.New("unknown kind", cls, e)
		} else {
			ret = RecordOf(k.NewRecord())
		}

	case affine.RecordList:
		if _, e := ks.GetKindByName(cls); e != nil {
			err = errutil.New("unknown kind", cls, e)
		} else {
			ret = RecordsFrom(nil, cls)
		}

	default:
		err = errutil.New("default value requested for unhandled affinity", aff)
	}
	return
}
