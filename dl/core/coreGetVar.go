package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *GetFromVar) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromVar(run, op.Name, affine.Bool, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromVar(run, op.Name, affine.Number, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromVar(run, op.Name, affine.Text, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromVar(run, op.Name, affine.List, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromVar) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromVar(run, op.Name, affine.Record, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// FIX: convert and warn instead of error on field affinity checks
func getFromVar(run rt.Runtime, name rt.TextEval, aff affine.Affinity, path []Dot) (ret g.Value, err error) {
	if name, e := safe.GetText(run, name); e != nil {
		err = e
	} else if val, e := run.GetField(meta.Variables, name.String()); e != nil {
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
