package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// Kinds database ( primarily for generating default values )
type Kinds interface {
	KindByName(name string) (*Kind, error)
}

// DefaultFor generates a zero value for the specified affinity;
// uses the passed Kinds to generate empty records when necessary.
func DefaultFor(ks Kinds, a affine.Affinity, subtype string) (ret Value, err error) {
	// return the default value for the
	switch a {
	case affine.Bool:
		ret = False

	case affine.Number:
		ret = Zero
	case affine.NumList:
		ret = ZeroList

	case affine.Text:
		ret = Empty
	case affine.TextList:
		ret = EmptyList

	case affine.Record:
		if n, e := ks.KindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret = n.NewRecord()
		}
	case affine.RecordList:
		if n, e := ks.KindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret = n.NewRecordSlice()
		}
	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}
