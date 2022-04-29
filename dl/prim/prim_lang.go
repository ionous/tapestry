// Code generated by "makeops"; edit at your own risk.
package prim

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn"
)

// Bool requires a predefined string.
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

func (op *Bool) Marshal(m jsn.Marshaler) error {
	return Bool_Marshal(m, op)
}

type Bool_Unboxed_Slice []bool

func (op *Bool_Unboxed_Slice) GetType() string { return Bool_Type }

func (op *Bool_Unboxed_Slice) Marshal(m jsn.Marshaler) error {
	return Bool_Unboxed_Repeats_Marshal(m, (*[]bool)(op))
}

func (op *Bool_Unboxed_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Bool_Unboxed_Slice) SetSize(cnt int) {
	var els []bool
	if cnt >= 0 {
		els = make(Bool_Unboxed_Slice, cnt)
	}
	(*op) = els
}

func (op *Bool_Unboxed_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Bool_Unboxed_Marshal(m, &(*op)[i])
}

func Bool_Unboxed_Repeats_Marshal(m jsn.Marshaler, vals *[]bool) error {
	return jsn.RepeatBlock(m, (*Bool_Unboxed_Slice)(vals))
}

func Bool_Unboxed_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]bool) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Bool_Unboxed_Repeats_Marshal(m, pv)
	}
	return
}

func Bool_Unboxed_Optional_Marshal(m jsn.Marshaler, val *bool) (err error) {
	var zero bool
	if enc := m.IsEncoding(); !enc || *val != zero {
		err = Bool_Unboxed_Marshal(m, val)
	}
	return
}

func Bool_Unboxed_Marshal(m jsn.Marshaler, val *bool) error {
	return m.MarshalValue(Bool_Type, jsn.BoxBool(val))
}

func Bool_Optional_Marshal(m jsn.Marshaler, val *Bool) (err error) {
	var zero Bool
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = Bool_Marshal(m, val)
	}
	return
}

func Bool_Marshal(m jsn.Marshaler, val *Bool) (err error) {
	return m.MarshalValue(Bool_Type, jsn.MakeEnum(val, &val.Str))
}

type Bool_Slice []Bool

func (op *Bool_Slice) GetType() string { return Bool_Type }

func (op *Bool_Slice) Marshal(m jsn.Marshaler) error {
	return Bool_Repeats_Marshal(m, (*[]Bool)(op))
}

func (op *Bool_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Bool_Slice) SetSize(cnt int) {
	var els []Bool
	if cnt >= 0 {
		els = make(Bool_Slice, cnt)
	}
	(*op) = els
}

func (op *Bool_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Bool_Marshal(m, &(*op)[i])
}

func Bool_Repeats_Marshal(m jsn.Marshaler, vals *[]Bool) error {
	return jsn.RepeatBlock(m, (*Bool_Slice)(vals))
}

func Bool_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Bool) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Bool_Repeats_Marshal(m, pv)
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

func (op *Lines) Marshal(m jsn.Marshaler) error {
	return Lines_Marshal(m, op)
}

func Lines_Optional_Marshal(m jsn.Marshaler, val *Lines) (err error) {
	var zero Lines
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = Lines_Marshal(m, val)
	}
	return
}

func Lines_Marshal(m jsn.Marshaler, val *Lines) (err error) {
	return m.MarshalValue(Lines_Type, &val.Str)
}

type Lines_Slice []Lines

func (op *Lines_Slice) GetType() string { return Lines_Type }

func (op *Lines_Slice) Marshal(m jsn.Marshaler) error {
	return Lines_Repeats_Marshal(m, (*[]Lines)(op))
}

func (op *Lines_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Lines_Slice) SetSize(cnt int) {
	var els []Lines
	if cnt >= 0 {
		els = make(Lines_Slice, cnt)
	}
	(*op) = els
}

func (op *Lines_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Lines_Marshal(m, &(*op)[i])
}

func Lines_Repeats_Marshal(m jsn.Marshaler, vals *[]Lines) error {
	return jsn.RepeatBlock(m, (*Lines_Slice)(vals))
}

func Lines_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Lines) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Lines_Repeats_Marshal(m, pv)
	}
	return
}

// Number requires a user-specified number.
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

func (op *Number) Marshal(m jsn.Marshaler) error {
	return Number_Marshal(m, op)
}

type Number_Unboxed_Slice []float64

func (op *Number_Unboxed_Slice) GetType() string { return Number_Type }

func (op *Number_Unboxed_Slice) Marshal(m jsn.Marshaler) error {
	return Number_Unboxed_Repeats_Marshal(m, (*[]float64)(op))
}

func (op *Number_Unboxed_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Number_Unboxed_Slice) SetSize(cnt int) {
	var els []float64
	if cnt >= 0 {
		els = make(Number_Unboxed_Slice, cnt)
	}
	(*op) = els
}

