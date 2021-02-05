package scope

import (
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type TargetRecord struct {
	Target string
	Record *g.Record
}

func (k *TargetRecord) GetField(target, field string) (ret g.Value, err error) {
	if target != k.Target {
		err = g.UnknownTarget{target}
	} else {
		ret, err = k.Record.GetNamedField(field)
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *TargetRecord) SetField(target, field string, v g.Value) (err error) {
	if target != k.Target {
		err = g.UnknownTarget{target}
	} else {
		// FIX: if v is a record, then it becomes shared by k.Record;
		// that's probably incorrect.
		err = k.Record.SetNamedField(field, v)
	}
	return
}
