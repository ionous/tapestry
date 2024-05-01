package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ObjectRef) GetReference(run rt.Runtime) (ret dot.Endpoint, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = e
	} else if id, e := run.GetField(meta.ObjectId, name.String()); e != nil {
		err = e
	} else if fieldName, e := safe.GetText(run, op.Field); e != nil {
		err = e
	} else if path, e := ResolvePath(run, op.Dot); e != nil {
		err = e
	} else {
		field := dot.Field(fieldName.String())
		fullPath := append(dot.Path{field}, path...)
		ret, err = dot.FindEndpoint(run, id.String(), fullPath)
	}
	return
}

func (op *ObjectRef) GetBool(run rt.Runtime) (ret g.Value, err error) {
	var u g.Unknown
	if v, e := op.getValue(run, affine.Bool); e == nil {
		ret = v
	} else if errors.As(e, &u) && u.IsUnknownField() {
		// asking for a boolean field that doesn't exist?
		// we allow this so that any object can support trait requests
		// fix: this should somehow validate that there is such a trait however
		// [ ex. return "inapplicable trait" instead of "unknown field" ]
		// bonus points for determining this during weave when using literals
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
