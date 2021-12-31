package literal

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// LiteralValue marks script constants.
// ( same interface as rt.Assignment currently )
type LiteralValue interface {
	Affinity() affine.Affinity // alt: use a switch on the eval type to generate affinity.
}

// Affinity returns affine.Bool
func (op *BoolValue) Affinity() affine.Affinity {
	return affine.Bool
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Bool)
}

// GetBool implements rt.BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Bool)
	return
}

// Affinity returns affine.Number
func (op *NumValue) Affinity() affine.Affinity {
	return affine.Number
}

// Int converts to native int.
func (op *NumValue) Int() int {
	return int(op.Num)
}

// Float converts to native float.
func (op *NumValue) Float() float64 {
	return op.Num
}

// String returns a nicely formatted float, with no decimal point when possible.
func (op *NumValue) String() string {
	return strconv.FormatFloat(op.Num, 'g', -1, 64)
}

// GetNumber implements rt.NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Num)
	return
}

// Affinity returns affine.Text
func (op *TextValue) Affinity() affine.Affinity {
	return affine.Text
}

// String returns the text.
func (op *TextValue) String() string {
	return op.Text
}

// GetText implements interface rt.TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringOf(op.Text)
	return
}

// Affinity returns affine.NumList
func (op *NumValues) Affinity() affine.Affinity {
	return affine.NumList
}

// GetNumList implements rt.NumListEval providing the dl with a literal list of numbers.
func (op *NumValues) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// fix: would aliasing be better?
	dst := make([]float64, len(op.Values))
	copy(dst, op.Values)
	ret = g.FloatsOf(dst)
	return
}

// Affinity returns affine.TextList
func (op *TextValues) Affinity() affine.Affinity {
	return affine.TextList
}

// GetTextList implements rt.TextListEval providing the dl with a literal list of text.
func (op *TextValues) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	// fix: would aliasing be better?
	dst := make([]string, len(op.Values))
	copy(dst, op.Values)
	ret = g.StringsOf(dst)
	return
}

// Affinity returns affine.Record
func (op *RecordValue) Affinity() affine.Affinity {
	return affine.Record
}

// GetRecord implements interface rt.RecordEval providing the dl with a text literal.
func (op *RecordValue) GetRecord(run rt.Runtime) (ret g.Value, _ error) {
	panic("not implemented")
	return
}

// Affinity returns affine.RecordList
func (op *RecordValues) Affinity() affine.Affinity {
	return affine.RecordList
}

// GetNumList implements rt.RecordListEval providing the dl with a literal list of records.
func (op *RecordValues) GetRecordList(rt.Runtime) (ret g.Value, _ error) {
	panic("not implemented")
	return
}
