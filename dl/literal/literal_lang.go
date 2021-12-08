// Code generated by "makeops"; edit at your own risk.
package literal

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

// BoolValue Specify an explicit true or false.
// User implements: BoolEval, LiteralValue.
type BoolValue struct {
	Bool bool `if:"label=_,type=bool"`
}

func (*BoolValue) Compose() composer.Spec {
	return composer.Spec{
		Name: BoolValue_Type,
		Uses: composer.Type_Flow,
		Lede: "bool",
	}
}

const BoolValue_Type = "bool_value"

const BoolValue_Field_Bool = "$BOOL"

func (op *BoolValue) Marshal(m jsn.Marshaler) error {
	return BoolValue_Marshal(m, op)
}

type BoolValue_Slice []BoolValue

func (op *BoolValue_Slice) GetType() string { return BoolValue_Type }

func (op *BoolValue_Slice) Marshal(m jsn.Marshaler) error {
	return BoolValue_Repeats_Marshal(m, (*[]BoolValue)(op))
}

func (op *BoolValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *BoolValue_Slice) SetSize(cnt int) {
	var els []BoolValue
	if cnt >= 0 {
		els = make(BoolValue_Slice, cnt)
	}
	(*op) = els
}

func (op *BoolValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return BoolValue_Marshal(m, &(*op)[i])
}

func BoolValue_Repeats_Marshal(m jsn.Marshaler, vals *[]BoolValue) error {
	return jsn.RepeatBlock(m, (*BoolValue_Slice)(vals))
}

func BoolValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]BoolValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = BoolValue_Repeats_Marshal(m, pv)
	}
	return
}

type BoolValue_Flow struct{ ptr *BoolValue }

func (n BoolValue_Flow) GetType() string      { return BoolValue_Type }
func (n BoolValue_Flow) GetLede() string      { return "bool" }
func (n BoolValue_Flow) GetFlow() interface{} { return n.ptr }
func (n BoolValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*BoolValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func BoolValue_Optional_Marshal(m jsn.Marshaler, pv **BoolValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = BoolValue_Marshal(m, *pv)
	} else if !enc {
		var v BoolValue
		if err = BoolValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func BoolValue_Marshal(m jsn.Marshaler, val *BoolValue) (err error) {
	if err = m.MarshalBlock(BoolValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("", BoolValue_Field_Bool)
		if e0 == nil {
			e0 = value.Bool_Unboxed_Marshal(m, &val.Bool)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", BoolValue_Field_Bool))
		}
		m.EndBlock()
	}
	return
}

// FieldValue A fixed value of a record.
type FieldValue struct {
	Name  string       `if:"label=name,type=text"`
	Value LiteralValue `if:"label=value"`
}

func (*FieldValue) Compose() composer.Spec {
	return composer.Spec{
		Name: FieldValue_Type,
		Uses: composer.Type_Flow,
		Lede: "field",
	}
}

const FieldValue_Type = "field_value"

const FieldValue_Field_Name = "$NAME"
const FieldValue_Field_Value = "$VALUE"

func (op *FieldValue) Marshal(m jsn.Marshaler) error {
	return FieldValue_Marshal(m, op)
}

type FieldValue_Slice []FieldValue

func (op *FieldValue_Slice) GetType() string { return FieldValue_Type }

func (op *FieldValue_Slice) Marshal(m jsn.Marshaler) error {
	return FieldValue_Repeats_Marshal(m, (*[]FieldValue)(op))
}

func (op *FieldValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *FieldValue_Slice) SetSize(cnt int) {
	var els []FieldValue
	if cnt >= 0 {
		els = make(FieldValue_Slice, cnt)
	}
	(*op) = els
}

func (op *FieldValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return FieldValue_Marshal(m, &(*op)[i])
}

func FieldValue_Repeats_Marshal(m jsn.Marshaler, vals *[]FieldValue) error {
	return jsn.RepeatBlock(m, (*FieldValue_Slice)(vals))
}

func FieldValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]FieldValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = FieldValue_Repeats_Marshal(m, pv)
	}
	return
}

type FieldValue_Flow struct{ ptr *FieldValue }

