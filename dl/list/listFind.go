package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// return true if the value exists in the list
func (op *ListFind) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// fix: with autoconversion of int to bool and one-based indices -
	// only the GetNum variant would be needed.
	if i, e := op.getIndex(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i >= 0)
	}
	return
}

// returns 1 based index of the value in the list
func (op *ListFind) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.IntOf(i + 1)
	}
	return
}

// zero-based
func (op *ListFind) getIndex(run rt.Runtime) (ret int, err error) {
	if val, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if listAff, aff := els.Affinity(), val.Affinity(); aff != affine.Element(listAff) {
		err = errutil.New(listAff, "can't contain", aff)
	} else {
		switch aff {
		case affine.Num:
			ret = findFloat(els, val.Float())
		case affine.Text:
			ret = findString(els, val.String())
		default:
			err = errutil.New(aff, "not implemented")
		}
	}
	return
}

func findFloat(els rt.Value, match float64) (ret int) {
	ret = -1
	for i, n := range els.Floats() {
		if n == match { //epsilon?
			ret = i
			break
		}
	}
	return
}

func findString(els rt.Value, match string) (ret int) {
	ret = -1
	for i, n := range els.Strings() {
		if n == match {
			ret = i
			break
		}
	}
	return
}
