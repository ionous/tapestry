package list

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func (op *ListLen) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if els, e := op.List.GetList(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = g.IntOf(els.Len())
	}
	return
}
