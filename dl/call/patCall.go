package call

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *CallPattern) Execute(run rt.Runtime) error {
	_, err := op.determine(run, affine.None)
	return err
}

func (op *CallPattern) GetBool(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallPattern) GetNum(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Num)
}

func (op *CallPattern) GetText(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Text)
}

func (op *CallPattern) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *CallPattern) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *CallPattern) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *CallPattern) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.RecordList)
}

// note: at one point this would unwrap errors so that callers couldn't see them
// i no longer am sure why. doing stops game.Signals(s) ( ex SignalQuit ) from reaching the parser.
func (op *CallPattern) determine(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	name := inflect.Normalize(op.PatternName)
	if k, v, e := assign.ExpandArgs(run, op.Arguments); e != nil {
		err = cmd.ErrorCtx(op, name, e)
	} else if v, e := run.Call(name, aff, k, v); e != nil {
		err = cmd.ErrorCtx(op, name, e)
	} else {
		ret = v
	}
	return
}
