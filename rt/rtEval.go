package rt

import (
	"git.sr.ht/~ionous/iffy/affine"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// Execute runs a bit of code that has no return value.
type Execute interface {
	Execute(Runtime) error
}

// BoolEval represents the result of some true or false expression.
type BoolEval interface {
	GetBool(Runtime) (g.Value, error)
}

// NumberEval represents the result of some numeric expression.
type NumberEval interface {
	GetNumber(Runtime) (g.Value, error)
}

// TextEval represents the result of some expression which creates a string.
type TextEval interface {
	GetText(Runtime) (g.Value, error)
}

// RecordEval yields access to a set of fields and their values.
type RecordEval interface {
	GetRecord(Runtime) (g.Value, error)
}

// TextListEval represents the computation of a series of numeric values.
type NumListEval interface {
	GetNumList(Runtime) (g.Value, error)
}

// TextListEval represents the computation of a series of strings.
type TextListEval interface {
	GetTextList(Runtime) (g.Value, error)
}

// RecordListEval represents the computation of a series of a set of fields.
type RecordListEval interface {
	GetRecordList(Runtime) (g.Value, error)
}

func AffineOfEval(eval interface{}) (ret affine.Affinity) {
	switch eval.(type) {
	case BoolEval:
		ret = affine.Bool
	case NumberEval:
		ret = affine.Number
	case TextEval:
		ret = affine.Text
	case NumListEval:
		ret = affine.NumList
	case TextListEval:
		ret = affine.TextList
	case RecordEval:
		ret = affine.Record
	case RecordListEval:
		ret = affine.RecordList
	}
	return
}
