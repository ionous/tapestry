package scope

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// creates a scope from a record
// requires the runtime in order to create default values for empty fields
func FromRecord(run rt.Runtime, rec *g.Record) rt.Scope {
	if rec == nil {
		panic("null record")
	}
	return recordScope{run, rec}
}

type recordScope struct {
	run rt.Runtime
	rec *g.Record
}

func (rs recordScope) FieldByName(field string) (ret g.Value, err error) {
	run, rec:= rs.run, rs.rec
	if v, e := rec.GetNamedField(field); e != nil {
		err = e
	} else if v.Affinity() != affine.Record || v.Record() != nil {
		ret = v
	} else if k, e := run.GetKindByName(v.Type()); e != nil {
		err = e
	} else {
		newVal := g.RecordOf(k.NewRecord())
		rec.SetNamedField(field, newVal)
		ret = newVal
	}
	return
}

func (rs recordScope) SetFieldByName(field string, val g.Value) error {
	rec:= rs.rec
	return rec.SetNamedField(field, val)
}

// todo: example, flag object or db for save.
// for now, simply verify that the field exists.
func (rs recordScope) SetFieldDirty(field string) (err error) {
	rec:= rs.rec
	_, err = rec.GetNamedField(field)
	return
}
