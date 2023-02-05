package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *GetFromVar) GetBool(run rt.Runtime) (g.Value, error) {
	return GetFromVariable(run, op.Name, affine.Bool, op.Path)
}

func (op *GetFromVar) GetNumber(run rt.Runtime) (g.Value, error) {
	return GetFromVariable(run, op.Name, affine.Number, op.Path)
}

func (op *GetFromVar) GetText(run rt.Runtime) (g.Value, error) {
	return GetFromVariable(run, op.Name, affine.Text, op.Path)
}

func (op *GetFromVar) GetList(run rt.Runtime) (g.Value, error) {
	return GetFromVariable(run, op.Name, affine.List, op.Path)
}

func (op *GetFromVar) GetRecord(run rt.Runtime) (g.Value, error) {
	return GetFromVariable(run, op.Name, affine.Record, op.Path)
}

// FIX: instead of check, convert
func GetFromVariable(run rt.Runtime, name rt.TextEval, aff affine.Affinity, path []PathEval) (ret g.Value, err error) {
	if name, e := safe.GetText(run, name); e != nil {
		err = e
	} else if val, e := run.GetField(meta.Variables, name.String()); e != nil {
		err = e
	} else if val, e := PickValue(run, val, path); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	return
}
