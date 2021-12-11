package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type ListSource interface {
	Affinity() affine.Affinity
	GetListSource(run rt.Runtime) (ret g.Value, err error)
}

func GetListSource(run rt.Runtime, src ListSource) (ret g.Value, err error) {
	if src == nil {
		err = errutil.New("missing source list")
	} else {
		ret, err = src.GetListSource(run)
	}
	return
}
func (*FromNumList) Affinity() affine.Affinity { return affine.NumList }

func (*FromTxtList) Affinity() affine.Affinity { return affine.TextList }

func (*FromRecList) Affinity() affine.Affinity { return affine.RecordList }

func (op *FromNumList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *FromRecList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *FromTxtList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