func (n FieldValue_Flow) GetType() string      { return FieldValue_Type }
func (n FieldValue_Flow) GetLede() string      { return "field" }
func (n FieldValue_Flow) GetFlow() interface{} { return n.ptr }
func (n FieldValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*FieldValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func FieldValue_Optional_Marshal(m jsn.Marshaler, pv **FieldValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = FieldValue_Marshal(m, *pv)
	} else if !enc {
		var v FieldValue
		if err = FieldValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func FieldValue_Marshal(m jsn.Marshaler, val *FieldValue) (err error) {
	if err = m.MarshalBlock(FieldValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("name", FieldValue_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", FieldValue_Field_Name))
		}
		e1 := m.MarshalKey("value", FieldValue_Field_Value)
		if e1 == nil {
			e1 = LiteralValue_Marshal(m, &val.Value)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", FieldValue_Field_Value))
		}
		m.EndBlock()
	}
	return
}

// FieldValues A series of values all for the same record.
type FieldValues struct {
	Values []FieldValue `if:"label=_"`
}

func (*FieldValues) Compose() composer.Spec {
	return composer.Spec{
		Name: FieldValues_Type,
		Uses: composer.Type_Flow,
		Lede: "fields",
	}
}

const FieldValues_Type = "field_values"

const FieldValues_Field_Values = "$VALUES"

func (op *FieldValues) Marshal(m jsn.Marshaler) error {
	return FieldValues_Marshal(m, op)
}

type FieldValues_Slice []FieldValues

func (op *FieldValues_Slice) GetType() string { return FieldValues_Type }

func (op *FieldValues_Slice) Marshal(m jsn.Marshaler) error {
	return FieldValues_Repeats_Marshal(m, (*[]FieldValues)(op))
}

func (op *FieldValues_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *FieldValues_Slice) SetSize(cnt int) {
	var els []FieldValues
	if cnt >= 0 {
		els = make(FieldValues_Slice, cnt)
	}
	(*op) = els
}

func (op *FieldValues_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return FieldValues_Marshal(m, &(*op)[i])
}

func FieldValues_Repeats_Marshal(m jsn.Marshaler, vals *[]FieldValues) error {
	return jsn.RepeatBlock(m, (*FieldValues_Slice)(vals))
}

func FieldValues_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]FieldValues) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = FieldValues_Repeats_Marshal(m, pv)
	}
	return
}

type FieldValues_Flow struct{ ptr *FieldValues }

func (n FieldValues_Flow) GetType() string      { return FieldValues_Type }
func (n FieldValues_Flow) GetLede() string      { return "fields" }
func (n FieldValues_Flow) GetFlow() interface{} { return n.ptr }
func (n FieldValues_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*FieldValues); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func FieldValues_Optional_Marshal(m jsn.Marshaler, pv **FieldValues) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = FieldValues_Marshal(m, *pv)
	} else if !enc {
		var v FieldValues
		if err = FieldValues_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func FieldValues_Marshal(m jsn.Marshaler, val *FieldValues) (err error) {
	if err = m.MarshalBlock(FieldValues_Flow{val}); err == nil {
		e0 := m.MarshalKey("", FieldValues_Field_Values)
		if e0 == nil {
			e0 = FieldValue_Repeats_Marshal(m, &val.Values)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", FieldValues_Field_Values))
		}
		m.EndBlock()
	}
	return
}

const LiteralValue_Type = "literal_value"

var LiteralValue_Optional_Marshal = LiteralValue_Marshal

type LiteralValue_Slot struct{ ptr *LiteralValue }

func (at LiteralValue_Slot) GetType() string              { return LiteralValue_Type }
func (at LiteralValue_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at LiteralValue_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(LiteralValue)
	return
}

func LiteralValue_Marshal(m jsn.Marshaler, ptr *LiteralValue) (err error) {
	slot := LiteralValue_Slot{ptr}
	if err = m.MarshalBlock(slot); err == nil {
		if a, ok := slot.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}

type LiteralValue_Slice []LiteralValue

func (op *LiteralValue_Slice) GetType() string { return LiteralValue_Type }

func (op *LiteralValue_Slice) Marshal(m jsn.Marshaler) error {
	return LiteralValue_Repeats_Marshal(m, (*[]LiteralValue)(op))
}

func (op *LiteralValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *LiteralValue_Slice) SetSize(cnt int) {
	var els []LiteralValue
	if cnt >= 0 {
		els = make(LiteralValue_Slice, cnt)
	}
	(*op) = els
}

func (op *LiteralValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return LiteralValue_Marshal(m, &(*op)[i])
}

func LiteralValue_Repeats_Marshal(m jsn.Marshaler, vals *[]LiteralValue) error {
	return jsn.RepeatBlock(m, (*LiteralValue_Slice)(vals))
}

func LiteralValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]LiteralValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = LiteralValue_Repeats_Marshal(m, pv)
	}
	return
}

