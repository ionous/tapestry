package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListEmpty) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = cmd.Error(op, e)
	} else if !affine.IsList(els.Affinity()) {
		err = cmd.Error(op, errors.New("not a list"))
	} else {
		ret = rt.BoolOf(els.Len() == 0)
	}
	return
}

func (op *ListLength) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = cmd.Error(op, e)
	} else if !affine.IsList(els.Affinity()) {
		err = cmd.Error(op, errors.New("not a list"))
	} else {
		ret = rt.IntOf(els.Len())
	}
	return
}
