package core

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *ChooseAction) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseAction) Branch(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseAction) ifDoElse(run rt.Runtime) (err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = e
	} else if b.Bool() {
		err = op.Do.Execute(run)
	} else if branch := op.Else; branch != nil {
		err = branch.Branch(run)
	}
	return
}
