package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/scope"
)

// ChooseValue creates a local assignment for use with evaluation.
// if:txt:and:do:else: #name "value" {#name=="value"}  {...} {...}
type ChooseValue struct {
	Assign string
	From   rt.Assignment `if:"selector"`
	Filter rt.BoolEval   `if:"pb=and,selector=and"`
	Do     Activity
	Else   Brancher `if:"selector,optional"`
}

func (*ChooseValue) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "if", Role: composer.Command},
		Desc:   "Choose value: an if statement with local assignment.",
	}
}

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
	if v, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewSingleValue(op.Assign, v))
		//
		if ok, e := safe.GetOptionalBool(run, op.Filter, true); e != nil {
			err = e
		} else {
			if ok.Bool() {
				err = op.Do.Execute(run)
			} else if branch := op.Else; branch != nil {
				err = branch.Branch(run)
			}
		}
		run.PopScope()
	}
	return
}
