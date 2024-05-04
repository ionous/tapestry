package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// creates a scope from a record
func FromRecord(rec *rt.Record) rt.Scope {
	return recordScope{rec, make(map[string]bool)}
}

type recordScope struct {
	rec     *rt.Record
	changed map[string]bool
}

func (rs recordScope) FieldChanged(field string) bool {
	return rs.changed[field]
}

func (rs recordScope) FieldByName(field string) (rt.Value, error) {
	return rs.rec.GetNamedField(field)
}

func (rs recordScope) SetFieldByName(field string, val rt.Value) error {
	rs.changed[field] = true
	return rs.rec.SetNamedField(field, val)
}
