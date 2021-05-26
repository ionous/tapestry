package core

import (
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"

	g "git.sr.ht/~ionous/iffy/rt/generic"
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
	name := lang.Underscore(op.Pattern.Value())
	if depth, e := run.GetField(object.Running, name); e != nil {
		err = cmdError(op, e)
	} else {
		depth := depth.Int()
		ret = g.IntOf(depth)
	}
	return
}
