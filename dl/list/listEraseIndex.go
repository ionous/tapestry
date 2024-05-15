package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
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
) (ret rt.Value, err error) {
	if at, e := assign.GetReference(run, target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if rub, e := safe.GetOptionalNumber(run, count, 0); e != nil {
		err = e
	} else if startOne, e := safe.GetNumber(run, atIndex); e != nil {
		err = e
	} else {
		start, listLen := startOne.Int(), vs.Len()
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
		if v, e := vs.Splice(start, end, nil); e != nil {
			err = e
		} else {
			ret = v
		}
	}
	return
}
