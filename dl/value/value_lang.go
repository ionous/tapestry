// Code generated by "makeops"; edit at your own risk.
package value

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/reader"
	"git.sr.ht/~ionous/iffy/jsn"
)

// Bool requires a user-specified string.
type Bool struct {
	Str string
}

func (op *Bool) String() string {
	return op.Str
}

const Bool_True = "$TRUE"
const Bool_False = "$FALSE"

func (*Bool) Compose() composer.Spec {
	return composer.Spec{
		Name: Bool_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			Bool_True, Bool_False,
		},
		Strings: []string{
			"true", "false",
		},
	}
}

const Bool_Type = "bool"

func (op *Bool) Marshal(n jsn.Marshaler) {
	Bool_Marshal(n, op)
}

type Bool_Unboxed_Slice []bool

func (op *Bool_Unboxed_Slice) GetSize() int    { return len(*op) }
func (op *Bool_Unboxed_Slice) SetSize(cnt int) { (*op) = make(Bool_Unboxed_Slice, cnt) }

func Bool_Unboxed_Repeats_Marshal(n jsn.Marshaler, vals *[]bool) {
	if n.RepeatValues(Bool_Type, (*Bool_Unboxed_Slice)(vals)) {
		for _, el := range *vals {
			Bool_Unboxed_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Bool_Unboxed_Optional_Marshal(n jsn.Marshaler, val *bool) {
	var zero bool
	if enc := n.IsEncoding(); !enc || *val != zero {
		Bool_Unboxed_Marshal(n, val)
	}
}

func Bool_Unboxed_Marshal(n jsn.Marshaler, val *bool) {
	Bool_Marshal(n, &Bool{jsn.BoxBool(val)})
}

func Bool_Optional_Marshal(n jsn.Marshaler, val *Bool) {
	var zero Bool
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		Bool_Marshal(n, val)
	}
}

func Bool_Marshal(n jsn.Marshaler, val *Bool) {
	n.MarshalValue(Bool_Type, jsn.MakeEnum(val, &val.Str))
}

type Bool_Slice []Bool

func (op *Bool_Slice) GetSize() int    { return len(*op) }
func (op *Bool_Slice) SetSize(cnt int) { (*op) = make(Bool_Slice, cnt) }

func Bool_Repeats_Marshal(n jsn.Marshaler, vals *[]Bool) {
	if n.RepeatValues(Bool_Type, (*Bool_Slice)(vals)) {
		for _, el := range *vals {
			Bool_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// Lines requires a user-specified string.
type Lines struct {
	Str string
}

func (op *Lines) String() string {
	return op.Str
}

func (*Lines) Compose() composer.Spec {
	return composer.Spec{
		Name:        Lines_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const Lines_Type = "lines"

func (op *Lines) Marshal(n jsn.Marshaler) {
	Lines_Marshal(n, op)
}

func Lines_Optional_Marshal(n jsn.Marshaler, val *Lines) {
	var zero Lines
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		Lines_Marshal(n, val)
	}
}

func Lines_Marshal(n jsn.Marshaler, val *Lines) {
	n.MarshalValue(Lines_Type, &val.Str)
}

type Lines_Slice []Lines

func (op *Lines_Slice) GetSize() int    { return len(*op) }
func (op *Lines_Slice) SetSize(cnt int) { (*op) = make(Lines_Slice, cnt) }

func Lines_Repeats_Marshal(n jsn.Marshaler, vals *[]Lines) {
	if n.RepeatValues(Lines_Type, (*Lines_Slice)(vals)) {
		for _, el := range *vals {
			Lines_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// Number requires a user-specified string.
type Number struct {
	Num float64
}

func (*Number) Compose() composer.Spec {
	return composer.Spec{
		Name: Number_Type,
		Uses: composer.Type_Num,
	}
}

const Number_Type = "number"

func (op *Number) Marshal(n jsn.Marshaler) {
	Number_Marshal(n, op)
}

type Number_Unboxed_Slice []float64

func (op *Number_Unboxed_Slice) GetSize() int    { return len(*op) }
func (op *Number_Unboxed_Slice) SetSize(cnt int) { (*op) = make(Number_Unboxed_Slice, cnt) }

func Number_Unboxed_Repeats_Marshal(n jsn.Marshaler, vals *[]float64) {
	if n.RepeatValues(Number_Type, (*Number_Unboxed_Slice)(vals)) {
		for _, el := range *vals {
			Number_Unboxed_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Number_Unboxed_Optional_Marshal(n jsn.Marshaler, val *float64) {
	var zero float64
	if enc := n.IsEncoding(); !enc || *val != zero {
		Number_Unboxed_Marshal(n, val)
	}
}

func Number_Unboxed_Marshal(n jsn.Marshaler, val *float64) {
	Number_Marshal(n, &Number{jsn.BoxFloat64(val)})
}

func Number_Optional_Marshal(n jsn.Marshaler, val *Number) {
	var zero Number
	if enc := n.IsEncoding(); !enc || val.Num != zero.Num {
		Number_Marshal(n, val)
	}
}

func Number_Marshal(n jsn.Marshaler, val *Number) {
	n.MarshalValue(Number_Type, &val.Num)
}

type Number_Slice []Number

func (op *Number_Slice) GetSize() int    { return len(*op) }
func (op *Number_Slice) SetSize(cnt int) { (*op) = make(Number_Slice, cnt) }

func Number_Repeats_Marshal(n jsn.Marshaler, vals *[]Number) {
	if n.RepeatValues(Number_Type, (*Number_Slice)(vals)) {
		for _, el := range *vals {
			Number_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// PatternName requires a user-specified string.
type PatternName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PatternName) String() string {
	return op.Str
}

func (*PatternName) Compose() composer.Spec {
	return composer.Spec{
		Name:        PatternName_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const PatternName_Type = "pattern_name"

func (op *PatternName) Marshal(n jsn.Marshaler) {
	PatternName_Marshal(n, op)
}

func PatternName_Optional_Marshal(n jsn.Marshaler, val *PatternName) {
	var zero PatternName
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		PatternName_Marshal(n, val)
	}
}

func PatternName_Marshal(n jsn.Marshaler, val *PatternName) {
	n.SetCursor(val.At.Offset)
	n.MarshalValue(PatternName_Type, &val.Str)
}

type PatternName_Slice []PatternName

func (op *PatternName_Slice) GetSize() int    { return len(*op) }
func (op *PatternName_Slice) SetSize(cnt int) { (*op) = make(PatternName_Slice, cnt) }

func PatternName_Repeats_Marshal(n jsn.Marshaler, vals *[]PatternName) {
	if n.RepeatValues(PatternName_Type, (*PatternName_Slice)(vals)) {
		for _, el := range *vals {
			PatternName_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// RelationName requires a user-specified string.
type RelationName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *RelationName) String() string {
	return op.Str
}

func (*RelationName) Compose() composer.Spec {
	return composer.Spec{
		Name:        RelationName_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const RelationName_Type = "relation_name"

func (op *RelationName) Marshal(n jsn.Marshaler) {
	RelationName_Marshal(n, op)
}

func RelationName_Optional_Marshal(n jsn.Marshaler, val *RelationName) {
	var zero RelationName
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		RelationName_Marshal(n, val)
	}
}

func RelationName_Marshal(n jsn.Marshaler, val *RelationName) {
	n.SetCursor(val.At.Offset)
	n.MarshalValue(RelationName_Type, &val.Str)
}

type RelationName_Slice []RelationName

func (op *RelationName_Slice) GetSize() int    { return len(*op) }
func (op *RelationName_Slice) SetSize(cnt int) { (*op) = make(RelationName_Slice, cnt) }

func RelationName_Repeats_Marshal(n jsn.Marshaler, vals *[]RelationName) {
	if n.RepeatValues(RelationName_Type, (*RelationName_Slice)(vals)) {
		for _, el := range *vals {
			RelationName_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// Text requires a user-specified string.
type Text struct {
	Str string
}

func (op *Text) String() string {
	return op.Str
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:        Text_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const Text_Type = "text"

func (op *Text) Marshal(n jsn.Marshaler) {
	Text_Marshal(n, op)
}

type Text_Unboxed_Slice []string

func (op *Text_Unboxed_Slice) GetSize() int    { return len(*op) }
func (op *Text_Unboxed_Slice) SetSize(cnt int) { (*op) = make(Text_Unboxed_Slice, cnt) }

func Text_Unboxed_Repeats_Marshal(n jsn.Marshaler, vals *[]string) {
	if n.RepeatValues(Text_Type, (*Text_Unboxed_Slice)(vals)) {
		for _, el := range *vals {
			Text_Unboxed_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Text_Unboxed_Optional_Marshal(n jsn.Marshaler, val *string) {
	var zero string
	if enc := n.IsEncoding(); !enc || *val != zero {
		Text_Unboxed_Marshal(n, val)
	}
}

func Text_Unboxed_Marshal(n jsn.Marshaler, val *string) {
	Text_Marshal(n, &Text{jsn.BoxString(val)})
}

func Text_Optional_Marshal(n jsn.Marshaler, val *Text) {
	var zero Text
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		Text_Marshal(n, val)
	}
}

func Text_Marshal_Customized(n jsn.Marshaler, val *Text) {
	n.MarshalValue(Text_Type, &val.Str)
}

type Text_Slice []Text

func (op *Text_Slice) GetSize() int    { return len(*op) }
func (op *Text_Slice) SetSize(cnt int) { (*op) = make(Text_Slice, cnt) }

func Text_Repeats_Marshal(n jsn.Marshaler, vals *[]Text) {
	if n.RepeatValues(Text_Type, (*Text_Slice)(vals)) {
		for _, el := range *vals {
			Text_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// VariableName requires a user-specified string.
type VariableName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *VariableName) String() string {
	return op.Str
}

func (*VariableName) Compose() composer.Spec {
	return composer.Spec{
		Name:        VariableName_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const VariableName_Type = "variable_name"

func (op *VariableName) Marshal(n jsn.Marshaler) {
	VariableName_Marshal(n, op)
}

func VariableName_Optional_Marshal(n jsn.Marshaler, val *VariableName) {
	var zero VariableName
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		VariableName_Marshal(n, val)
	}
}

func VariableName_Marshal(n jsn.Marshaler, val *VariableName) {
	n.SetCursor(val.At.Offset)
	n.MarshalValue(VariableName_Type, &val.Str)
}

type VariableName_Slice []VariableName

func (op *VariableName_Slice) GetSize() int    { return len(*op) }
func (op *VariableName_Slice) SetSize(cnt int) { (*op) = make(VariableName_Slice, cnt) }

func VariableName_Repeats_Marshal(n jsn.Marshaler, vals *[]VariableName) {
	if n.RepeatValues(VariableName_Type, (*VariableName_Slice)(vals)) {
		for _, el := range *vals {
			VariableName_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

var Slats = []composer.Composer{
	(*Bool)(nil),
	(*Lines)(nil),
	(*Number)(nil),
	(*PatternName)(nil),
	(*RelationName)(nil),
	(*Text)(nil),
	(*VariableName)(nil),
}
