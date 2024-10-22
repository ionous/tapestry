package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *ListMap) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListMap) remap(run rt.Runtime) (err error) {
	pat := inflect.Normalize(op.PatternName)
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if src, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(src.Affinity()) {
		err = errors.New("not a list")
	} else {
		var changes int
		aff := affine.Element(vs.Affinity())
		for it := safe.ListIt(src); it.HasNext(); {
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else if newVal, e := run.Call(pat, aff, nil, []rt.Value{inVal}); e != nil {
				// note: this treats "no result" as an error because its
				// trying to map *all* of the elements from one list into another
				err = e
				break
			} else if e := vs.Appends(newVal); e != nil {
				err = e
				break
			} else {
				changes++
			}
		}
	}
	return
}
