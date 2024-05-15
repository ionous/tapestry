package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *CopyValue) Execute(run rt.Runtime) (err error) {
	if e := op.copyValue(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *CopyValue) copyValue(run rt.Runtime) (err error) {
	if tgt, e := GetReference(run, op.Target); e != nil {
		err = e
	} else if src, e := GetReference(run, op.Source); e != nil {
		err = e
	} else if v, e := src.GetValue(); e != nil {
		err = e
	} else {
		err = tgt.SetValue(v)
	}
	return
}
