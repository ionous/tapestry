package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/scope"
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

// rewrite the current scope as the passed scope
// returns the previous scope
// ( used for function calls; they "hide" all variables from the caller stack )
func (run *Runner) replaceScope(top rt.Scope) (ret scope.Chain) {
	ret, run.scope = run.scope, scope.Chain{Scope: top}
	return ret
}

func (run *Runner) restoreScope(oldScope scope.Chain) {
	run.scope = oldScope
}
