package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
}

// ChooseMore provides an else-if clause.
// Like ChooseAction it chooses a branch based on an if statement
// with an appropriate lede for the composer.
type ChooseMore ChooseAction

// ChooseMoreValue provides an else-if assignment clause.
// Like ChooseValue it chooses a branch based on an if statement and a new value.
type ChooseMoreValue ChooseValue

// ChooseNothingElse provides an else clause.
type ChooseNothingElse struct {
	Do Activity `if:"selector"`
}

func (*ChooseMore) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "elseIf", Role: composer.Selector},
	}
}

func (*ChooseMoreValue) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "elseIf", Role: composer.Selector},
	}
}

func (*ChooseNothingElse) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "elseDo", Role: composer.Selector},
	}
}

func (op *ChooseMore) Branch(run rt.Runtime) (err error) {
	if e := (*ChooseAction)(op).ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseMoreValue) Branch(run rt.Runtime) (err error) {
	if e := (*ChooseValue)(op).ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseNothingElse) Branch(run rt.Runtime) (err error) {
	if e := op.Do.Execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}
