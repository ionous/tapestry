package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *Assignment) IsValid() bool {
	return len(op.Choice) > 0 && op.Value != nil
}

// tdb? does Assignment being a "swap" really buy anything?
// if anything, arent swaps double specifications?
func (op *Assignment) GetValue(run rt.Runtime) (ret g.Value, err error) {
	switch src := op.Value.(type) {
	case *FromBool:
		ret, err = safe.GetBool(run, src.Val)
	case *FromNumber:
		ret, err = safe.GetNumber(run, src.Val)
	case *FromText:
		ret, err = safe.GetText(run, src.Val)
	case *FromRecord:
		ret, err = safe.GetRecord(run, src.Val)
	case *FromNumList:
		ret, err = safe.GetNumList(run, src.Val)
	case *FromTextList:
		ret, err = safe.GetTextList(run, src.Val)
	case *FromRecordList:
		ret, err = safe.GetRecordList(run, src.Val)
	default:
		err = safe.MissingEval("source value")
	}
	return
}

func (op *Assignment) GetList(run rt.Runtime) (ret g.Value, err error) {
	if els, e := op.GetValue(run); e != nil {
		err = e
	} else if listAff := els.Affinity(); !affine.IsList(listAff) {
		err = errutil.New(listAff, "isn't a list")
	} else {
		ret = els
	}
	return
}

func (op *Assignment) Affinity() (ret affine.Affinity) {
	switch src := op.Choice; src {
	case Assignment_Bool_Opt:
		ret = affine.Bool
	case Assignment_Number_Opt:
		ret = affine.Number
	case Assignment_Text_Opt:
		ret = affine.Text
	case Assignment_Record_Opt:
		ret = affine.Record
	case Assignment_NumList_Opt:
		ret = affine.NumList
	case Assignment_TextList_Opt:
		ret = affine.TextList
	case Assignment_RecordList_Opt:
		ret = affine.RecordList
	}
	return
}
