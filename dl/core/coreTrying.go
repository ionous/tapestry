package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/pattern"

	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// Trying runs a pattern.
// It implements every evaluation,
// erroring if the value requested doesnt support the error returned.
type Trying struct {
	Pattern   pattern.PatternName `if:"selector"`
	Arguments *Arguments          `if:"optional"`
	As        string              // fix: variable definition field
	Do        Activity
	Else      Activity
}

func (*Trying) Compose() composer.Spec {
	return composer.Spec{
		Group: "patterns",
		Desc:  "Runs a pattern, and potentially returns a value.",
		Spec:  "Trying: {pattern%name:pattern_name}{?arguments} as:{as:text} do:{do:activity} else:{else:activity}",
		// Fluent: &composer.Fluid{Name: "trying", Role: composer.Command},
		Stub: true,
	}
}

func (op *Trying) Execute(run rt.Runtime) (err error) {
	if e := op.trying(run, ""); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Trying) trying(run rt.Runtime, aff affine.Affinity) (err error) {
	var args []rt.Arg
	if op.Arguments != nil { // FIX!!!!!!!!
		for _, a := range op.Arguments.Args {
			args = append(args, rt.Arg{a.Name, a.From})
		}
	}
	name := op.Pattern.String()
	if v, e := run.Call(name, aff, args); e == nil {
		ls := g.NewAnonymousRecord(run, []g.Field{
			{Name: op.As, Affinity: v.Affinity(), Type: v.Type()},
		})
		run.PushScope(g.RecordOf(ls))
		err = op.Do.Execute(run)
		run.PopScope()
	} else if errors.Is(e, rt.NoResult{}) {
		err = op.Else.Execute(run)
	} else {
		err = cmdError(op, e)
	}
	return
}
