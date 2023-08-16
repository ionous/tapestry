package qna

import (
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

// helper for tracking the start and end of patterns
type callState struct {
	run       *Runner
	name      string
	prevScope scope.Chain
}

func (run *Runner) saveCallState() callState {
	return callState{
		run:       run,
		prevScope: run.replaceScope(scope.Empty{}),
	}
}

func (r *callState) restore() {
	r.run.restoreScope(r.prevScope)
	r.setPattern("")
}

func (r *callState) setPattern(name string) {
	if r.name != "" {
		r.run.currentPatterns.stoppedPattern(r.name)
		r.name = ""
	}
	if name != "" {
		r.run.currentPatterns.startedPattern(name)
		r.name = name
	}
}
