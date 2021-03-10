package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"github.com/ionous/errutil"
)

// Trying runs a pattern.
// It implements every evaluation,
// erroring if the value requested doesnt support the error returned.
type Trying struct {
	Name      pattern.PatternName `if:"selector"`
	Arguments *Arguments          `if:"optional"`
	As        string              // fix: variable definition field; fix: should be optional.
	Filter    rt.BoolEval         `if:"optional,selector=and"`
	Do        Activity
	Else      Brancher `if:"selector,optional"`
}

func (*Trying) Compose() composer.Spec {
	return composer.Spec{
		Group: "patterns",
		Desc:  "Trying: Runs a pattern, and potentially returns a value.",
		// Spec:  "Trying: {pattern%name:pattern_name}{?arguments} as:{as:text} {filter?bool_eval} do:{do:activity} {else:activity}",
		Fluent: &composer.Fluid{Name: "trying", Role: composer.Command},
		Stub:   true,
	}
}

func (op *Trying) Execute(run rt.Runtime) (err error) {
	if e := op.trying(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Trying) trying(run rt.Runtime) (err error) {
	var args []rt.Arg
	if op.Arguments != nil { // FIX!!!!!!!!
		for _, a := range op.Arguments.Args {
			args = append(args, rt.Arg{a.Name, a.From})
		}
	}
	name := op.Name.String()
	const anyaff affine.Affinity = "" // pass empty string for any affinity
	if v, e := run.Call(name, anyaff, args); e == nil {
		if hasReturn, hasLocal := v != nil, len(op.As) > 0; hasReturn != hasLocal {
			if hasReturn {
				err = errutil.New("expected a local value")
			} else {
				err = errutil.New("expected a return value")
			}
		} else {
			run.PushScope(scope.NewSingleValue(op.As, v))
			if ok, e := safe.GetOptionalBool(run, op.Filter, true); e != nil {
				err = e
			} else {
				if ok.Bool() {
					err = op.Do.Execute(run)
				} else if branch := op.Else; branch != nil {
					err = branch.Branch(run)
				}
			}
			run.PopScope()
		}
	} else if errors.Is(e, rt.NoResult{}) {
		if branch := op.Else; branch != nil {
			err = branch.Branch(run)
		}
	} else {
		err = e
	}
	return
}
