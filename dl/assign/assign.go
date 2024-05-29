package assign

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
	if pos, e := GetReference(run, op.Value); e != nil {
		err = e
	} else {
		ret, err = pos.GetValue()
	}
	return
}

func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetBool(run, op.Value)
}
func (op *FromNumber) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
	return safe.GetNumber(run, op.Value)
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
		case *FromNumber:
			ret = affine.Number
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
		ret = &FromNumber{Value: v}
	case *literal.TextValue:
		ret = &FromText{Value: v}
	case *literal.RecordValue:
		ret = &FromRecord{Value: v}
	case *literal.NumValues:
		ret = &FromNumList{Value: v}
	case *literal.TextValues:
		ret = &FromTextList{Value: v}
	case *literal.RecordList:
		ret = &FromRecordList{Value: v}
	default:
		log.Panicf("unknown literal %T", v)
	}
	return
}

func Object(name string, path ...any) *ObjectDot {
	return &ObjectDot{
		Name: literal.T(name),
		Dot:  MakeDot(path...),
	}
}

// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func Variable(name string, path ...any) *VariableDot {
	return &VariableDot{
		Name: literal.T(name),
		Dot:  MakeDot(path...),
	}
}

func MakeDot(path ...any) (ret []Dot) {
	if cnt := len(path); cnt > 0 {
		out := make([]Dot, len(path))
		for i, p := range path {
			switch el := p.(type) {
			case string:
				out[i] = &AtField{Field: literal.T(el)}
			case int:
				out[i] = &AtIndex{Index: literal.I(el)}
			case Dot:
				out[i] = el
			default:
				log.Panicf("expected an int or string element; got %T", el)
			}
		}
		ret = out
	}
	return
}
