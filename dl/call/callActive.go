package call

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *ActiveScene) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	return run.GetField(meta.Domain, op.Name)
}

// GetBool returns the first matching bool evaluation.
func (op *ActivePattern) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if depth, e := op.GetNum(run); e != nil {
		err = e
	} else {
		depth := depth.Int()
		ret = rt.BoolOf(depth > 0)
	}
	return
}

func (op *ActivePattern) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	name := inflect.Normalize(op.PatternName)
	if depth, e := run.GetField(meta.PatternRunning, name); e != nil {
		err = cmd.Error(op, e)
	} else {
		depth := depth.Int()
		ret = rt.IntOf(depth)
	}
	return
}
