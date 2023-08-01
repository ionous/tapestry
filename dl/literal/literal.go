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
	GetLiteralValue(rt.Runtime) (g.Value, error)
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Value)
}

func (op *BoolValue) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetBool(run)
}

// GetBool implements rt.BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Value)
	return
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

func (op *NumValue) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetNumber(run)
}

// GetNumber implements rt.NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Value)
	return
}

// String returns the text.
func (op *TextValue) String() string {
	return op.Value
}

func (op *TextValue) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetText(run)
}

// GetText implements interface rt.TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringFrom(op.Value, op.Kind)
	return
}

func (op *NumValues) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetNumList(run)
}

// GetNumList implements rt.NumListEval providing the dl with a literal list of numbers.
func (op *NumValues) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.FloatsOf(op.Values)
	return
}

func (op *TextValues) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetTextList(run)
}

// GetTextList implements rt.TextListEval providing the dl with a literal list of text.
func (op *TextValues) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.StringsOf(op.Values)
	return
}

func (op *RecordValue) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetRecord(run)
}

// GetRecord implements interface rt.RecordEval providing the dl with a structured literal.
func (op *RecordValue) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.Cache.GetRecord(run, op.Kind, op.Fields)
}

func (op *RecordList) GetLiteralValue(run rt.Runtime) (g.Value, error) {
	return op.GetRecordList(run)
}

// GetNumList implements rt.RecordListEval providing the dl with a literal list of records.
func (op *RecordList) GetRecordList(run rt.Runtime) (ret g.Value, _ error) {
	return op.Cache.GetRecords(run, op.Kind, op.Records)
}

// unimplemented: panics.
func (op *FieldList) String() (ret string) {
	panic("field values are not intended to be comparable")
}

// unimplemented: panics.
func (op *FieldList) GetLiteralValue(run rt.Runtime) (g.Value, error) {
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
