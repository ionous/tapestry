package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *Slice) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.sliceList(run, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.FloatsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.sliceList(run, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.StringsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, t, e := op.sliceList(run, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.RecordsOf(t, nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) sliceList(run rt.Runtime, aff affine.Affinity) (retVal g.Value, retType string, err error) {
	if els, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = e
	} else if e := safe.Check(els, aff); e != nil {
		err = e
	} else if i, j, e := op.getIndices(run, els.Len()); e != nil {
		err = e
	} else {
		if i >= 0 && j >= i {
			retVal, err = els.Slice(i, j)
		}
		if err == nil {
			retType = els.Type()
		}
	}
	return
}

// reti is < 0 to indicate an empty list
func (op *Slice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := safe.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if j, e := safe.GetOptionalNumber(run, op.End, 0); e != nil {
		err = e
	} else {
		reti = clipStart(i.Int(), cnt)
		retj = clipEnd(j.Int(), cnt)
	}
	return
}
