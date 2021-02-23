package rel

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Relate struct {
	Object   rt.TextEval `if:"selector"`
	ToObject rt.TextEval `if:"selector=to"`
	Via      Relation
}

func (*Relate) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "relate", Role: composer.Command},
		Group:  "relations",
		Desc:   "Relate: Relate two nouns.",
	}
}

func (op *Relate) Execute(run rt.Runtime) (err error) {
	if e := op.setRelation(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Relate) setRelation(run rt.Runtime) (err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = e
	} else if b, e := safe.ObjectText(run, op.ToObject); e != nil {
		err = e
	} else {
		err = run.RelateTo(a.String(), b.String(), op.Via.String())
	}
	return
}
