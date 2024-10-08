package object

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ObjectDot) GetReference(run rt.Runtime) (ret rt.Reference, err error) {
	if name, e := safe.ObjectText(run, op.NounName); e != nil {
		err = e
	} else if path, e := resolveDots(run, op.Dot); e != nil {
		err = e
	} else if str := name.String(); len(str) == 0 {
		err = fmt.Errorf("object is nothing")
	} else {
		pos := dot.MakeReference(run, str)
		ret, err = dot.Path(pos, path)
	}
	return
}

func (op *ObjectDot) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if len(op.Dot) > 0 {
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
			err = cmd.Error(op, e)
		}
	} else {
		switch obj, e := safe.ObjectText(run, op.NounName); e.(type) {
		case rt.Unknown:
			ret = rt.False // no such object
		case nil:
			if len(obj.String()) == 0 {
				ret = rt.False // the eval returned the empty string
			} else {
				ret = rt.True
			}
		default:
			err = cmd.Error(op, e)
		}
	}
	return
}

func (op *ObjectDot) GetNum(run rt.Runtime) (rt.Value, error) {
	return op.getValue(run, affine.Num)
}

// as a special case, if there are no dot parts, return the id of the object
func (op *ObjectDot) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if len(op.Dot) > 0 {
		ret, err = op.getValue(run, affine.Text)
	} else {
		ret, err = safe.ObjectText(run, op.NounName)
	}
	return
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
	if pos, e := safe.GetReference(run, op); e != nil {
		err = e
	} else if val, e := pos.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	if err != nil {
		err = cmd.ErrorCtx(op, "get value of "+aff.String(), err)
	}
	return
}
