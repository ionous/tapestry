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

func (rs *recordScope) SetFieldDirty(field string) (err error) {
	rec := (*g.Record)(rs)
	_, err = rec.GetNamedField(field)
	return
}