// NumValue Specify a particular number.
// User implements: NumberEval, LiteralValue.
type NumValue struct {
	Num float64 `if:"label=_,type=number"`
}

func (*NumValue) Compose() composer.Spec {
	return composer.Spec{
		Name: NumValue_Type,
		Uses: composer.Type_Flow,
		Lede: "num",
	}
}

const NumValue_Type = "num_value"

const NumValue_Field_Num = "$NUM"

func (op *NumValue) Marshal(m jsn.Marshaler) error {
	return NumValue_Marshal(m, op)
}

type NumValue_Slice []NumValue

func (op *NumValue_Slice) GetType() string { return NumValue_Type }

func (op *NumValue_Slice) Marshal(m jsn.Marshaler) error {
	return NumValue_Repeats_Marshal(m, (*[]NumValue)(op))
}

func (op *NumValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *NumValue_Slice) SetSize(cnt int) {
	var els []NumValue
	if cnt >= 0 {
		els = make(NumValue_Slice, cnt)
	}
	(*op) = els
}

func (op *NumValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return NumValue_Marshal(m, &(*op)[i])
}

func NumValue_Repeats_Marshal(m jsn.Marshaler, vals *[]NumValue) error {
	return jsn.RepeatBlock(m, (*NumValue_Slice)(vals))
}

func NumValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]NumValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = NumValue_Repeats_Marshal(m, pv)
	}
	return
}

type NumValue_Flow struct{ ptr *NumValue }

func (n NumValue_Flow) GetType() string      { return NumValue_Type }
func (n NumValue_Flow) GetLede() string      { return "num" }
func (n NumValue_Flow) GetFlow() interface{} { return n.ptr }
func (n NumValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*NumValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func NumValue_Optional_Marshal(m jsn.Marshaler, pv **NumValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = NumValue_Marshal(m, *pv)
	} else if !enc {
		var v NumValue
		if err = NumValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func NumValue_Marshal(m jsn.Marshaler, val *NumValue) (err error) {
	if err = m.MarshalBlock(NumValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("", NumValue_Field_Num)
		if e0 == nil {
			e0 = value.Number_Unboxed_Marshal(m, &val.Num)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", NumValue_Field_Num))
		}
		m.EndBlock()
	}
	return
}

// NumValues Number List: Specify a list of numbers.
// User implements: NumListEval, LiteralValue.
type NumValues struct {
	Values []float64 `if:"label=_,type=number"`
}

func (*NumValues) Compose() composer.Spec {
	return composer.Spec{
		Name: NumValues_Type,
		Uses: composer.Type_Flow,
		Lede: "nums",
	}
}

const NumValues_Type = "number_values"

const NumValues_Field_Values = "$VALUES"

func (op *NumValues) Marshal(m jsn.Marshaler) error {
	return NumValues_Marshal(m, op)
}

type NumValues_Slice []NumValues

func (op *NumValues_Slice) GetType() string { return NumValues_Type }

func (op *NumValues_Slice) Marshal(m jsn.Marshaler) error {
	return NumValues_Repeats_Marshal(m, (*[]NumValues)(op))
}

func (op *NumValues_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *NumValues_Slice) SetSize(cnt int) {
	var els []NumValues
	if cnt >= 0 {
		els = make(NumValues_Slice, cnt)
	}
	(*op) = els
}

func (op *NumValues_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return NumValues_Marshal(m, &(*op)[i])
}

func NumValues_Repeats_Marshal(m jsn.Marshaler, vals *[]NumValues) error {
	return jsn.RepeatBlock(m, (*NumValues_Slice)(vals))
}

func NumValues_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]NumValues) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = NumValues_Repeats_Marshal(m, pv)
	}
	return
}

