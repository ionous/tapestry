package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func FromRecord(rec *g.Record) rt.Scope {
	return (*recordScope)(rec)
}

type recordScope g.Record

func (rs *recordScope) FieldByName(field string) (g.Value, error) {
	rec := (*g.Record)(rs)
	return rec.GetNamedField(field)
}

func (rs *recordScope) SetFieldByName(field string, val g.Value) error {
	rec := (*g.Record)(rs)
	return rec.SetNamedField(field, val)
}

// todo: example, flag object or db for save.
// for now, simply verify that the field exists.
func (rs *recordScope) SetFieldDirty(field string) (err error) {
	rec := (*g.Record)(rs)
	_, err = rec.GetNamedField(field)
	return
}
