package object

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *VariableDot) GetReference(run rt.Runtime) (ret rt.Reference, err error) {
	if varName, e := safe.GetText(run, op.VariableName); e != nil {
		err = e
	} else if path, e := resolveDots(run, op.Dot); e != nil {
		err = e
	} else {
		// we create a reference containing the variable name
		at := dot.MakeReference(run, meta.Variables)
		if at, e := at.Dot(dot.Field(varName.String())); e != nil {
			err = e
		} else {
			// walk the full path ( if any ) to expand it.
			ret, err = dot.Path(at, path)
		}
	}
	return
}

func (op *VariableDot) GetBool(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Bool)
}

func (op *VariableDot) GetNum(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Num)
}

func (op *VariableDot) GetText(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Text)
}

func (op *VariableDot) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Record)
}

func (op *VariableDot) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.NumList)
}

func (op *VariableDot) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.TextList)
}

func (op *VariableDot) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.RecordList)
}

// uses GetReference, which expands the full path to get the targeted value.
func (op *VariableDot) getValue(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
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
