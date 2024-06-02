package list

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

func (op *Erasing) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *Erasing) popping(run rt.Runtime) (err error) {
	if els, e := eraseIndex(run, op.Count, op.Target, op.AtIndex); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewSingleValue(op.As, els))
		err = safe.RunAll(run, op.Exe)
		run.PopScope()
	}
	return
}

func (op *ErasingEdge) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ErasingEdge) popping(run rt.Runtime) (err error) {
	if vs, e := eraseEdge(run, op.Target, op.AtEdge); e != nil {
		err = e
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		run.PushScope(scope.NewSingleValue(op.As, vs.Index(0)))
		err = safe.RunAll(run, op.Exe)
		run.PopScope()
	}
	return
}
