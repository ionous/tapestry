package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

func (op *ListReduce) Execute(run rt.Runtime) (err error) {
	if e := op.reduce(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListReduce) reduce(run rt.Runtime) (err error) {
	pat := inflect.Normalize(op.PatternName)
	if at, e := assign.GetReference(run, op.Target); e != nil {
		err = e
	} else if accum, e := at.GetValue(); e != nil {
		err = e
	} else if fromList, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(fromList.Affinity()) {
		err = errutil.New("not a list")
	} else {
		changed := false
		for it := g.ListIt(fromList); it.HasNext() && err == nil; {
			if inVal, e := it.GetNext(); e != nil {
				err = e
			} else {
				if newVal, e := run.Call(pat, accum.Affinity(), nil, []g.Value{inVal, accum}); e == nil {
					// update the accumulating value for next time
					accum = newVal
					changed = true
				} else if !errors.Is(e, rt.NoResult) {
					// if there was no result, just keep going with what we had
					// for other errors, break.
					err = e
				}
			}
		}
		// did we have a successful result at some point?
		if err == nil && changed {
			err = at.SetValue(accum)
		}
	}
	return
}
