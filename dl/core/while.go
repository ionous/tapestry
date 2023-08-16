package core

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// MaxLoopError provides both an error and a counter
type MaxLoopError int

func (e MaxLoopError) Error() string { return "nearly infinite loop detected" }

var MaxLoopIterations MaxLoopError = 0xbad

func (op *While) Execute(run rt.Runtime) (err error) {
	if len(op.Exe) > 0 {
	LoopBreak:
		for i := 0; ; i++ {
			if i >= int(MaxLoopIterations) {
				err = cmdError(op, MaxLoopIterations)
				break
			} else if keepGoing, e := safe.GetBool(run, op.True); e != nil {
				err = cmdError(op, e)
				break
			} else if !keepGoing.Bool() {
				// all done
				break
			} else {
				// run the loop:
				if e := safe.RunAll(run, op.Exe); e != nil {
					var i DoInterrupt
					if !errors.As(e, &i) {
						err = cmdError(op, e)
						break LoopBreak
					} else if !i.KeepGoing {
						break LoopBreak
					}
				}
			}
		}
	}
	return
}
