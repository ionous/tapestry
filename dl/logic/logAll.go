package logic

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Never) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	ret = rt.BoolOf(false)
	return
}

func (op *Always) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	ret = rt.BoolOf(true)
	return
}

func (op *AllTrue) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first statement to return false.
	if i, cnt, e := resolve(run, op.Test, false); e != nil {
		err = cmd.Error(op, e)
	} else if i < cnt {
		ret = rt.False
	} else {
		ret = rt.True // return true, resolve never found a false statement
	}
	return
}

func (op *AnyTrue) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	// stop on the first statement to return true.
	if i, cnt, e := resolve(run, op.Test, true); e != nil {
		err = cmd.Error(op, e)
	} else if i < cnt {
		ret = rt.True
	} else {
		ret = rt.False // return false, resolve never found a true statement
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
