package logic

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Not) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.BoolOf(!val.Bool())
	}
	return
}
