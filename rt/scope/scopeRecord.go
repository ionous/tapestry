package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// creates a scope from a record
// requires kinds so that unit tests can create sub-records safely
func FromRecord(ks rt.Kinds, rec *rt.Record) rt.Scope {
	return recordScope{ks, rec, make(map[string]bool)}
}

type recordScope struct {
	ks      rt.Kinds
	rec     *rt.Record
	changed map[string]bool
}

func (rs recordScope) FieldChanged(field string) bool {
	return rs.changed[field]
}

func (rs recordScope) FieldByName(field string) (rt.Value, error) {
	return safe.EnsureField(rs.ks, rs.rec, field)
}

func (rs recordScope) SetFieldByName(field string, val rt.Value) error {
	rs.changed[field] = true
	return rs.rec.SetNamedField(field, val)
}
