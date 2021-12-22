package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *HasTrait) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.False
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else {
		trait := trait.String()
		if p, e := obj.FieldByName(trait); e == nil {
			ret = p
		} else if errors.Is(e, g.UnknownField(obj.String(), trait)) {
			ret = g.False
		} else {
			err = cmdError(op, e)
		}
	}
	return
}

func (op *SetTrait) Execute(run rt.Runtime) (err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if obj != nil {
		if e := obj.SetFieldByName(trait.String(), g.BoolOf(true)); e != nil {
			err = cmdError(op, e)
		}
	}
	return
}
