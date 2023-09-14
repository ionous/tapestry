package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *CopyValue) Execute(run rt.Runtime) (err error) {
	if e := op.copyValue(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *CopyValue) copyValue(run rt.Runtime) (err error) {
	if tgt, e := GetReference(run, op.Target); e != nil {
		err = e
	} else if src, e := GetReference(run, op.Source); e != nil {
		err = e
	} else if root, e := src.GetRootValue(run); e != nil {
		err = e
	} else if v, e := root.getValue(run); e != nil {
		err = e
	} else {
		err = tgt.SetValue(run, v)
	}
	return
}
