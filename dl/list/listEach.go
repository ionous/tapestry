package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt/scope"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListRepeat) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListRepeat) forEach(run rt.Runtime) (err error) {
	if vs, e := safe.GetAssignment(run, op.List); e != nil {
		err = e
	} else if !affine.IsList(vs.Affinity()) {
		err = errors.New("not a list")
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		if it := op.As; len(it) == 0 {
			// fix: a more robust check? an assembly check?
			err = errors.New("list iterator was undefined")
		} else {
			// could cache this -- just trying to keep it simple right now.
			// maybe this type could live in the db.
			const el, index, first, last = 0, 1, 2, 3
			itAff, itType := affine.Element(vs.Affinity()), vs.Type()
			ls := rt.NewRecord(&rt.Kind{Fields: []rt.Field{
				{Name: it, Affinity: itAff, Type: itType},
				{Name: "index", Affinity: affine.Num},
				{Name: "first", Affinity: affine.Bool},
				{Name: "last", Affinity: affine.Bool},
			}})
			run.PushScope(scope.FromRecord(run, ls))
			for i := 0; i < cnt; i++ {
				at := vs.Index(i)
				if e := ls.SetIndexedField(el, at); e != nil {
					err = e // ^ note: the record index doesnt copy the value.
					break
				} else if e := ls.SetIndexedField(index, rt.IntOf(i+1)); e != nil {
					err = e
					break
				} else if e := ls.SetIndexedField(first, rt.BoolOf(i == 0)); e != nil {
					err = e
					break
				} else if e := ls.SetIndexedField(last, rt.BoolOf((i+1) == cnt)); e != nil {
					err = e
					break
				} else if e := safe.RunAll(run, op.Exe); e != nil {
					var i logic.DoInterrupt
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
