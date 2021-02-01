package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// Activity wraps a block of multiple execute statements
// This is primarily useful for the composer so it can display blocks in a uniform manner.
type Activity struct {
	Exe []rt.Execute
}

func (*Activity) Compose() composer.Spec {
	return composer.Spec{
		Group: "hidden",
		Spec:  "{exe*execute}",
	}
}

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
