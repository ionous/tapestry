package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Not) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(!val.Bool())
	}
	return
}
