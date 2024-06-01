package logic

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Not) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(!val.Bool())
	}
	return
}

func (op *Never) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	ret = rt.BoolOf(false)
	return
}

func (op *NotValue) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if val, e := safe.GetAssignment(run, op.Value); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(safe.Falsy(val))
	}
	return
}

// check if all conditions return false.
func (op *NotAll) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first condition to return true.
	if i, cnt, e := resolve(run, op.Test, true); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i > 0 && i == cnt)
	}
	return
}

// check if any conditions return false.
func (op *NotAny) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first condition to return false.
	if i, cnt, e := resolve(run, op.Test, false); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i < cnt)
	}
	return
}
