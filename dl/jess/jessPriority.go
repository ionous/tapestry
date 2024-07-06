package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Scheduler interface {
	// can return error if out of phase
	Schedule(weaver.Phase, func(weaver.Weaves, rt.Runtime) error) error
	SchedulePos(weaver.Phase, mdl.Source, func(weaver.Weaves, rt.Runtime) error) error
}

type Context struct {
	Query
	Scheduler
	pos mdl.Source
}

// run the passed function now or in the future.
func (ctx Context) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	return ctx.Scheduler.SchedulePos(when, ctx.pos, cb)
}
