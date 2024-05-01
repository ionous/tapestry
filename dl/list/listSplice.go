package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListSplice) Execute(run rt.Runtime) (err error) {
	if _, _, e := op.spliceList(run, ""); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListSplice) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.spliceList(run, affine.NumList); e != nil {
		err = CmdError(op, e)
	} else if v == nil {
		ret = g.FloatsOf(nil)
	} else {
		ret = v
	}
	return
}
func (op *ListSplice) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.spliceList(run, affine.TextList); e != nil {
		err = CmdError(op, e)
	} else if v == nil {
		ret = g.StringsOf(nil)
	} else {
		ret = v
	}
	return
}
func (op *ListSplice) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, t, e := op.spliceList(run, affine.RecordList); e != nil {
		err = CmdError(op, e)
	} else if v == nil {
		ret = g.RecordsFrom(nil, t)
	} else {
		ret = v
	}
	return
}

// modify a list by adding and removing elements.
// fix: aff?
func (op *ListSplice) spliceList(run rt.Runtime, _ affine.Affinity) (retVal g.Value, retType string, err error) {
	if at, e := assign.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if ins, e := safe.GetAssignment(run, op.Insert); e != nil {
		err = e
	} else if !IsAppendable(ins, vs) {
		err = insertError{ins, vs}
	} else if i, j, e := op.getIndices(run, vs.Len()); e != nil {
		err = e
	} else {
		if i >= 0 && j >= i {
			retVal, err = vs.Splice(i, j, ins)
		}
		if err == nil {
			retType = vs.Type()
		}
	}
	return
}

func (op *ListSplice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := safe.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if rng, e := safe.GetOptionalNumber(run, op.Remove, float64(cnt)); e != nil {
		err = e
	} else {
		reti = clipStart(i.Int(), cnt)
		retj = clipRange(reti, rng.Int(), cnt)
	}
	return
}
