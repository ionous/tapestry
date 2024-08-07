package jess

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type JessContext struct {
	Query
	Scheduler
}

type ParallelMatcher interface {
	ParallelMatcher() ParallelMatcher
	typeinfo.Instance
}

func After(p weaver.Phase) weaver.Phase {
	return p + 1
}

func (jc JessContext) Try(z weaver.Phase, cb func(weaver.Weaves, rt.Runtime), reject func(error)) {
	if e := jc.Schedule(z, func(w weaver.Weaves, run rt.Runtime) (_ error) {
		cb(w, run)
		return
	}); e != nil {
		reject(e)
	}
}
