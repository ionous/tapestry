package list

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

func (op *ListErasing) Execute(run rt.Runtime) (err error) {
	if e := op.erasingIndex(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListErasing) erasingIndex(run rt.Runtime) (err error) {
	if els, e := eraseIndex(run, op.Count, op.Target, op.Index); e != nil {
		err = e
	} else if cnt, otherwise := els.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else {
		run.PushScope(scope.NewSingleValue(op.As, els))
		err = safe.RunAll(run, op.Exe)
		run.PopScope()
	}
	return
}

func (op *ListPopping) Execute(run rt.Runtime) (err error) {
	if e := op.erasingEdge(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListPopping) erasingEdge(run rt.Runtime) (err error) {
	if vs, e := eraseEdge(run, op.Target, op.Edge); e != nil {
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
