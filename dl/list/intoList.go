package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type ListTarget interface {
	GetListTarget(run rt.Runtime) (g.Value, error)
}

func (op *IntoNumList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoRecList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoTxtList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
