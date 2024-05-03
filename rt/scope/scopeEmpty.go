package scope

import "git.sr.ht/~ionous/tapestry/rt"

type Empty struct{}

func (k Empty) FieldByName(field string) (rt.Value, error) {
	return nil, rt.UnknownVariable(field)
}

func (k Empty) SetFieldByName(field string, val rt.Value) error {
	return rt.UnknownVariable(field)
}

func (k Empty) SetFieldDirty(field string) error {
	return rt.UnknownVariable(field)
}
