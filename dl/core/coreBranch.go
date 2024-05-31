package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"github.com/ionous/errutil"
)

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
}

// else statements can only be run in the context of previous decisions
// ( include things like the branch of an empty while statements )
func (op *ChooseNothingElse) Branch(run rt.Runtime) (err error) {
	// not having special local variables for the else block
	// makes "pick()" simpler; it also happens to match go-lang.
	if e := safe.RunAll(run, op.Exe); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseBranch) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseBranch) Branch(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseBranch) ifDoElse(run rt.Runtime) (err error) {
	var pushes int
	if exe, e := op.PickBranch(run, &pushes); e != nil {
		err = e
	} else {
		err = safe.RunAll(run, exe)
	}
	safe.PopSeveral(run, pushes)
	return
}

// loop through descendant branches to find which block we should run.
// rationale: patterns sometimes need to update counters and *not* run the block which guard them.
// this makes rule processing and normal branch selection the same.
func (op *ChooseBranch) PickBranch(run rt.Runtime, pushes *int) (ret []rt.Execute, err error) {
Pick:
	for next := op; next != nil; {
		if matched, e := next.eval(run, pushes); e != nil {
			err = e
			break
		} else {
			if matched {
				// see if there's another branch to evaluate;
				// if not, return the statements to execute
				if n := PickTree(next.Exe); n != nil {
					next = n
				} else {
					ret = next.Exe
					break
				}
			} else {
				switch otherwise := next.Else.(type) {
				case nil:
					break Pick
				case *ChooseBranch:
					next = otherwise
				case *ChooseNothingElse:
					if n := PickTree(otherwise.Exe); n != nil {
						next = n
					} else {
						ret = otherwise.Exe
						break Pick
					}
				default:
					err = errutil.Fmt("unknown branch %T", otherwise)
					break Pick
				}
			}
		}
	}
	return
}

func (op *ChooseBranch) eval(run rt.Runtime, pushed *int) (matched bool, err error) {
	if ks, vs, e := assign.ExpandArgs(run, op.Args); e != nil {
		err = e
	} else {
		if len(op.Args) > 0 {
			run.PushScope(scope.NewPairs(ks, vs))
			*pushed++
		}
		if b, e := safe.GetBool(run, op.Condition); e != nil {
			err = e
		} else {
			matched = b.Bool()
		}
	}
	return
}

// scan for an initial branching statement in a "case like" switch,.
// that statement indicates a filter for a rule;
// returns nil if no such statement exists
func PickTree(exe []rt.Execute) (ret *ChooseBranch) {
FindTree:
	for i, el := range exe {
		switch el := el.(type) {
		case *debug.DebugLog:
			// skip debug logs when trying to find a tree
			// FIX: run those logs?
			// if so, we have to slice them out
			// ( ex. run exe[0:i] once the branch is chosen )
			// or could return an index for ChooseBranch i suppose, -1 on nothing found
			// and have the caller run the excess.
		case *ChooseBranch:
			// found a branch; we'll use it if we can.
			// if there's a following sibling statement of any kind --
			// other than a debug log, but even another branch --
			// then we're not really a switch like series of if-else(s)
			for _, sib := range exe[i+1:] {
				if _, ok := sib.(*debug.DebugLog); !ok {
					break FindTree
				}
			}
			ret = el
		default:
			// some other statement
			break FindTree
		}
	}
	return
}
