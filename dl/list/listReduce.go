package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

func (op *ListReduce) Execute(run rt.Runtime) (err error) {
	if e := op.reduce(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListReduce) reduce(run rt.Runtime) (err error) {
	pat := op.PatternName
	if tgt, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if fromList, e := assign.GetValue(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(fromList.Affinity()) {
		err = errutil.New("not a list")
	} else {
		const (
			inArg = iota
			outArg
		)
		outVal := tgt.RootValue
		for it := g.ListIt(fromList); it.HasNext() && err == nil; {
			if inVal, e := it.GetNext(); e != nil {
				err = e
			} else if rec, e := assign.MakeRecord(run, pat); e != nil {
				err = e // created a fresh record so it has blank default values
			} else if e := rec.SetIndexedField(inArg, inVal); e != nil {
				err = e
			} else if e := rec.SetIndexedField(outArg, outVal); e != nil {
				err = e
			} else {
				outAff := rec.Kind().Field(outArg).Affinity
				if newVal, e := run.Call(rec, outAff); e == nil {
					// update the accumulating value for next time
					outVal = newVal
				} else if !errors.Is(e, rt.NoResult{}) {
					// if there was no result, just keep going with what we had
					// for other errors, break.
					err = e
				}
			}
		}
		// did we have a successful result at some point?
		if err == nil && outVal != tgt.RootValue {
			err = tgt.SetValue(run, outVal)
		}
	}
	return
}
