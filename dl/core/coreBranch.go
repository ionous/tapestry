package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
}

func (op *ChooseNothingElse) Branch(run rt.Runtime) (err error) {
	if e := op.doElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseNothingElse) doElse(run rt.Runtime) (err error) {
	if ks, vs, e := assign.ExpandArgs(run, op.Args); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewPairs(ks, vs))
		err = safe.RunAll(run, op.Exe)
		run.PopScope()
	}
	return
}

func (op *ChooseBranch) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseBranch) Branch(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseBranch) ifDoElse(run rt.Runtime) (err error) {
	if ks, vs, e := assign.ExpandArgs(run, op.Args); e != nil {
		err = e
	} else {
		var pushPop bool
		if pushPop = len(op.Args) > 0; pushPop {
			run.PushScope(scope.NewPairs(ks, vs))
		}
		if b, e := safe.GetBool(run, op.If); e != nil {
			err = e
		} else if b.Bool() {
			err = safe.RunAll(run, op.Exe)
		} else if branch := op.Else; branch != nil {
			err = branch.Branch(run)
		}
		if pushPop {
			run.PopScope()
		}
	}
	return
}
