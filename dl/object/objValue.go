package object

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SetValue) Execute(run rt.Runtime) (err error) {
	if e := op.setValue(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *SetValue) setValue(run rt.Runtime) (err error) {
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e // note: things are easier to debug if this grabs the target first
	} else if newValue, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if e := at.SetValue(newValue); e != nil {
		err = e
	}
	return
}

func (op *SetState) Execute(run rt.Runtime) (err error) {
	if e := op.setState(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *SetState) setState(run rt.Runtime) (err error) {
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = e
	} else if at, e := at.Dot(dot.Field(trait.String())); e != nil {
		err = e
	} else {
		err = at.SetValue(rt.BoolOf(true))
	}
	return
}
