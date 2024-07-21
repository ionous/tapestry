package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListPopNum) Execute(run rt.Runtime) error {
	return popEdge(run, op, affine.NumList, op.Target, op.Edge, nil)
}
func (op *ListPopNum) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	err = popEdge(run, op, affine.NumList, op.Target, op.Edge, &ret)
	return
}

func (op *ListPopText) Execute(run rt.Runtime) error {
	return popEdge(run, op, affine.TextList, op.Target, op.Edge, nil)
}

func (op *ListPopText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	err = popEdge(run, op, affine.TextList, op.Target, op.Edge, &ret)
	return
}

func (op *ListPopRecord) Execute(run rt.Runtime) (err error) {
	return popEdge(run, op, affine.RecordList, op.Target, op.Edge, nil)
}

func (op *ListPopRecord) GetRecord(run rt.Runtime) (ret rt.Value, err error) {
	err = popEdge(run, op, affine.Record, op.Target, op.Edge, &ret)
	return
}

func popEdge(
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
			if e := vs.Splice(idx, idx+1, nil, cutList); e != nil {
				err = e
			} else if cutList != nil {
				*cutList = (*cutList).Index(0)
			}
		}
	}
	if err != nil {
		err = cmd.Error(op, err)
	}
	return
}
