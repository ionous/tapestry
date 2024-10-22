package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RenderValue) RenderEval(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if v, e := safe.GetAssignment(run, op.Value); e != nil {
		err = cmd.Error(op, e)
	} else if e := safe.Check(v, hint); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}