type NumValues_Flow struct{ ptr *NumValues }

func (n NumValues_Flow) GetType() string      { return NumValues_Type }
func (n NumValues_Flow) GetLede() string      { return "nums" }
func (n NumValues_Flow) GetFlow() interface{} { return n.ptr }
func (n NumValues_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*NumValues); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func NumValues_Optional_Marshal(m jsn.Marshaler, pv **NumValues) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = NumValues_Marshal(m, *pv)
	} else if !enc {
		var v NumValues
		if err = NumValues_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func NumValues_Marshal(m jsn.Marshaler, val *NumValues) (err error) {
	if err = m.MarshalBlock(NumValues_Flow{val}); err == nil {
		e0 := m.MarshalKey("", NumValues_Field_Values)
		if e0 == nil {
			e0 = value.Number_Unboxed_Repeats_Marshal(m, &val.Values)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", NumValues_Field_Values))
		}
		m.EndBlock()
	}
	return
}

// RecordValue Specify a record composed of literal values.
// User implements: RecordEval, LiteralValue.
type RecordValue struct {
	Kind   string      `if:"label=_,type=text"`
	Fields *FieldValue `if:"label=fields,optional"`
}

func (*RecordValue) Compose() composer.Spec {
	return composer.Spec{
		Name: RecordValue_Type,
		Uses: composer.Type_Flow,
		Lede: "rec",
	}
}

const RecordValue_Type = "record_value"

const RecordValue_Field_Kind = "$KIND"
const RecordValue_Field_Fields = "$FIELDS"

func (op *RecordValue) Marshal(m jsn.Marshaler) error {
	return RecordValue_Marshal(m, op)
}

type RecordValue_Slice []RecordValue

func (op *RecordValue_Slice) GetType() string { return RecordValue_Type }

func (op *RecordValue_Slice) Marshal(m jsn.Marshaler) error {
	return RecordValue_Repeats_Marshal(m, (*[]RecordValue)(op))
}

func (op *RecordValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RecordValue_Slice) SetSize(cnt int) {
	var els []RecordValue
	if cnt >= 0 {
		els = make(RecordValue_Slice, cnt)
	}
	(*op) = els
}

func (op *RecordValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RecordValue_Marshal(m, &(*op)[i])
}

func RecordValue_Repeats_Marshal(m jsn.Marshaler, vals *[]RecordValue) error {
	return jsn.RepeatBlock(m, (*RecordValue_Slice)(vals))
}

func RecordValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RecordValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = RecordValue_Repeats_Marshal(m, pv)
	}
	return
}

type RecordValue_Flow struct{ ptr *RecordValue }

func (n RecordValue_Flow) GetType() string      { return RecordValue_Type }
func (n RecordValue_Flow) GetLede() string      { return "rec" }
func (n RecordValue_Flow) GetFlow() interface{} { return n.ptr }
func (n RecordValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RecordValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RecordValue_Optional_Marshal(m jsn.Marshaler, pv **RecordValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RecordValue_Marshal(m, *pv)
	} else if !enc {
		var v RecordValue
		if err = RecordValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RecordValue_Marshal(m jsn.Marshaler, val *RecordValue) (err error) {
	if err = m.MarshalBlock(RecordValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RecordValue_Field_Kind)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Kind)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RecordValue_Field_Kind))
		}
		e1 := m.MarshalKey("fields", RecordValue_Field_Fields)
		if e1 == nil {
			e1 = FieldValue_Optional_Marshal(m, &val.Fields)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RecordValue_Field_Fields))
		}
		m.EndBlock()
	}
	return
}

// RecordValues Specify a series of records, all of the same kind.
// User implements: RecordEval, LiteralValue.
type RecordValues struct {
	Kind       string        `if:"label=_,type=text"`
	Containing []FieldValues `if:"label=containing"`
}

