package list

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *ListSet) Execute(run rt.Runtime) (err error) {
	if e := op.setAt(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ListSet) setAt(run rt.Runtime) (err error) {
	if els, e := safe.List(run, op.List); e != nil {
		err = e
	} else if onedex, e := safe.GetNumber(run, op.Index); e != nil {
		err = e
	} else if el, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = e
	} else if !IsInsertable(el, els) {
		err = insertError{el, els}
	} else if i, e := safe.Range(onedex.Int()-1, 0, els.Len()); e != nil {
		err = e
	} else {
		els.SetIndex(i, el)
	}
	return
}
