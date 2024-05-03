package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SetValue) Execute(run rt.Runtime) (err error) {
	if e := op.setValue(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SetValue) setValue(run rt.Runtime) (err error) {
	if ref, e := GetReference(run, op.Target); e != nil {
		err = e // note: things are easier to debug if this grabs the target frst
	} else if newValue, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if e := ref.SetValue(newValue); e != nil {
		err = e
	}
	return
}

func (op *SetTrait) Execute(run rt.Runtime) (err error) {
	if e := op.setTrait(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SetTrait) setTrait(run rt.Runtime) (err error) {
	if tgt, e := safe.GetText(run, op.Target); e != nil {
		err = e
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = e
	} else if obj, e := run.GetField(meta.ObjectId, tgt.String()); e != nil {
		err = e
	} else {
		err = run.SetField(obj.String(), trait.String(), rt.BoolOf(true))
	}
	return
}
