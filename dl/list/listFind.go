package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// return true if the value exists in the list
func (op *ListFind) GetBool(run rt.Runtime) (ret g.Value, err error) {
	// fix: with autoconversion of int to bool and one-based indices -
	// only the GetNumber variant would be needed.
	if i, e := op.getIndex(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = g.BoolOf(i >= 0)
	}
	return
}

// returns 1 based index of the value in the list
func (op *ListFind) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = g.IntOf(i + 1)
	}
	return
}

// zero-based
func (op *ListFind) getIndex(run rt.Runtime) (ret int, err error) {
	if val, e := op.Value.GetValue(run); e != nil {
		err = e
	} else if els, e := op.List.GetList(run); e != nil {
		err = e
	} else if listAff, aff := els.Affinity(), val.Affinity(); aff != affine.Element(listAff) {
		err = errutil.New(listAff, "can't contain", aff)
	} else {
		switch aff {
		case affine.Number:
			ret = findFloat(els, val.Float())
		case affine.Text:
			ret = findString(els, val.String())
		default:
			err = errutil.New(aff, "not implemented")
		}
	}
	return
}

func findFloat(els g.Value, match float64) (ret int) {
	ret = -1
	for i, n := range els.Floats() {
		if n == match { //epsilon?
			ret = i
			break
		}
	}
	return
}

func findString(els g.Value, match string) (ret int) {
	ret = -1
	for i, n := range els.Strings() {
		if n == match {
			ret = i
			break
		}
	}
	return
}
