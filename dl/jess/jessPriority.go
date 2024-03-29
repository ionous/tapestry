package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Scheduler interface {
	// can return error if out of phase
	Schedule(phase weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) error
}

type Context struct {
	Query
	Scheduler
}

func NewContext(q Query, n Scheduler) Context {
	return Context{q, n}
}
