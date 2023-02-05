package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *GetFromObj) GetBool(run rt.Runtime) (g.Value, error) {
	return GetFromObject(run, op.Name, op.Field, affine.Bool, op.Path)
}

func (op *GetFromObj) GetNumber(run rt.Runtime) (g.Value, error) {
	return GetFromObject(run, op.Name, op.Field, affine.Number, op.Path)
}

func (op *GetFromObj) GetText(run rt.Runtime) (g.Value, error) {
	return GetFromObject(run, op.Name, op.Field, affine.Text, op.Path)
}

func (op *GetFromObj) GetList(run rt.Runtime) (g.Value, error) {
	return GetFromObject(run, op.Name, op.Field, affine.List, op.Path)
}

func (op *GetFromObj) GetRecord(run rt.Runtime) (g.Value, error) {
	return GetFromObject(run, op.Name, op.Field, affine.Record, op.Path)
}

// FIX: instead of check, convert
func GetFromObject(run rt.Runtime, name rt.TextEval, field rt.TextEval, aff affine.Affinity, path []PathEval) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, name); e != nil {
		err = e
	} else if obj == nil {
		err = g.NothingObject
	} else if fieldName, e := safe.GetText(run, field); e != nil {
		err = e
	} else if memberValue, e := obj.FieldByName(fieldName.String()); e != nil {
		err = e
	} else if val, e := PickValue(run, memberValue, path); e != nil {
		err = e
	} else {
		ret = val
	}
	return
}
