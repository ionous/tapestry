package qna

import (
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

// helper for tracking the start and end of patterns
type callState struct {
	run       *Runner
	name      string
	prevScope scope.Chain
}

func (run *Runner) saveCallState(top *g.Record) callState {
	state := callState{
		run:       run,
		prevScope: run.scope.ReplaceScope(scope.FromRecord(top)), // scope.Empty{}
	}
	state.setPattern(top.Kind().Name())
	return state
}

func (r *callState) restore() {
	r.run.scope.RestoreScope(r.prevScope)
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
