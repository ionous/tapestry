package literal

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// LiteralValue marks script constants.
type LiteralValue interface {
	GetLiteralValue(g.Kinds) (g.Value, error)
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Value)
}

func (op *BoolValue) GetLiteralValue(g.Kinds) (ret g.Value, _ error) {
	ret = g.BoolFrom(op.Value, op.Kind)
	return
}

// GetBool implements rt.BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(run rt.Runtime) (g.Value, error) {
	return op.GetLiteralValue(run)
}

// Int converts to native int.
func (op *NumValue) Int() int {
	return int(op.Value)
}

// Float converts to native float.
func (op *NumValue) Float() float64 {
	return op.Value
}

// String returns a nicely formatted float, with no decimal point when possible.
func (op *NumValue) String() string {
	return strconv.FormatFloat(op.Value, 'g', -1, 64)
}

func (op *NumValue) GetLiteralValue(g.Kinds) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Value)
	return
}

// GetNumber implements rt.NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.GetLiteralValue(run)
}

// String returns the text.
func (op *TextValue) String() string {
	return op.Value
}

func (op *TextValue) GetLiteralValue(g.Kinds) (ret g.Value, _ error) {
	ret = g.StringFrom(op.Value, op.Kind)
	return
}

// GetText implements interface rt.TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (g.Value, error) {
	return op.GetLiteralValue(run)
}

func (op *NumValues) GetLiteralValue(g.Kinds) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.FloatsOf(op.Values)
	return
}

// GetNumList implements rt.NumListEval providing the dl with a literal list of numbers.
func (op *NumValues) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.GetLiteralValue(run)
}

func (op *TextValues) GetLiteralValue(g.Kinds) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.StringsOf(op.Values)
	return
}

// GetTextList implements rt.TextListEval providing the dl with a literal list of text.
func (op *TextValues) GetTextList(run rt.Runtime) (ret g.Value, _ error) {
	return op.GetLiteralValue(run)
}

func (op *RecordValue) GetLiteralValue(kinds g.Kinds) (ret g.Value, err error) {
	return op.Cache.GetRecord(kinds, op.Kind, op.Fields)
}

// GetRecord implements interface rt.RecordEval providing the dl with a structured literal.
func (op *RecordValue) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.GetLiteralValue(run)
}

func (op *RecordList) GetLiteralValue(kinds g.Kinds) (g.Value, error) {
	return op.Cache.GetRecords(kinds, op.Kind, op.Records)
}

// GetNumList implements rt.RecordListEval providing the dl with a literal list of records.
func (op *RecordList) GetRecordList(run rt.Runtime) (ret g.Value, _ error) {
	return op.GetLiteralValue(run)
}

// unimplemented: panics.
func (op *FieldList) String() (ret string) {
	panic("field values are not intended to be comparable")
}

// unimplemented: panics.
func (op *FieldList) GetLiteralValue(g.Kinds) (g.Value, error) {
	panic("field values should only be used in record literals")
}

func GetAffinity(a LiteralValue) (ret affine.Affinity) {
	if a != nil {
		switch a.(type) {
		case *BoolValue:
			ret = affine.Bool
		case *NumValue:
			ret = affine.Number
		case *TextValue:
			ret = affine.Text
		case *RecordValue:
			ret = affine.Record
		case *NumValues:
			ret = affine.NumList
		case *TextValues:
			ret = affine.TextList
		case *RecordList:
			ret = affine.RecordList
		default:
			panic(errutil.Fmt("unknown Assignment %T", a))
		}
	}
	return
}
