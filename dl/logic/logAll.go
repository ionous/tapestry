package logic

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Always) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	ret = rt.BoolOf(true)
	return
}

func (op *IsValue) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if val, e := safe.GetAssignment(run, op.Value); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(safe.Truthy(val))
	}
	return
}

// check if all conditions return true.
func (op *IsAll) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first condition to return false.
	if i, cnt, e := resolve(run, op.Test, false); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i > 0 && i == cnt)
	}
	return
}

// check if any conditions return true.
func (op *IsAny) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first condition to return true.
	if i, cnt, e := resolve(run, op.Test, true); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(i < cnt)
	}
	return
}

func resolve(run rt.Runtime, evals []rt.BoolEval, breakOn bool) (i, cnt int, err error) {
	for i, cnt = 0, len(evals); i < cnt; i++ {
		if ok, e := safe.GetBool(run, evals[i]); e != nil {
			err = e
			break
		} else if ok.Bool() == breakOn {
			break
		}
	}
	return i, cnt, err
}
