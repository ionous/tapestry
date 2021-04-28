package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/scope"
)

/**
 * removing: count
 * numFrom/txt/rec:varName,
 * atIndex:
 * as: string, do:{}
 */
type Erasing struct {
	EraseIndex
	As string // fix: variable definition field
	Do core.Activity
}

func (op *Erasing) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erasing", Role: composer.Command},
		Group:  "list",
		Desc: `Erasing indices: Erase elements from the front or back of a list.
Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.`,
	}
}

func (op *Erasing) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Erasing) popping(run rt.Runtime) (err error) {
	if els, e := op.pop(run); e != nil {
		err = e
	} else {
		run.PushScope(scope.NewSingleValue(op.As, els))
		err = op.Do.Execute(run)
		run.PopScope()
	}
	return
}

/**
 * removing: count
 * numFrom/txt/rec:varName,
 * atIndex:
 * as: string, do:{}
 */
type ErasingEdge struct {
	EraseEdge
	As   string // fix: variable definition field
	Do   core.Activity
	Else core.Brancher `if:"selector,optional"`
}

func (op *ErasingEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erasing", Role: composer.Command},
		Group:  "list",
		Desc: `Erasing at edge: Erase one element from the front or back of a list.
Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.`,
	}
}

func (op *ErasingEdge) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ErasingEdge) popping(run rt.Runtime) (err error) {
	if vs, e := op.pop(run); e != nil {
		err = e
	} else if cnt, otherwise := vs.Len(), op.Else; otherwise != nil && cnt == 0 {
		err = otherwise.Branch(run)
	} else if cnt > 0 {
		run.PushScope(scope.NewSingleValue(op.As, vs.Index(0)))
		err = op.Do.Execute(run)
		run.PopScope()
	}
	return
}
