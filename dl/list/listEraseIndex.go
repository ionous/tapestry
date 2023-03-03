package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *EraseIndex) Execute(run rt.Runtime) (err error) {
	if _, e := eraseIndex(run, op.Count, op.Target, op.AtIndex); e != nil {
		err = CmdError(op, e)
	}
	return
}

func eraseIndex(run rt.Runtime,
	count rt.NumberEval,
	target assign.Address,
	atIndex rt.NumberEval,
) (ret g.Value, err error) {
	if rub, e := safe.GetOptionalNumber(run, count, 0); e != nil {
		err = e
	} else if root, e := assign.GetRootValue(run, target); e != nil {
		err = e
	} else if els, e := root.GetList(run); e != nil {
		err = e
	} else if startOne, e := safe.GetNumber(run, atIndex); e != nil {
		err = e
	} else {
		start, listLen := startOne.Int(), els.Len()
		if start < 0 {
			start += listLen // wrap negative starts
		} else {
			start -= 1 // adjust to zero based
		}
		var end int
		if start >= listLen {
			start, end = 0, 0 // (still) out of bounds? do nothing.
		} else if rub := rub.Int(); rub <= 0 {
			start, end = 0, 0 // zero and negative removal means remove nothing
		} else {
			// If length + start is less than 0, begin from index 0.
			if start < 0 {
				start = 0
			}
			// too many elements means remove all.
			end = start + rub
			if end > listLen {
				end = listLen
			}
		}
		if v, e := els.Splice(start, end, nil); e != nil {
			err = e
		} else {
			root.SetDirty(run)
			ret = v
		}
	}
	return
}
