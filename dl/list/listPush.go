package list

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListPush) Execute(run rt.Runtime) (err error) {
	if e := op.push(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListPush) push(run rt.Runtime) (err error) {
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if ins, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if !IsAppendable(ins, vs) {
		err = insertError{ins, vs}
	} else {
		if atFront, e := safe.GetOptionalBool(run, op.AtEdge, false); e != nil {
			err = e
		} else {
			if !atFront.Bool() {
				err = vs.Appends(ins)
			} else {
				_, err = vs.Splice(0, 0, ins)
			}
		}
	}
	return
}
