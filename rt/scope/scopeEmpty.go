package scope

import "git.sr.ht/~ionous/tapestry/rt"

type Empty struct{}

func (Empty) FieldByName(field string) (rt.Value, error) {
	return nil, rt.UnknownVariable(field)
}

func (Empty) SetFieldByName(field string, val rt.Value) error {
	return rt.UnknownVariable(field)
}

func (Empty) FieldChanged(field string) bool {
	return false
}
