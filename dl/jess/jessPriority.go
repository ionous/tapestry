package jess

import (
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
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

// run the passed function now or in the future.
func (ctx Context) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	return ctx.Scheduler.SchedulePos(ctx.pos, when, cb)
}
