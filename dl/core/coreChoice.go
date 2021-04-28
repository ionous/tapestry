package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// ChooseAction runs one block of instructions or another based on the results of a conditional evaluation.
// If:do:, If:do:else:
type ChooseAction struct {
	If   rt.BoolEval `if:"selector,placeholder=choose"`
	Do   Activity
	Else Brancher `if:"selector,optional"`
}

func (*ChooseAction) Compose() composer.Spec {
	return composer.Spec{
		Lede:   "if",
		Fluent: &composer.Fluid{Name: "if", Role: composer.Command},
		Desc:   "Choose action: an if statement.",
	}
}

func (op *ChooseAction) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseAction) Branch(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseAction) ifDoElse(run rt.Runtime) (err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = e
	} else if b.Bool() {
		err = op.Do.Execute(run)
	} else if branch := op.Else; branch != nil {
		err = branch.Branch(run)
	}
	return
}
