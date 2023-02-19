package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func (op *ListMap) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListMap) remap(run rt.Runtime) (err error) {
	pat := op.UsingPattern
	if src, e := op.List.GetList(run); e != nil {
		err = e
	} else if root, e := op.Target.GetRootValue(run); e != nil {
		err = e
	} else if tgt, e := root.GetList(run); e != nil {
		err = e
	} else {
		const (
			inArg = iota
		)
		aff := affine.Element(tgt.Affinity())
		for it := g.ListIt(src); it.HasNext() && err == nil; {
			if inVal, e := it.GetNext(); e != nil {
				err = e
			} else if rec, e := core.MakeRecord(run, pat); e != nil {
				err = e // created a fresh record so it has blank default values
			} else if e := rec.SetIndexedField(inArg, inVal); e != nil {
				err = e
			} else {
				if newVal, e := run.Call(rec, aff); e != nil {
					// note: we treat no result as an error because
					// we are trying to map *all* of the elements from one list into another
					err = e
				} else {
					err = tgt.Appends(newVal)
				}
			}
		}
	}
	return
}
