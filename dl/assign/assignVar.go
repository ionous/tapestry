package assign

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *VariableDot) GetReference(run rt.Runtime) (ret dot.Endpoint, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = e
	} else if path, e := ResolvePath(run, op.Dot); e != nil {
		err = e
	} else {
		ret, err = resolveVariable(run, name.String(), path)
	}
	return
}

func resolveVariable(run rt.Runtime, name string, path dot.Path) (ret dot.Endpoint, err error) {
	field := dot.Field(name)
	fullPath := append(dot.Path{field}, path...)
	return dot.FindEndpoint(run, meta.Variables, fullPath)
}

func (op *VariableDot) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.Bool); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetNumber(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.Number); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.Text); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetRecord(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.Record); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetNumList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.NumList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.TextList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) GetRecordList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getValue(run, affine.RecordList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *VariableDot) getValue(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	if src, e := op.GetReference(run); e != nil {
		err = e
	} else if val, e := src.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	return
}
