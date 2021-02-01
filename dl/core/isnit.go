package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// IsNotTrue returns the opposite of a boolean eval.
type IsNotTrue struct {
	Test rt.BoolEval `if:"selector"`
}

func (*IsNotTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:   "not",
		Fluent: &composer.Fluid{Name: "not", Role: composer.Function},
		Group:  "logic",
		Desc:   "Is Not: Returns the opposite value.",
	}
}

func (op *IsNotTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(!val.Bool())
	}
	return
}
