package core

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *Activity) Empty() bool {
	return len(op.Exe) == 0
}

// Execute statements
func (op *Activity) Execute(run rt.Runtime) (err error) {
	if e := safe.RunAll(run, op.Exe); e != nil {
		err = cmdError(op, e)
	}
	return
}
