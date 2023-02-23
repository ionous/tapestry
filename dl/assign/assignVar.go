package assign

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// always returns meta.Variables as its object
func (op *VariableRef) GetObjectName(run rt.Runtime) (string, error) {
	return meta.Variables, nil
}

func (op *VariableRef) GetFieldName(run rt.Runtime) (ret string, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else {
		ret = name.String()
	}
	return
}

func (op *VariableRef) GetPath() []Dot {
	return op.Dot
}

func (op *VariableRef) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Bool); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Number); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.Record); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getValue(run, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableRef) getValue(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if src, e := GetRootValue(run, op); e != nil {
		err = e
	} else {
		ret, err = src.GetCheckedValue(run, aff)
	}
	return
}
