package list

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ListTextAt) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if els, e := safe.GetTextList(run, op.List); e != nil {
		err = cmd.Error(op, e)
	} else if v, e := getValueAt(run, els, op.Index); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *ListNumAt) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if els, e := safe.GetNumList(run, op.List); e != nil {
		err = cmd.Error(op, e)
	} else if v, e := getValueAt(run, els, op.Index); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func getValueAt(run rt.Runtime, els rt.Value, num rt.NumEval) (ret rt.Value, err error) {
	if idx, e := safe.GetNum(run, num); e != nil {
		err = e
	} else if idx := idx.Int(); idx <= 0 {
		err = fmt.Errorf("underflow %d", idx)
	} else if cnt := els.Len(); idx > cnt {
		err = fmt.Errorf("overflow %d of %d elements", idx, cnt)
	} else {
		ret = els.Index(idx - 1)
	}
	return
}
