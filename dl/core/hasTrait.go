package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *HasTrait) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else {
		trait := trait.String()
		if p, e := run.GetField(obj, trait); e == nil {
			ret = p
		} else if errors.Is(e, g.UnknownField(obj, trait)) {
			ret = g.False
		} else {
			err = cmdError(op, e)
		}
	}
	return
}

func (op *SetTrait) Execute(run rt.Runtime) (err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if e := run.SetField(obj.String(), trait.String(), g.BoolOf(true)); e != nil {
		err = cmdError(op, e)
	}
	return
}
