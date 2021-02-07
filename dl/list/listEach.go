package list

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type Each struct {
	List core.Assignment `if:"selector=across"`
	As   ListIterator    `if:"selector"`
	Do   core.Activity
	Else *ElseIfEmpty `if:"optional,selector"`
}

type ElseIfEmpty struct {
	Do core.Activity `if:"selector"`
}

func (op *Each) Compose() composer.Spec {
	return composer.Spec{
		// fix: once the fluent interface is decided, rename the command, and remove the name.
		// alt: the compact format could use the fluent names and then the actual command name doesnt matter much
		Name:   "list_each",
		Group:  "list",
		Fluent: &composer.Fluid{Name: "repeating", Role: composer.Command},
		Desc:   `For each in list: Loops over the elements in the passed list, or runs the 'else' activity if empty.`,
		Locals: []string{"index", "first", "last"},
	}
}

func (op *ElseIfEmpty) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_empty_do",
		Group:  "list",
		Fluent: &composer.Fluid{Name: "elseIfEmptyDo", Role: composer.Selector},
		Desc:   `Runs an activity when a list is empty.`,
	}
}

func (op *Each) Execute(run rt.Runtime) (err error) {
	if e := op.forEach(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Each) forEach(run rt.Runtime) (err error) {
	if vs, e := core.GetAssignedValue(run, op.List); e != nil {
		err = e
	} else {
		if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
			err = op.Else.Do.Execute(run)
		} else if cnt > 0 {
			//
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
	}
	return
}
