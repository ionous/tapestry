package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RenderValue) RenderEval(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if v, e := safe.GetAssignment(run, op.Value); e != nil {
		err = CmdError(op, e)
	} else if safe.Check(v, hint); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}
