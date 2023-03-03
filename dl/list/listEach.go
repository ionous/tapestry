package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListEach) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListEach) forEach(run rt.Runtime) (err error) {
	if vs, e := assign.GetSafeAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(vs.Affinity()) {
		err = errutil.New("not a list")
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		if it := op.As; len(it) == 0 {
			// fix: a more robust check? an assembly check?
			err = errutil.New("list iterator was undefined")
		} else {
			// could cache this -- just trying to keep it simple right now.
			// maybe this type could live in the db.
			const el, index, first, last = 0, 1, 2, 3
			itAff, itType := affine.Element(vs.Affinity()), vs.Type()
			ls := g.NewAnonymousRecord(run, []g.Field{
				{Name: it, Affinity: itAff, Type: itType},
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
				} else if e := safe.RunAll(run, op.Does); e != nil {
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
