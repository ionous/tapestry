package list

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListPopNum) Execute(run rt.Runtime) error {
	return popList(run, op, affine.NumList, op.Target, op.Edge, nil)
}
func (op *ListPopNum) GetNum(run rt.Runtime) (rt.Value, error) {
	return popValue(run, op, affine.NumList, op.Target, op.Edge)
}

func (op *ListPopText) Execute(run rt.Runtime) error {
	return popList(run, op, affine.TextList, op.Target, op.Edge, nil)
}
func (op *ListPopText) GetText(run rt.Runtime) (rt.Value, error) {
	return popValue(run, op, affine.TextList, op.Target, op.Edge)
}

func (op *ListPopRecord) Execute(run rt.Runtime) error {
	return popList(run, op, affine.RecordList, op.Target, op.Edge, nil)
}
func (op *ListPopRecord) GetRecord(run rt.Runtime) (rt.Value, error) {
	return popValue(run, op, affine.Record, op.Target, op.Edge)
}

func popValue(
	run rt.Runtime,
	op typeinfo.Instance,
	aff affine.Affinity,
	tgt rt.Address,
	atFront rt.BoolEval) (ret rt.Value, err error) {
	var vs rt.Value
	if e := popList(run, op, aff, tgt, atFront, &vs); e != nil {
		err = e
	} else if vs.Len() == 0 {
		err = cmd.Error(op, errors.New("popped empty list"))
	} else {
		ret = vs.Index(0)
	}
	return
}

// cutList will contain a *list* of values
func popList(
	run rt.Runtime,
	op typeinfo.Instance,
	aff affine.Affinity,
	tgt rt.Address,
	atFront rt.BoolEval,
	cutList *rt.Value,
) (err error) {
	if at, e := safe.GetReference(run, tgt); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(vs, aff); e != nil {
		err = e
	} else if cnt := vs.Len(); cnt > 0 {
		var idx int
		if atFront, e := safe.GetOptionalBool(run, atFront, false); e != nil {
			err = e
		} else {
			if !atFront.Bool() {
				idx = cnt - 1
			}
			err = vs.Splice(idx, idx+1, nil, cutList)
		}
	}
	if err != nil {
		err = cmd.Error(op, err)
	}
	return
}
