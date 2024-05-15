package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ObjectDot) GetReference(run rt.Runtime) (ret dot.Endpoint, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = e
	} else if id, e := run.GetField(meta.ObjectId, name.String()); e != nil {
		err = e
	} else if path, e := ResolvePath(run, op.Dot); e != nil {
		err = e
	} else {
		ret, err = dot.FindEndpoint(run, id.String(), path)
	}
	return
}

func (op *ObjectDot) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	var u rt.Unknown
	if v, e := op.getValue(run, affine.Bool); e == nil {
		ret = v
	} else if errors.As(e, &u) && u.IsUnknownField() {
		// asking for a boolean field that doesn't exist?
		// we allow this so that any object can support trait requests
		// fix: this should somehow validate that there is such a trait however
		// [ ex. return "inapplicable trait" instead of "unknown field" ]
		// bonus points for determining this during weave when using literals
		ret = rt.False
	} else {
		err = CmdError(op, e)
	}
	return
}

func (op *ObjectDot) GetNumber(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Number)
}

func (op *ObjectDot) GetText(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Text)
}

func (op *ObjectDot) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Record)
}

func (op *ObjectDot) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.NumList)
}

func (op *ObjectDot) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.TextList)
}

func (op *ObjectDot) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.RecordList)
}

func (op *ObjectDot) getValue(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	if at, e := GetReference(run, op); e != nil {
		err = e
	} else if val, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	if err != nil {
		err = CmdErrorCtx(op, "get value of "+aff.String(), err)
	}
	return
}
