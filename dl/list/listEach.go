package list

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListEach) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ListEach) forEach(run rt.Runtime) (err error) {
	if vs, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = e
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		if it := op.As; it == nil {
			err = errutil.New("list iterator was undefined")
		} else if itAff, elAff := it.Affinity(), affine.Element(vs.Affinity()); itAff != elAff {
			err = errutil.New("iterator of %s doesnt support elements of %s", itAff, elAff)
		} else {
			// could cache this -- just trying to keep it simple right now.
			// hopefully could live right in the db.
			const el, index, first, last = 0, 1, 2, 3
			ls := g.NewAnonymousRecord(run, []g.Field{
				{Name: it.Name(), Affinity: itAff, Type: vs.Type()},
				{Name: "index", Affinity: affine.Number},
				{Name: "first", Affinity: affine.Bool},
				{Name: "last", Affinity: affine.Bool},
			})
			run.PushScope(g.RecordOf(ls))
			for i := 0; i < cnt; i++ {
				at := vs.Index(i)
				if e := ls.SetIndexedField(el, at); e != nil {
					err = e
					break
				} else if e := ls.SetIndexedField(index, g.IntOf(i+1)); e != nil {
					err = e
					break
				} else if e := ls.SetIndexedField(first, g.BoolOf(i == 0)); e != nil {
					err = e
					break
				} else if e := ls.SetIndexedField(last, g.BoolOf((i+1) == cnt)); e != nil {
					err = e
					break
				} else if e := op.Do.Execute(run); e != nil {
					var i core.DoInterrupt
					if !errors.As(e, &i) {
						err = e
						break
					} else if !i.KeepGoing {
						break
					}
				}
			}
			run.PopScope()
		}
	}

	return
}
