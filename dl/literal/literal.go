package literal

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// LiteralValue marks script constants.
// ( same interface as core.Assignment currently )
type LiteralValue interface {
	Affinity() affine.Affinity // alt: use a switch on the eval type to generate affinity.
	GetAssignedValue(rt.Runtime) (g.Value, error)
}

// Affinity returns affine.Bool
func (op *BoolValue) Affinity() affine.Affinity {
	return affine.Bool
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Value)
}

func (op *BoolValue) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetBool(run)
}

// GetBool implements rt.BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Value)
	return
}

// Affinity returns affine.Number
func (op *NumValue) Affinity() affine.Affinity {
	return affine.Number
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

func (op *NumValue) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetNumber(run)
}

// GetNumber implements rt.NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Value)
	return
}

// Affinity returns affine.Text
func (op *TextValue) Affinity() affine.Affinity {
	return affine.Text
}

// String returns the text.
func (op *TextValue) String() string {
	return op.Value
}

func (op *TextValue) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetText(run)
}

// GetText implements interface rt.TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringOf(op.Value)
	return
}

// Affinity returns affine.NumList
func (op *NumValues) Affinity() affine.Affinity {
	return affine.NumList
}

func (op *NumValues) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetNumList(run)
}

// GetNumList implements rt.NumListEval providing the dl with a literal list of numbers.
func (op *NumValues) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.FloatsOf(op.Values)
	return
}

// Affinity returns affine.TextList
func (op *TextValues) Affinity() affine.Affinity {
	return affine.TextList
}

func (op *TextValues) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetTextList(run)
}

// GetTextList implements rt.TextListEval providing the dl with a literal list of text.
func (op *TextValues) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	// fix? would copying be better?
	ret = g.StringsOf(op.Values)
	return
}

// Affinity returns affine.Record
func (op *RecordValue) Affinity() affine.Affinity {
	return affine.Record
}

func (op *RecordValue) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetRecord(run)
}

// GetRecord implements interface rt.RecordEval providing the dl with a structured literal.
func (op *RecordValue) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.Cache.GetRecord(run, op.Kind, op.Fields)
}

// Affinity returns affine.RecordList
func (op *RecordList) Affinity() affine.Affinity {
	return affine.RecordList
}

func (op *RecordList) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.GetRecordList(run)
}

// GetNumList implements rt.RecordListEval providing the dl with a literal list of records.
func (op *RecordList) GetRecordList(run rt.Runtime) (ret g.Value, _ error) {
	return op.Cache.GetRecords(run, op.Kind, op.Records)
}

// unimplemented: returns empty string.
func (op *FieldList) Affinity() affine.Affinity {
	return ""
}

// unimplemented: panics.
func (op *FieldList) String() (ret string) {
	panic("field values are not intended to be comparable")
}

// unimplemented: panics.
func (op *FieldList) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	panic("field values should only be used in record literals")
}
