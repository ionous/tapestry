package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListMap) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListMap) remap(run rt.Runtime) (err error) {
	pat := op.PatternName
	if src, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(src.Affinity()) {
		err = errutil.New("not a list")
	} else if root, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if tgt, e := root.GetList(run); e != nil {
		err = e
	} else {
		const (
			inArg = iota
		)
		var changes int
		aff := affine.Element(tgt.Affinity())
		for it := g.ListIt(src); it.HasNext() && err == nil; {
			if inVal, e := it.GetNext(); e != nil {
				err = e
			} else if rec, e := assign.MakeRecord(run, pat); e != nil {
				err = e // created a fresh record so it has blank default values
			} else if e := rec.SetIndexedField(inArg, inVal); e != nil {
				err = e
			} else if newVal, e := run.Call(rec, aff); e != nil {
				// note: this treats "no result" as an error because its
				// trying to map *all* of the elements from one list into another
				err = e
			} else if e := tgt.Appends(newVal); e != nil {
				err = e
			} else {
				changes++
			}
		}
		if err == nil && changes > 0 {
			// Appends doesn't inform the caller of a result; so we have to.
			root.SetDirty(run)
		}
	}
	return
}
