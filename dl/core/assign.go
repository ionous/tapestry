package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Assign) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = cmdError(op, e)
	} else if e := run.SetField(meta.Variables, op.Var.String(), v); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *FromBool) Affinity() affine.Affinity {
	return affine.Bool
}
func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (op *FromNum) Affinity() affine.Affinity {
	return affine.Number
}
func (op *FromNum) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetNumber(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (op *FromText) Affinity() affine.Affinity {
	return affine.Text
}
func (op *FromText) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetText(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (op *FromRecord) Affinity() affine.Affinity {
	return affine.Record
}
func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.GetRecord(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (op *FromNumbers) Affinity() affine.Affinity {
	return affine.NumList
}
func (op *FromNumbers) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetNumList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (op *FromTexts) Affinity() affine.Affinity {
	return affine.TextList
}
func (op *FromTexts) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetTextList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (op *FromRecords) Affinity() affine.Affinity {
	return affine.RecordList
}
func (op *FromRecords) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if objs, e := safe.GetRecordList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = objs
	}
	return
}
