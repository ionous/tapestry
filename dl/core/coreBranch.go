package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
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
	if e := safe.RunAll(run, op.Does); e != nil {
		err = cmdError(op, e)
	}
	return
}
