package logic

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

func (op *ChooseNum) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	var pop bool
	if ks, vs, e := call.ExpandArgs(run, op.Args); e != nil {
		err = cmd.Error(op, e)
	} else {
		if pop = len(vs) > 0; pop {
			run.PushScope(scope.NewPairs(ks, vs))
		}
		if b, e := safe.GetBool(run, op.If); e != nil {
			err = cmd.Error(op, e)
		} else {
			var next rt.NumEval
			if b.Bool() {
				next = op.True
			} else {
				next = op.False
			}
			if v, e := safe.GetOptionalNumber(run, next, 0); e != nil {
				err = cmd.Error(op, e)
			} else {
				ret = v
			}
		}
		// on success or error
		if pop {
			run.PopScope()
		}
	}
	return
}

func (op *ChooseText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var pop bool
	if ks, vs, e := call.ExpandArgs(run, op.Args); e != nil {
		err = cmd.Error(op, e)
	} else {
		if pop = len(vs) > 0; pop {
			run.PushScope(scope.NewPairs(ks, vs))
		}
		if b, e := safe.GetBool(run, op.If); e != nil {
			err = cmd.Error(op, e)
		} else {
			var next rt.TextEval
			if b.Bool() {
				next = op.True
			} else {
				next = op.False
			}
			if v, e := safe.GetOptionalText(run, next, ""); e != nil {
				err = cmd.Error(op, e)
			} else {
				ret = v
			}
		}
		// on success or error
		if pop {
			run.PopScope()
		}
	}
	return
}
