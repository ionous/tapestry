package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// add the passed set of variables to the pool of current variables
// ( implements rt.Runtime )
func (run *Runner) PushScope(top rt.Scope) {
	run.scope.PushScope(top)
}

// remove most recently pushed set of variables.
// ( implements rt.Runtime )
func (run *Runner) PopScope() {
	run.scope.PopScope()
}