func (*RecordValues) Compose() composer.Spec {
	return composer.Spec{
		Name: RecordValues_Type,
		Uses: composer.Type_Flow,
		Lede: "recs",
	}
}

const RecordValues_Type = "record_values"

const RecordValues_Field_Kind = "$KIND"
const RecordValues_Field_Containing = "$CONTAINING"

func (op *RecordValues) Marshal(m jsn.Marshaler) error {
	return RecordValues_Marshal(m, op)
}

type RecordValues_Slice []RecordValues

func (op *RecordValues_Slice) GetType() string { return RecordValues_Type }

func (op *RecordValues_Slice) Marshal(m jsn.Marshaler) error {
	return RecordValues_Repeats_Marshal(m, (*[]RecordValues)(op))
}

func (op *RecordValues_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RecordValues_Slice) SetSize(cnt int) {
	var els []RecordValues
	if cnt >= 0 {
		els = make(RecordValues_Slice, cnt)
	}
	(*op) = els
}

func (op *RecordValues_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RecordValues_Marshal(m, &(*op)[i])
}

func RecordValues_Repeats_Marshal(m jsn.Marshaler, vals *[]RecordValues) error {
	return jsn.RepeatBlock(m, (*RecordValues_Slice)(vals))
}

func RecordValues_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RecordValues) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = RecordValues_Repeats_Marshal(m, pv)
	}
	return
}

type RecordValues_Flow struct{ ptr *RecordValues }

func (n RecordValues_Flow) GetType() string      { return RecordValues_Type }
func (n RecordValues_Flow) GetLede() string      { return "recs" }
func (n RecordValues_Flow) GetFlow() interface{} { return n.ptr }
func (n RecordValues_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RecordValues); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RecordValues_Optional_Marshal(m jsn.Marshaler, pv **RecordValues) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RecordValues_Marshal(m, *pv)
	} else if !enc {
		var v RecordValues
		if err = RecordValues_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RecordValues_Marshal(m jsn.Marshaler, val *RecordValues) (err error) {
	if err = m.MarshalBlock(RecordValues_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RecordValues_Field_Kind)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Kind)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RecordValues_Field_Kind))
		}
		e1 := m.MarshalKey("containing", RecordValues_Field_Containing)
		if e1 == nil {
			e1 = FieldValues_Repeats_Marshal(m, &val.Containing)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RecordValues_Field_Containing))
		}
		m.EndBlock()
	}
	return
}

// TextValue Specify a small bit of text.
// User implements: TextEval, LiteralValue.
type TextValue struct {
	Text string `if:"label=_,type=text"`
}

func (*TextValue) Compose() composer.Spec {
	return composer.Spec{
		Name: TextValue_Type,
		Uses: composer.Type_Flow,
		Lede: "txt",
	}
}

const TextValue_Type = "text_value"

const TextValue_Field_Text = "$TEXT"

func (op *TextValue) Marshal(m jsn.Marshaler) error {
	return TextValue_Marshal(m, op)
}

type TextValue_Slice []TextValue

func (op *TextValue_Slice) GetType() string { return TextValue_Type }

func (op *TextValue_Slice) Marshal(m jsn.Marshaler) error {
	return TextValue_Repeats_Marshal(m, (*[]TextValue)(op))
}

func (op *TextValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TextValue_Slice) SetSize(cnt int) {
	var els []TextValue
	if cnt >= 0 {
		els = make(TextValue_Slice, cnt)
	}
	(*op) = els
}

func (op *TextValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TextValue_Marshal(m, &(*op)[i])
}

func TextValue_Repeats_Marshal(m jsn.Marshaler, vals *[]TextValue) error {
	return jsn.RepeatBlock(m, (*TextValue_Slice)(vals))
}

func TextValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TextValue) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = TextValue_Repeats_Marshal(m, pv)
	}
	return
}

type TextValue_Flow struct{ ptr *TextValue }

