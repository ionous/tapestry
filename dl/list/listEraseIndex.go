package list

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListErase) Execute(run rt.Runtime) (err error) {
	if _, e := eraseIndex(run, op.Count, op.Target, op.Index); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func eraseIndex(run rt.Runtime,
	count rt.NumEval,
	target rt.Address,
	atIndex rt.NumEval,
) (ret rt.Value, err error) {
	if at, e := safe.GetReference(run, target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if rub, e := safe.GetOptionalNumber(run, count, 1); e != nil {
		err = e
	} else if startOne, e := safe.GetOptionalNumber(run, atIndex, 1); e != nil {
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
