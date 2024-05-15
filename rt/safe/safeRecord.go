package safe

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
)

// helper for trying to create a record using a name from the passed typeset.
func NewRecord(ks rt.Kinds, cls string) (ret *rt.Record, err error) {
	if k, e := ks.GetKindByName(cls); e != nil {
		err = e
	} else {
		ret = rt.NewRecord(k)
	}
	return
}

// pulls a field out of this record, creating a sub-record if necessary
func EnsureField(ks rt.Kinds, rec *rt.Record, field string) (ret rt.Value, err error) {
	if v, e := rec.GetNamedField(field); e == nil {
		ret = v
	} else {
		var z rt.NilRecord
		if !errors.As(e, &z) {
			err = e
		} else {
			// getting the field failed because points to a nil record;
			// so create and store that record. somebody has to.
			if subRec, e := NewRecord(ks, z.Class); e != nil {
				err = e
			} else {
				newVal := rt.RecordOf(subRec)
				if e := rec.SetIndexedField(z.Field, newVal); e != nil {
					err = e
				} else {
					ret = newVal
				}
			}
		}
	}
	return
}
