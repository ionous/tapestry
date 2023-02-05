package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *GetFromVar) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := GetFromVariable(run, name.String(), affine.Bool, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := GetFromVariable(run, name.String(), affine.Number, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetText(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := GetFromVariable(run, name.String(), affine.Text, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetList(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := GetFromVariable(run, name.String(), affine.List, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := GetFromVariable(run, name.String(), affine.Record, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// FIX: instead of check, convert
func GetFromVariable(run rt.Runtime, name string, aff affine.Affinity, path []Dot) (ret g.Value, err error) {
	if val, e := run.GetField(meta.Variables, name); e != nil {
		err = e
	} else if val, e := Peek(run, val, path); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	return
}
