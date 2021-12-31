package list

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *PutEdge) Execute(run rt.Runtime) (err error) {
	if e := op.push(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *PutEdge) push(run rt.Runtime) (err error) {
	if ins, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = e
	} else if els, e := op.Into.GetListTarget(run); e != nil {
		err = e
	} else if !IsAppendable(ins, els) {
		err = insertError{ins, els}
	} else {
		if atFront, e := safe.GetOptionalBool(run, op.AtEdge, false); e != nil {
			err = e
		} else if !atFront.Bool() {
			err = els.Appends(ins)
		} else {
			_, err = els.Splice(0, 0, ins)
		}
	}
	return
}
