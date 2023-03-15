package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

func (op *ChooseValue) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseValue) Branch(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseValue) ifDoElse(run rt.Runtime) (err error) {
	if v, e := safe.GetAssignment(run, op.From); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewSingleValue(op.Assign, v))
		//
		if ok, e := safe.GetOptionalBool(run, op.Filter, true); e != nil {
			err = e
		} else {
			if ok.Bool() {
				err = safe.RunAll(run, op.Does)
			} else if branch := op.Else; branch != nil {
				err = branch.Branch(run)
			}
		}
		run.PopScope()
	}
	return
}
