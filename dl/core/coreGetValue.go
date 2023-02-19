package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func (op *GetValue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Bool); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Number); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Record); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetValue) getValue(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if src, e := op.Source.GetRootValue(run); e != nil {
		err = e
	} else {
		ret, err = src.GetCheckedValue(run, aff)
	}
	return
}
