package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ChooseNum) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = cmdError(op, e)
	} else {
		var next rt.NumberEval
		if b.Bool() {
			next = op.True
		} else {
			next = op.False
		}
		if v, e := safe.GetOptionalNumber(run, next, 0); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	}
	return
}

func (op *ChooseText) GetText(run rt.Runtime) (ret g.Value, err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = cmdError(op, e)
	} else {
		var next rt.TextEval
		if b.Bool() {
			next = op.True
		} else {
			next = op.False
		}
		if v, e := safe.GetOptionalText(run, next, ""); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	}
	return
}
