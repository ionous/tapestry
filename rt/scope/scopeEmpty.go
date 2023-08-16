package scope

import g "git.sr.ht/~ionous/tapestry/rt/generic"

type Empty struct{}

func (k Empty) FieldByName(field string) (g.Value, error) {
	return nil, g.UnknownVariable(field)
}

func (k Empty) SetFieldByName(field string, val g.Value) error {
	return g.UnknownVariable(field)
}

func (k Empty) SetFieldDirty(field string) error {
	return g.UnknownVariable(field)
}
