package scope

import g "git.sr.ht/~ionous/iffy/rt/generic"

type SingleValue struct {
	name string
	val  g.Value
}

func NewSingleValue(name string, val g.Value) *SingleValue {
	return &SingleValue{name, val}
}

func (k *SingleValue) SetValue(val g.Value) {
	k.val = val
}

func (k *SingleValue) FieldByName(field string) (ret g.Value, err error) {
	if found := len(field) > 0 && field == k.name; !found {
		err = g.UnknownVariable(field)
	} else {
		ret = k.val
	}
	return
}

func (k *SingleValue) SetFieldByName(field string, val g.Value) (err error) {
	if found := len(field) > 0 && field == k.name; !found {
		err = g.UnknownVariable(field)
	} else {
		k.val = val
	}
	return
}
