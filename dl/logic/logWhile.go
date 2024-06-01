package logic

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

// MaxLoopError provides both an error and a counter
type MaxLoopError int

func (e MaxLoopError) Error() string { return "nearly infinite loop detected" }

var MaxLoopIterations MaxLoopError = 0xbad

func (op *Repeat) Execute(run rt.Runtime) (err error) {
	var pop bool
	if ks, vs, e := assign.ExpandArgs(run, op.Initial); e != nil {
		err = cmd.Error(op, e)
	} else {
		if pop = len(vs) > 0; pop {
			run.PushScope(scope.NewPairs(ks, vs))
		}
		if e := op.loop(run); e != nil {
			err = cmd.Error(op, e)
		}
		if pop {
			run.PopScope()
		}
	}
	return
}

func (op *Repeat) loop(run rt.Runtime) (err error) {
	keepGoing := len(op.Exe) > 0
	for i := 0; keepGoing && err == nil; i++ {
		if i >= int(MaxLoopIterations) {
			err = MaxLoopIterations
		} else {
			if ks, vs, e := assign.ExpandArgs(run, op.Args); e != nil {
				err = e
			} else {
				var pop bool
				if pop = len(vs) > 0; pop {
					run.PushScope(scope.NewPairs(ks, vs))
				}
				keepGoing, err = runOnce(run, op.Condition, op.Exe)
				if pop {
					run.PopScope()
				}
			}
		}
	}
	return
}

// return true to keep going
func runOnce(run rt.Runtime, c rt.BoolEval, exe []rt.Execute) (okay bool, err error) {
	if keepGoing, e := safe.GetBool(run, c); e != nil {
		err = e
	} else if keepGoing.Bool() {
		if e := safe.RunAll(run, exe); e == nil {
			okay = true
		} else {
			var i DoInterrupt
			if !errors.As(e, &i) {
				err = e
			} else {
				okay = i.KeepGoing
			}
		}
	}
	return
}
