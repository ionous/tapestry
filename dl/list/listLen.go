package list

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListLen) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetList(run, op.List); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(v.Len())
	}
	return
}
