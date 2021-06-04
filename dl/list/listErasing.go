package list

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/scope"
)

func (op *Erasing) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Erasing) popping(run rt.Runtime) (err error) {
	if els, e := eraseIndex(run, op.Count, op.From, op.AtIndex); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewSingleValue(op.As.String(), els))
		err = op.Do.Execute(run)
		run.PopScope()
	}
	return
}

func (op *ErasingEdge) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ErasingEdge) popping(run rt.Runtime) (err error) {
	if vs, e := eraseEdge(run, op.From, op.AtEdge); e != nil {
		err = e
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		run.PushScope(scope.NewSingleValue(op.As.String(), vs.Index(0)))
		err = op.Do.Execute(run)
		run.PopScope()
	}
	return
}
