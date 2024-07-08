package jess

import (
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Scheduler interface {
	// can return error if out of phase
	Schedule(weaver.Phase, func(weaver.Weaves, rt.Runtime) error) error
	SchedulePos(compact.Source, weaver.Phase, func(weaver.Weaves, rt.Runtime) error) error
}

type Context struct {
	Query
	Scheduler
	pos compact.Source
}

// overrides the normal find noun to map "you" to the object "self"
// fix: it'd be nice if the mapping of "you" to "self" was handled by script
// ( ex. registering the appropriate names )
func (ctx Context) FindNoun(name []match.TokenValue, pkind *string) (string, int) {
	if len(name) == 1 {
		n := name[0] // a copy
		if n.Token == match.String && n.Hash() == keywords.You {
			n.Value = PlayerSelf         // replace in the copy
			name = []match.TokenValue{n} // leave the original data as is.
		}
	}
	return ctx.Query.FindNoun(name, pkind)
}

// run the passed function now or in the future.
func (ctx Context) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	return ctx.Scheduler.SchedulePos(ctx.pos, when, cb)
}
