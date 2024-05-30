package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// GetBool returns the first matching bool evaluation.
func (op *During) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if depth, e := op.GetNum(run); e != nil {
		err = e
	} else {
		depth := depth.Int()
		ret = rt.BoolOf(depth > 0)
	}
	return
}

func (op *During) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	name := inflect.Normalize(op.PatternName)
	if depth, e := run.GetField(meta.PatternRunning, name); e != nil {
		err = cmdError(op, e)
	} else {
		depth := depth.Int()
		ret = rt.IntOf(depth)
	}
	return
}
