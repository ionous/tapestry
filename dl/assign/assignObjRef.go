package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ObjectRef) GetObjectName(run rt.Runtime) (ret string, err error) {
	// note: ObjectText can return a valid empty string; and here i think we want to error
	// so doing this manually.
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = CmdErrorCtx(op, "get object text", e)
	} else if id, e := run.GetField(meta.ObjectId, name.String()); e != nil {
		err = CmdErrorCtx(op, "get object name", e)
	} else {
		ret = id.String()
	}
	return
}

func (op *ObjectRef) GetFieldName(run rt.Runtime) (ret string, err error) {
	if name, e := safe.GetText(run, op.Field); e != nil {
		err = CmdErrorCtx(op, "get field text", e)
	} else {
		ret = name.String()
	}
	return
}

func (op *ObjectRef) GetPath() []Dot {
	return op.Dot
}

func (op *ObjectRef) GetBool(run rt.Runtime) (ret g.Value, err error) {
	var u g.Unknown
	if v, e := op.getValue(run, affine.Bool); e == nil {
		ret = v
	} else if errors.As(e, &u) && u.IsUnknownField() {
		// asking for a boolean field that doesn't exist?
		// we allow this so that any object can support trait requests
		ret = g.False
	} else {
		err = cmdError(op, e)
	}
	return
}

func (op *ObjectRef) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.Number)
}

func (op *ObjectRef) GetText(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.Text)
}

func (op *ObjectRef) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.Record)
}

func (op *ObjectRef) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.NumList)
}

func (op *ObjectRef) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.TextList)
}

func (op *ObjectRef) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.getValue(run, affine.RecordList)
}

func (op *ObjectRef) getValue(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if src, e := GetRootValue(run, op); e != nil {
		err = e
	} else {
		ret, err = src.GetCheckedValue(run, aff)
	}
	if err != nil {
		err = CmdErrorCtx(op, "get value of "+aff.String(), err)
	}
	return
}
