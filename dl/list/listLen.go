package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListLen) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if !affine.IsList(v.Affinity()) {
		err = cmdError(op, errutil.New("not a list"))
	} else {
		ret = g.IntOf(v.Len())
	}
	return
}
