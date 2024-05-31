package render

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// RenderRef reads a value using a name which might refer to a variable or an object.
// If its an object, the dot will reference some particular field, otherwise turns the object into an id.
// Intended for internal use.

func (op *RenderRef) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.Bool); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.Num); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

// GetText handles unpacking a text variable,
// turning an object variable into an id, or
// looking for an object of the passed name ( if no variable of the name exists. )
func (op *RenderRef) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.Text); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) GetRecord(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.Record); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) GetNumList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.NumList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.TextList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) GetRecordList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, affine.RecordList); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) RenderEval(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if v, e := op.renderRef(run, hint); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) renderRef(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = fmt.Errorf("%w getting text", e)
	} else if val, e := object.GetNamedValue(run, name.String(), op.Dot); e != nil {
		err = e
	} else {
		ret, err = safe.ConvertValue(run, val, hint)
	}
	return
}
