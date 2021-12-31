package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type ListTarget interface {
	GetListTarget(run rt.Runtime) (g.Value, error)
}

func (op *IntoNumList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), affine.NumList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoRecList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoTxtList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), affine.TextList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
