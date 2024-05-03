package scope

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// creates a scope from a record
// requires the runtime in order to create default values for empty fields
func FromRecord(run rt.Runtime, rec *rt.Record) rt.Scope {
	if rec == nil {
		panic("null record")
	}
	return recordScope{run, rec}
}

type recordScope struct {
	run rt.Runtime
	rec *rt.Record
}

func (rs recordScope) FieldByName(field string) (ret rt.Value, err error) {
	run, rec := rs.run, rs.rec
	//
	if v, e := rec.GetNamedField(field); e != nil {
		err = e
	} else if v.Affinity() != affine.Record {
		ret = v
	} else if _, ok := v.Record(); ok {
		ret = v
	} else {
		// fix:
		if k, e := run.GetKindByName(v.Type()); e != nil {
			err = e
		} else {
			newVal := rt.RecordOf(rt.NewRecord(k))
			rec.SetNamedField(field, newVal)
			ret = newVal
		}
	}
	return
}

func (rs recordScope) SetFieldByName(field string, val rt.Value) error {
	rec := rs.rec
	return rec.SetNamedField(field, val)
}

// todo: example, flag object or db for save.
// for now, simply verify that the field exists.
func (rs recordScope) SetFieldDirty(field string) (err error) {
	rec := rs.rec
	_, err = rec.GetNamedField(field)
	return
}
