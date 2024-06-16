package list

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// return true if the value exists in the list
func (op *ListFind) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i >= 0) // getIndex() returns a zero-based index.
	}
	return
}

// returns 1 based index of the value in the list
func (op *ListFind) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.IntOf(i + 1) // convert zero to a one-based index.
	}
	return
}

// returns a zero-based index.
func (op *ListFind) getIndex(run rt.Runtime) (ret int, err error) {
	if val, e := safe.GetAssignment(run, op.Value); e != nil {
		err = e
	} else if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if listAff, aff := els.Affinity(), val.Affinity(); aff != affine.Element(listAff) {
		err = fmt.Errorf("%s can't contain %s", listAff, aff)
	} else {
		switch aff {
		case affine.Num:
			ret = findFloat(els, val.Float())
		case affine.Text:
			ret = findString(els, val.String())
		default:
			err = fmt.Errorf("%s not implemented")
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
