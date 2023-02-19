package core

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func (op *CallPattern) Execute(run rt.Runtime) error {
	_, err := op.determine(run, affine.None)
	return err
}

func (op *CallPattern) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallPattern) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

func (op *CallPattern) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
}

func (op *CallPattern) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *CallPattern) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *CallPattern) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *CallPattern) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.RecordList)
}

func (op *CallPattern) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	pat := op.Pattern.String()
	if rec, e := MakeRecord(run, pat, op.Arguments...); e != nil {
		err = cmdError(op, e)
	} else if v, e := run.Call(rec, aff); e != nil && !errors.Is(e, rt.NoResult{}) {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
