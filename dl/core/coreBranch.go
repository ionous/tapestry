package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
	// used for determine the state of pattern guards during weave
	// returns the else branch or nil;
	// returns true if a continuing statement is possible ( always true when the else branch exists)
	// or false if it the statement cannot be continued ( ex. ending in an pure "else" )
	Descend() (next Brancher, canContinue bool)
}

func (op *ChooseMore) Descend() (Brancher, bool) {
	return op.Else, true
}

func (op *ChooseMore) Branch(run rt.Runtime) (err error) {
	if e := (*ChooseAction)(op).ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseMoreValue) Descend() (Brancher, bool) {
	return op.Else, true
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

func (op *ChooseNothingElse) Descend() (Brancher, bool) {
	return nil, false
}
