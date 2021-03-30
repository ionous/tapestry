package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/pattern"

	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// During determines whether a pattern is running.
type During struct {
	Pattern pattern.PatternName // a text eval here would be like a function pointer maybe..
}

func (*During) Compose() composer.Spec {
	return composer.Spec{
		Group:  "patterns",
		Desc:   "During: Runs a pattern, and potentially returns a value.",
		Fluent: &composer.Fluid{Name: "during", Role: composer.Command},
	}
}

// GetBool returns the first matching bool evaluation.
func (op *During) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if depth, e := op.GetNumber(run); e != nil {
		err = e
	} else {
		depth := depth.Int()
		ret = g.BoolOf(depth > 0)
	}
	return
}

func (op *During) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	name := lang.Underscore(op.Pattern.String())
	if depth, e := run.GetField(object.Running, name); e != nil {
		err = cmdError(op, e)
	} else {
		depth := depth.Int()
		ret = g.IntOf(depth)
	}
	return
}
