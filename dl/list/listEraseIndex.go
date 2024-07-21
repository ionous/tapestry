package list

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListErase) Execute(run rt.Runtime) (err error) {
	if e := eraseIndex(run, op.Count, op.Target, op.Start, nil); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func eraseIndex(run rt.Runtime,
	count rt.NumEval,
	target rt.Address,
	atIndex rt.NumEval,
	cutList *rt.Value,
) (err error) {
	if at, e := safe.GetReference(run, target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if start, e := safe.GetOptionalInt(run, atIndex, 1); e != nil {
		err = e // start is a one based offset here.
	} else {
		var end int
		size := vs.Len() // if count isnt specified; erase everything
		if count, e := safe.GetOptionalInt(run, count, size); e != nil {
			err = e
		} else {
			if count <= 0 {
				start, end = 0, 0 // zero and negative count means remove nothing
			} else {
				if start < 0 {
					start += size // wrap negative starts
				} else {
					start -= 1 // adjust to zero based
				}
				if start >= size {
					start, end = 0, 0 // (still) out of bounds? do nothing.
				} else {
					// If length + start is less than 0, begin from index 0.
					if start < 0 {
						start = 0
					}
					// clip too many elements
					end = start + count
					if end > size {
						end = size
					}
				}
			}
			// always splice unless there was a critical error.
			err = vs.Splice(start, end, nil, cutList)
		}
	}
	return
}
