package core

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// BoolValue specifies a simple true/false value.
type BoolValue struct {
	Bool bool
}

// NumValue specifies a number value.
type NumValue struct {
	Num float64
}

// Text specifies a string value.
type TextValue struct {
	Text string
}

// Lines specifies a potentially multi-line string value.
type Lines struct {
	Lines string
}

// NumList specifies multiple float values.
type NumList struct {
	Values []float64
}

// TextList specifies multiple strings.
type TextList struct {
	Values []string
}

// Compose returns a spec for use by the composer editor.
func (*BoolValue) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "bool",
		Name:  "bool_value",
		Spec:  "{bool}",
		Group: "literals",
		Desc:  "Bool Value: specify an explicit true or false value.",
	}
}

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Bool)
	return
}

// String uses strconv.FormatBool.
func (op *BoolValue) String() string {
	return strconv.FormatBool(op.Bool)
}

func (*NumValue) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "num",
		Name:  "num_value",
		Group: "literals",
		Spec:  "{num:number}",
		Desc:  "Number Value: Specify a particular number.",
	}
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

func (*TextValue) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "txt",
		Name:  "text_value",
		Spec:  "{text}",
		Group: "literals",
		Desc:  "Text Value: specify a small bit of text.",
		Stub:  true,
	}
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

func (*Lines) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "here",
		Name:  "lines_value",
		Spec:  "{lines|quote}",
		Group: "literals",
		Desc:  "Lines Value: specify one or more lines of text.",
	}
}

// String returns the lines.
func (op *Lines) String() string {
	return op.Lines
}

func (*NumList) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "nums",
		Group: "literals",
		Desc:  "Number List: Specify a list of multiple numbers.",
	}
}

func (op *NumList) GetNumList(rt.Runtime) (ret g.Value, _ error) {
	// note: this generates a new slice pointing to the op.Values memory;
	// fix: should this be a copy? or, maybe mark this as read-only
	ret = g.FloatsOf(op.Values)
	return
}

func (*TextList) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "txts",
		Group: "literals",
		Desc:  "Text List: specifies multiple string values.",
		Spec:  "text {values*text|comma-and}",
	}
}

func (op *TextList) GetTextList(rt.Runtime) (ret g.Value, _ error) {
	ret = g.StringsOf(op.Values)
	return
}
