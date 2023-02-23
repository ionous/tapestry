package core

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// GetBool returns the first matching bool evaluation.
func (op *During) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if depth, e := op.GetNumber(run); e != nil {
		err = e
	} else {
		depth := depth.Int()
		ret = g.BoolOf(depth > 0)
	}
	return
}

func (op *During) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	name := lang.Underscore(op.PatternName)
	if depth, e := run.GetField(meta.PatternRunning, name); e != nil {
		err = cmdError(op, e)
	} else {
		depth := depth.Int()
		ret = g.IntOf(depth)
	}
	return
}
