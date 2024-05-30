package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListLen) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = CmdError(op, e)
	} else if !affine.IsList(els.Affinity()) {
		err = CmdError(op, errutil.New("not a list"))
	} else {
		ret = rt.IntOf(els.Len())
	}
	return
}
