package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

func (op *RenderValue) RenderEval(run rt.Runtime, hint affine.Affinity) (ret g.Value, err error) {
	if aff := op.Value.Affinity(); hint != aff {
		err = CmdError(op, errutil.New("mismatched affinity", aff, hint))
	} else if v, e := op.Value.GetValue(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}
