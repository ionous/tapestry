package core

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Bool)
	return
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Bool)
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Num)
	return
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

// GetText implements interface TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringOf(op.Text)
	return
}

// String returns the text.
func (op *TextValue) String() string {
	return op.Text
}

func (op *Numbers) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// fix: would aliasing be better?
	dst := make([]float64, len(op.Values))
	copy(dst, op.Values)
	ret = g.FloatsOf(dst)
	return
}

func (op *Texts) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	// fix: would aliasing be better?
	dst := make([]string, len(op.Values))
	copy(dst, op.Values)
	ret = g.StringsOf(dst)
	return
}
