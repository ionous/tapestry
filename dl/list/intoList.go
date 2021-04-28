package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type ListTarget interface {
	GetListTarget(run rt.Runtime) (g.Value, error)
}

type IntoNumList struct {
	Var core.Variable `if:"selector"`
}
type IntoTxtList struct {
	Var core.Variable `if:"selector"`
}
type IntoRecList struct {
	Var core.Variable `if:"selector"`
}

func (*IntoNumList) Compose() composer.Spec {
	return composer.Spec{
		Lede:   "nums",
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of numbers",
	}
}

func (*IntoTxtList) Compose() composer.Spec {
	return composer.Spec{
		Lede:   "txts",
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of text",
	}
}

func (*IntoRecList) Compose() composer.Spec {
	return composer.Spec{
		Lede:   "recs",
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of records",
	}
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
