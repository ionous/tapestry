package core

import (
	rtx "git.sr.ht/~ionous/iffy/idl/rtx"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

func (op *Always) GetBool(run rt.Runtime) (ret g.Value, err error) {
	ret = g.BoolOf(true)
	return
}

func (op *AllTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	// stop on the first statement to return false.
	if i, cnt, e := resolve(run, op.Test, false); e != nil {
		err = cmdError(op, op.Struct, e)
	} else if i < cnt {
		ret = g.False
	} else {
		ret = g.True // return true, resolve never found a false statement
	}
	return
}

func (op *AnyTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	// stop on the first statement to return true.
	if i, cnt, e := resolve(run, op.Test, true); e != nil {
		err = cmdError(op, op.Struct, e)
	} else if i < cnt {
		ret = g.True
	} else {
		ret = g.False // return false, resolve never found a true statement
	}
	return
}

func resolve(run rt.Runtime, fn func() (rtx.BoolEval_List, error), breakOn bool) (i, cnt int, err error) {
	if evals, e := fn(); e != nil {
		err = e
	} else {
		for i, cnt = 0, evals.Len(); i < cnt; i++ {
			op := evals.At(i)
			if ok, e := op.GetBool(run); e != nil {
				err = e
				break
			} else if ok.Bool() == breakOn {
				break
			}
		}
	}
	return
}
