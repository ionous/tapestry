package call

import (
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *FromExe) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	err = safe.RunAll(run, op.Exe)
	return // what should we return...?
}

func (op *FromAddress) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	if pos, e := safe.GetReference(run, op.Value); e != nil {
		err = e
	} else {
		ret, err = pos.GetValue()
	}
	return
}

func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetBool(run, op.Value)
}
func (op *FromNum) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetNum(run, op.Value)
}
func (op *FromText) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetText(run, op.Value)
}
func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetRecord(run, op.Value)
}
func (op *FromNumList) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetNumList(run, op.Value)
}
func (op *FromTextList) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetTextList(run, op.Value)
}
func (op *FromRecordList) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetRecordList(run, op.Value)
}

func GetAffinity(a rt.Assignment) (ret affine.Affinity) {
	if a != nil {
		switch a.(type) {
		case *FromExe:
			ret = affine.None
		case *FromBool:
			ret = affine.Bool
		case *FromNum:
			ret = affine.Num
		case *FromText:
			ret = affine.Text
		case *FromRecord:
			ret = affine.Record
		case *FromNumList:
			ret = affine.NumList
		case *FromTextList:
			ret = affine.TextList
		case *FromRecordList:
			ret = affine.RecordList
		default:
			log.Panicf("unknown Assignment %T", a)
		}
	}
	return
}

// turn a literal into an assignment
// ( whats's the right package for this function? )
func Literal(v literal.LiteralValue) (ret rt.Assignment) {
	switch v := v.(type) {
	case *literal.BoolValue:
		ret = &FromBool{Value: v}
	case *literal.NumValue:
		ret = &FromNum{Value: v}
	case *literal.TextValue:
		ret = &FromText{Value: v}
	case *literal.RecordValue:
		ret = &FromRecord{Value: v}
	case *literal.NumList:
		ret = &FromNumList{Value: v}
	case *literal.TextList:
		ret = &FromTextList{Value: v}
	case *literal.RecordList:
		ret = &FromRecordList{Value: v}
	default:
		log.Panicf("unknown literal %T", v)
	}
	return
}
