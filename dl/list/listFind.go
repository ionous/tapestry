package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListFindBool) GetBool(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

// returns 1 based index
func (op *ListFindBool) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

// fix: with autoconversion of int to bool and one-based indices -
// only the GetNumber variant would be needed.
func (op *ListFindNumber) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(i >= 0)
	}
	return
}

// returns 1 based index
func (op *ListFindNumber) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(i + 1)
	}
	return
}

// zero-based
func (op *ListFindNumber) getIndex(run rt.Runtime) (ret int, err error) {
	if val, e := safe.GetNumber(run, op.Number); e != nil {
		err = e
	} else if vs, e := getList(run, op.InList, affine.Number); e != nil {
		err = e
	} else {
		ret = findFloat(vs, val.Float())
	}
	return
}

func (op *ListFindText) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(i >= 0)
	}
	return
}

// returns 1 based index
func (op *ListFindText) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.getIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(i + 1)
	}
	return
}

func (op *ListFindText) getIndex(run rt.Runtime) (ret int, err error) {
	if val, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if vs, e := getList(run, op.InList, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = findString(vs, val.String())
	}
	return
}

func (op *ListFindList) GetBool(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

// returns 1 based index
func (op *ListFindList) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

func (op *ListFindRecord) GetBool(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

// returns 1 based index
func (op *ListFindRecord) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	panic("not implemented")
}

func getList(run rt.Runtime, eval rt.ListEval, aff affine.Affinity) (ret g.Value, err error) {
	if vs, e := safe.GetList(run, eval); e != nil {
		err = e
	} else if el := affine.Element(vs.Affinity()); len(el) == 0 {
		err = errutil.New("not a list")
	} else if aff != el {
		err = errutil.New("expected", el, "have", aff)
	} else {
		ret = vs
	}
	return
}

func findFloat(vs g.Value, match float64) (ret int) {
	ret = -1
	for i, n := range vs.Floats() {
		if n == match { //epsilon?
			ret = i
			break
		}
	}
	return
}

func findString(vs g.Value, match string) (ret int) {
	ret = -1
	for i, n := range vs.Strings() {
		if n == match {
			ret = i
			break
		}
	}
	return
}
