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
	span rt.NumEval,
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
	} else if start, e := safe.GetOptionalInt(run, atIndex, 0); e != nil {
		err = e // ^ the default of 0 here indicates unspecified
	} else {
		cnt := vs.Len()
		if rng, e := safe.GetOptionalInt(run, span, cnt); e != nil {
			err = e // ^ if the span isn't specified; erase everything.
		} else {
			i := clipStart(start, cnt)  // these turn one-based indices to zero-based
			j := clipRange(i, rng, cnt) // converts range to an ending index
			if i >= 0 && j >= i {       // negative after clip indicate no-op
				err = vs.Splice(i, j, nil, cutList)
			}
		}
	}
	return
}
