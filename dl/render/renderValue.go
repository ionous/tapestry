package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RenderValue) RenderEval(run rt.Runtime, hint affine.Affinity) (ret g.Value, err error) {
	if v, e := assign.GetSafeAssignment(run, op.Value); e != nil {
		err = CmdError(op, e)
	} else if safe.Check(v, hint); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}
