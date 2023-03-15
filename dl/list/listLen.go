package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *ListLen) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if els, e := safe.GetAssignment(run, op.List); e != nil {
		err = CmdError(op, e)
	} else if !affine.IsList(els.Affinity()) {
		err = CmdError(op, errutil.New("not a list"))
	} else {
		ret = g.IntOf(els.Len())
	}
	return
}
