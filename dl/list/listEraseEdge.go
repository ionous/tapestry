package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

/*
 * erase: numEval
 * from: varName,
 * atIndex: num,
 */
type EraseEdge struct {
	From   ListSource `if:"selector"`
	AtEdge Edge       `if:"selector"`
}

func (*EraseEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase at edge: Remove one or more values from a list",
	}
}

func (op *EraseEdge) Execute(run rt.Runtime) (err error) {
	if _, e := op.pop(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *EraseEdge) pop(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := GetListSource(run, op.From); e != nil {
		err = e
	} else {
		if cnt := vs.Len(); cnt > 0 {
			var at int
			if !op.AtEdge.Front() {
				at = cnt - 1
			}
			ret, err = vs.Splice(at, at+1, nil)
		}
	}
	return
}