func (op *Number_Unboxed_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Number_Unboxed_Marshal(m, &(*op)[i])
}

func Number_Unboxed_Repeats_Marshal(m jsn.Marshaler, vals *[]float64) error {
	return jsn.RepeatBlock(m, (*Number_Unboxed_Slice)(vals))
}

func Number_Unboxed_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]float64) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Number_Unboxed_Repeats_Marshal(m, pv)
	}
	return
}

func Number_Unboxed_Optional_Marshal(m jsn.Marshaler, val *float64) (err error) {
	var zero float64
	if enc := m.IsEncoding(); !enc || *val != zero {
		err = Number_Unboxed_Marshal(m, val)
	}
	return
}

func Number_Unboxed_Marshal(m jsn.Marshaler, val *float64) error {
	return m.MarshalValue(Number_Type, jsn.BoxFloat64(val))
}

func Number_Optional_Marshal(m jsn.Marshaler, val *Number) (err error) {
	var zero Number
	if enc := m.IsEncoding(); !enc || val.Num != zero.Num {
		err = Number_Marshal(m, val)
	}
	return
}

func Number_Marshal(m jsn.Marshaler, val *Number) (err error) {
	return m.MarshalValue(Number_Type, &val.Num)
}

type Number_Slice []Number

func (op *Number_Slice) GetType() string { return Number_Type }

func (op *Number_Slice) Marshal(m jsn.Marshaler) error {
	return Number_Repeats_Marshal(m, (*[]Number)(op))
}

func (op *Number_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Number_Slice) SetSize(cnt int) {
	var els []Number
	if cnt >= 0 {
		els = make(Number_Slice, cnt)
	}
	(*op) = els
}

func (op *Number_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Number_Marshal(m, &(*op)[i])
}

func Number_Repeats_Marshal(m jsn.Marshaler, vals *[]Number) error {
	return jsn.RepeatBlock(m, (*Number_Slice)(vals))
}

func Number_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Number) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Number_Repeats_Marshal(m, pv)
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

func (op *Text) Marshal(m jsn.Marshaler) error {
	return Text_Marshal(m, op)
}

type Text_Unboxed_Slice []string

func (op *Text_Unboxed_Slice) GetType() string { return Text_Type }

func (op *Text_Unboxed_Slice) Marshal(m jsn.Marshaler) error {
	return Text_Unboxed_Repeats_Marshal(m, (*[]string)(op))
}

func (op *Text_Unboxed_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Text_Unboxed_Slice) SetSize(cnt int) {
	var els []string
	if cnt >= 0 {
		els = make(Text_Unboxed_Slice, cnt)
	}
	(*op) = els
}

func (op *Text_Unboxed_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Text_Unboxed_Marshal(m, &(*op)[i])
}

func Text_Unboxed_Repeats_Marshal(m jsn.Marshaler, vals *[]string) error {
	return jsn.RepeatBlock(m, (*Text_Unboxed_Slice)(vals))
}

func Text_Unboxed_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]string) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Text_Unboxed_Repeats_Marshal(m, pv)
	}
	return
}

func Text_Unboxed_Optional_Marshal(m jsn.Marshaler, val *string) (err error) {
	var zero string
	if enc := m.IsEncoding(); !enc || *val != zero {
		err = Text_Unboxed_Marshal(m, val)
	}
	return
}

func Text_Unboxed_Marshal(m jsn.Marshaler, val *string) error {
	return m.MarshalValue(Text_Type, jsn.BoxString(val))
}

func Text_Optional_Marshal(m jsn.Marshaler, val *Text) (err error) {
	var zero Text
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = Text_Marshal(m, val)
	}
	return
}

func Text_Marshal(m jsn.Marshaler, val *Text) (err error) {
	return m.MarshalValue(Text_Type, &val.Str)
}

type Text_Slice []Text

func (op *Text_Slice) GetType() string { return Text_Type }

func (op *Text_Slice) Marshal(m jsn.Marshaler) error {
	return Text_Repeats_Marshal(m, (*[]Text)(op))
}

func (op *Text_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Text_Slice) SetSize(cnt int) {
	var els []Text
	if cnt >= 0 {
		els = make(Text_Slice, cnt)
	}
	(*op) = els
}

func (op *Text_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Text_Marshal(m, &(*op)[i])
}

func Text_Repeats_Marshal(m jsn.Marshaler, vals *[]Text) error {
	return jsn.RepeatBlock(m, (*Text_Slice)(vals))
}

func Text_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Text) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Text_Repeats_Marshal(m, pv)
	}
	return
}

var Slats = []composer.Composer{
	(*Bool)(nil),
	(*Lines)(nil),
	(*Number)(nil),
	(*Text)(nil),
}

var Signatures = map[uint64]interface{}{}
