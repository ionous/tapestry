package assign

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

type Assignment interface {
	GetAssignedValue(run rt.Runtime) (ret g.Value, err error)
}

func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetBool(run, op.Value)
}
func (op *FromNumber) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetNumber(run, op.Value)
}
func (op *FromText) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetText(run, op.Value)
}
func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetRecord(run, op.Value)
}
func (op *FromNumList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetNumList(run, op.Value)
}
func (op *FromTextList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetTextList(run, op.Value)
}
func (op *FromRecordList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return safe.GetRecordList(run, op.Value)
}

// handles null assignments by returning "MissingEval" error
// ( cant live in package safe because package assign uses package safe )
func GetSafeAssignment(run rt.Runtime, a Assignment) (ret g.Value, err error) {
	if a == nil {
		err = safe.MissingEval("assigned value")
	} else {
		ret, err = a.GetAssignedValue(run)
	}
	return
}

func GetAffinity(a Assignment) (ret affine.Affinity) {
	if a != nil {
		switch a.(type) {
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
			panic(errutil.Fmt("unknown Assignment %T", a))
		}
	}
	return
}
