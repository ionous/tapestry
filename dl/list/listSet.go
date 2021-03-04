package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Set struct {
	List  string // variable name
	Index rt.NumberEval
	From  rt.Assignment
}

func (*Set) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_set",
		Group: "list",
		Desc:  "Set value in list: Overwrite an existing value in a list.",
	}
}

func (op *Set) Execute(run rt.Runtime) (err error) {
	if e := op.setAt(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Set) setAt(run rt.Runtime) (err error) {
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
