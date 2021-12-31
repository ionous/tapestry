package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

type counters map[string]int

func (c *counters) getCounter(name string) (ret g.Value, err error) {
	// fix: i think at some point we should have a global $counters object
	// with named fields for each counter; that would let save/load work normally.
	i := (*c)[name]
	ret = g.IntOf(i)
	return
}

func (c *counters) setCounter(name string, val g.Value) (err error) {
	if aff := val.Affinity(); aff != affine.Number {
		err = errutil.Fmt("counter %q expected a number got %s", name, aff)
	} else {
		(*c)[name] = val.Int()
	}
	return
}
