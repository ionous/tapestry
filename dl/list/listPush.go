package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
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
	if ins, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if root, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if els, e := root.GetList(run); e != nil {
		err = e
	} else if !IsAppendable(ins, els) {
		err = insertError{ins, els}
	} else {
		if atFront, e := safe.GetOptionalBool(run, op.AtEdge, false); e != nil {
			err = e
		} else {
			if !atFront.Bool() {
				err = els.Appends(ins)
			} else {
				_, err = els.Splice(0, 0, ins)
			}
			if err == nil {
				root.SetDirty(run)
			}
		}
	}
	return
}