func (n TextValue_Flow) GetType() string      { return TextValue_Type }
func (n TextValue_Flow) GetLede() string      { return "txt" }
func (n TextValue_Flow) GetFlow() interface{} { return n.ptr }
func (n TextValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TextValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TextValue_Optional_Marshal(m jsn.Marshaler, pv **TextValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TextValue_Marshal(m, *pv)
	} else if !enc {
		var v TextValue
		if err = TextValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TextValue_Marshal(m jsn.Marshaler, val *TextValue) (err error) {
	if err = m.MarshalBlock(TextValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("", TextValue_Field_Text)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Text)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TextValue_Field_Text))
		}
		m.EndBlock()
	}
	return
}

// TextValues Text List: Specifies a set of string values.
// User implements: TextListEval, LiteralValue.
type TextValues struct {
	Values []string `if:"label=_,type=text"`
}

func (*TextValues) Compose() composer.Spec {
	return composer.Spec{
		Name: TextValues_Type,
		Uses: composer.Type_Flow,
		Lede: "txts",
	}
}

const TextValues_Type = "text_values"

const TextValues_Field_Values = "$VALUES"

func (op *TextValues) Marshal(m jsn.Marshaler) error {
	return TextValues_Marshal(m, op)
}

type TextValues_Slice []TextValues

func (op *TextValues_Slice) GetType() string { return TextValues_Type }

func (op *TextValues_Slice) Marshal(m jsn.Marshaler) error {
	return TextValues_Repeats_Marshal(m, (*[]TextValues)(op))
}

func (op *TextValues_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TextValues_Slice) SetSize(cnt int) {
	var els []TextValues
	if cnt >= 0 {
		els = make(TextValues_Slice, cnt)
	}
	(*op) = els
}

func (op *TextValues_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TextValues_Marshal(m, &(*op)[i])
}

func TextValues_Repeats_Marshal(m jsn.Marshaler, vals *[]TextValues) error {
	return jsn.RepeatBlock(m, (*TextValues_Slice)(vals))
}

func TextValues_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TextValues) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = TextValues_Repeats_Marshal(m, pv)
	}
	return
}

type TextValues_Flow struct{ ptr *TextValues }

func (n TextValues_Flow) GetType() string      { return TextValues_Type }
func (n TextValues_Flow) GetLede() string      { return "txts" }
func (n TextValues_Flow) GetFlow() interface{} { return n.ptr }
func (n TextValues_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TextValues); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TextValues_Optional_Marshal(m jsn.Marshaler, pv **TextValues) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TextValues_Marshal(m, *pv)
	} else if !enc {
		var v TextValues
		if err = TextValues_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TextValues_Marshal(m jsn.Marshaler, val *TextValues) (err error) {
	if err = m.MarshalBlock(TextValues_Flow{val}); err == nil {
		e0 := m.MarshalKey("", TextValues_Field_Values)
		if e0 == nil {
			e0 = value.Text_Unboxed_Repeats_Marshal(m, &val.Values)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TextValues_Field_Values))
		}
		m.EndBlock()
	}
	return
}

var Slots = []interface{}{
	(*LiteralValue)(nil),
}

var Slats = []composer.Composer{
	(*BoolValue)(nil),
	(*FieldValue)(nil),
	(*FieldValues)(nil),
	(*NumValue)(nil),
	(*NumValues)(nil),
	(*RecordValue)(nil),
	(*RecordValues)(nil),
	(*TextValue)(nil),
	(*TextValues)(nil),
}

var Signatures = map[uint64]interface{}{
	1736897526516691909:  (*BoolValue)(nil),    /* Bool: */
	8393454129252355509:  (*FieldValue)(nil),   /* Field name:value: */
	2198313742266461362:  (*FieldValues)(nil),  /* Fields: */
	9668407916590545547:  (*NumValue)(nil),     /* Num: */
	17428560025310008846: (*NumValues)(nil),    /* Nums: */
	8364493114742675433:  (*RecordValue)(nil),  /* Rec: */
	7274569038616904990:  (*RecordValue)(nil),  /* Rec:fields: */
	5776881376101857802:  (*RecordValues)(nil), /* Recs:containing: */
	15892234395983060559: (*TextValue)(nil),    /* Txt: */
	10570907516103306122: (*TextValues)(nil),   /* Txts: */
}
