package object

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RecordDot) GetReference(run rt.Runtime) (ret rt.Reference, err error) {
	if rec, e := safe.GetRecord(run, op.Value); e != nil {
		err = e
	} else if path, e := resolveDots(run, op.Dot); e != nil {
		err = e
	} else {
		at := dot.MakeReferenceValue(run, rec)
		ret, err = dot.Path(at, path)
	}
	return
}

func (op *RecordDot) GetBool(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Bool)
}

func (op *RecordDot) GetNum(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Num)
}

func (op *RecordDot) GetText(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Text)
}

func (op *RecordDot) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Record)
}

func (op *RecordDot) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.NumList)
}

func (op *RecordDot) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.TextList)
}

func (op *RecordDot) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.RecordList)
}

func (op *RecordDot) getValue(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	if at, e := op.GetReference(run); e != nil {
		err = cmd.Error(op, e)
	} else if val, e := at.GetValue(); e != nil {
		err = cmd.Error(op, e)
	} else if e := safe.Check(val, aff); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = val
	}
	return
}
