package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *SetValue) Execute(run rt.Runtime) (err error) {
	if e := op.setValue(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SetValue) setValue(run rt.Runtime) (err error) {
	if newValue, e := GetValue(run, op.Value); e != nil {
		err = e
	} else if ref, e := GetReference(run, op.Target); e != nil {
		err = e
	} else if e := ref.SetValue(run, newValue); e != nil {
		err = e
	}
	return
}
