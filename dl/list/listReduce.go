package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *ListReduce) Execute(run rt.Runtime) (err error) {
	if e := op.reduce(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListReduce) reduce(run rt.Runtime) (err error) {
	pat := inflect.Normalize(op.PatternName)
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e
	} else if accum, e := at.GetValue(); e != nil {
		err = e
	} else if fromList, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(fromList.Affinity()) {
		err = errors.New("not a list")
	} else {
		for it := safe.ListIt(fromList); it.HasNext(); {
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if newVal, e := run.Call(pat, accum.Affinity(), nil, []rt.Value{inVal, accum}); e != nil {
				err = e
				break
			} else {
				accum = newVal // update the value for next loop
			}
		}
		if err == nil {
			err = at.SetValue(accum)
		}
	}
	return
}
