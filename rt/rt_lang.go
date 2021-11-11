// Code generated by "makeops"; edit at your own risk.
package rt

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

const Assignment_Type = "assignment"

var Assignment_Optional_Marshal = Assignment_Marshal

type Assignment_Slot struct{ ptr *Assignment }

func (at Assignment_Slot) GetType() string              { return Assignment_Type }
func (at Assignment_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at Assignment_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(Assignment)
	return
}

func Assignment_Marshal(m jsn.Marshaler, ptr *Assignment) (err error) {
	slot := Assignment_Slot{ptr}
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

type Assignment_Slice []Assignment

func (op *Assignment_Slice) GetType() string { return Assignment_Type }

func (op *Assignment_Slice) Marshal(m jsn.Marshaler) error {
	return Assignment_Repeats_Marshal(m, (*[]Assignment)(op))
}

func (op *Assignment_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Assignment_Slice) SetSize(cnt int) {
	var els []Assignment
	if cnt >= 0 {
		els = make(Assignment_Slice, cnt)
	}
	(*op) = els
}

func (op *Assignment_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Assignment_Marshal(m, &(*op)[i])
}

func Assignment_Repeats_Marshal(m jsn.Marshaler, vals *[]Assignment) error {
	return jsn.RepeatBlock(m, (*Assignment_Slice)(vals))
}

func Assignment_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Assignment) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Assignment_Repeats_Marshal(m, pv)
	}
	return
}

const BoolEval_Type = "bool_eval"

var BoolEval_Optional_Marshal = BoolEval_Marshal

type BoolEval_Slot struct{ ptr *BoolEval }

func (at BoolEval_Slot) GetType() string              { return BoolEval_Type }
func (at BoolEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at BoolEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(BoolEval)
	return
}

func BoolEval_Marshal(m jsn.Marshaler, ptr *BoolEval) (err error) {
	slot := BoolEval_Slot{ptr}
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

type BoolEval_Slice []BoolEval

func (op *BoolEval_Slice) GetType() string { return BoolEval_Type }

func (op *BoolEval_Slice) Marshal(m jsn.Marshaler) error {
	return BoolEval_Repeats_Marshal(m, (*[]BoolEval)(op))
}

func (op *BoolEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *BoolEval_Slice) SetSize(cnt int) {
	var els []BoolEval
	if cnt >= 0 {
		els = make(BoolEval_Slice, cnt)
	}
	(*op) = els
}

func (op *BoolEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return BoolEval_Marshal(m, &(*op)[i])
}

func BoolEval_Repeats_Marshal(m jsn.Marshaler, vals *[]BoolEval) error {
	return jsn.RepeatBlock(m, (*BoolEval_Slice)(vals))
}

func BoolEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]BoolEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = BoolEval_Repeats_Marshal(m, pv)
	}
	return
}

const Execute_Type = "execute"

var Execute_Optional_Marshal = Execute_Marshal

type Execute_Slot struct{ ptr *Execute }

func (at Execute_Slot) GetType() string              { return Execute_Type }
func (at Execute_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at Execute_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(Execute)
	return
}

func Execute_Marshal(m jsn.Marshaler, ptr *Execute) (err error) {
	slot := Execute_Slot{ptr}
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

type Execute_Slice []Execute

func (op *Execute_Slice) GetType() string { return Execute_Type }

func (op *Execute_Slice) Marshal(m jsn.Marshaler) error {
	return Execute_Repeats_Marshal(m, (*[]Execute)(op))
}

func (op *Execute_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Execute_Slice) SetSize(cnt int) {
	var els []Execute
	if cnt >= 0 {
		els = make(Execute_Slice, cnt)
	}
	(*op) = els
}

func (op *Execute_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Execute_Marshal(m, &(*op)[i])
}

func Execute_Repeats_Marshal(m jsn.Marshaler, vals *[]Execute) error {
	return jsn.RepeatBlock(m, (*Execute_Slice)(vals))
}

func Execute_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Execute) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Execute_Repeats_Marshal(m, pv)
	}
	return
}

// Fragment wrapper to hold assignments to values
type Fragment struct {
	Init Assignment `if:"label=_"`
}

func (*Fragment) Compose() composer.Spec {
	return composer.Spec{
		Name: Fragment_Type,
		Uses: composer.Type_Flow,
		Lede: "assign",
	}
}

const Fragment_Type = "fragment"

const Fragment_Field_Init = "$INIT"

func (op *Fragment) Marshal(m jsn.Marshaler) error {
	return Fragment_Marshal(m, op)
}

type Fragment_Slice []Fragment

func (op *Fragment_Slice) GetType() string { return Fragment_Type }

func (op *Fragment_Slice) Marshal(m jsn.Marshaler) error {
	return Fragment_Repeats_Marshal(m, (*[]Fragment)(op))
}

func (op *Fragment_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Fragment_Slice) SetSize(cnt int) {
	var els []Fragment
	if cnt >= 0 {
		els = make(Fragment_Slice, cnt)
	}
	(*op) = els
}

func (op *Fragment_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Fragment_Marshal(m, &(*op)[i])
}

func Fragment_Repeats_Marshal(m jsn.Marshaler, vals *[]Fragment) error {
	return jsn.RepeatBlock(m, (*Fragment_Slice)(vals))
}

func Fragment_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Fragment) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Fragment_Repeats_Marshal(m, pv)
	}
	return
}

type Fragment_Flow struct{ ptr *Fragment }

func (n Fragment_Flow) GetType() string      { return Fragment_Type }
func (n Fragment_Flow) GetLede() string      { return "assign" }
func (n Fragment_Flow) GetFlow() interface{} { return n.ptr }
func (n Fragment_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Fragment); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Fragment_Optional_Marshal(m jsn.Marshaler, pv **Fragment) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Fragment_Marshal(m, *pv)
	} else if !enc {
		var v Fragment
		if err = Fragment_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Fragment_Marshal(m jsn.Marshaler, val *Fragment) (err error) {
	if err = m.MarshalBlock(Fragment_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Fragment_Field_Init)
		if e0 == nil {
			e0 = Assignment_Marshal(m, &val.Init)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Fragment_Field_Init))
		}
		m.EndBlock()
	}
	return
}

// Handler triggers a series of statements when its filters are satisfied.
type Handler struct {
	Filter BoolEval `if:"label=when,optional"`
	Exe    Execute  `if:"label=do"`
}

func (*Handler) Compose() composer.Spec {
	return composer.Spec{
		Name: Handler_Type,
		Uses: composer.Type_Flow,
		Lede: "handle",
	}
}

const Handler_Type = "handler"

const Handler_Field_Filter = "$FILTER"
const Handler_Field_Exe = "$EXE"

func (op *Handler) Marshal(m jsn.Marshaler) error {
	return Handler_Marshal(m, op)
}

type Handler_Slice []Handler

func (op *Handler_Slice) GetType() string { return Handler_Type }

func (op *Handler_Slice) Marshal(m jsn.Marshaler) error {
	return Handler_Repeats_Marshal(m, (*[]Handler)(op))
}

func (op *Handler_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Handler_Slice) SetSize(cnt int) {
	var els []Handler
	if cnt >= 0 {
		els = make(Handler_Slice, cnt)
	}
	(*op) = els
}

func (op *Handler_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Handler_Marshal(m, &(*op)[i])
}

func Handler_Repeats_Marshal(m jsn.Marshaler, vals *[]Handler) error {
	return jsn.RepeatBlock(m, (*Handler_Slice)(vals))
}

func Handler_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Handler) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Handler_Repeats_Marshal(m, pv)
	}
	return
}

type Handler_Flow struct{ ptr *Handler }

func (n Handler_Flow) GetType() string      { return Handler_Type }
func (n Handler_Flow) GetLede() string      { return "handle" }
func (n Handler_Flow) GetFlow() interface{} { return n.ptr }
func (n Handler_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Handler); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Handler_Optional_Marshal(m jsn.Marshaler, pv **Handler) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Handler_Marshal(m, *pv)
	} else if !enc {
		var v Handler
		if err = Handler_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Handler_Marshal(m jsn.Marshaler, val *Handler) (err error) {
	if err = m.MarshalBlock(Handler_Flow{val}); err == nil {
		e0 := m.MarshalKey("when", Handler_Field_Filter)
		if e0 == nil {
			e0 = BoolEval_Optional_Marshal(m, &val.Filter)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Handler_Field_Filter))
		}
		e1 := m.MarshalKey("do", Handler_Field_Exe)
		if e1 == nil {
			e1 = Execute_Marshal(m, &val.Exe)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", Handler_Field_Exe))
		}
		m.EndBlock()
	}
	return
}

const NumListEval_Type = "num_list_eval"

var NumListEval_Optional_Marshal = NumListEval_Marshal

type NumListEval_Slot struct{ ptr *NumListEval }

func (at NumListEval_Slot) GetType() string              { return NumListEval_Type }
func (at NumListEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at NumListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(NumListEval)
	return
}

func NumListEval_Marshal(m jsn.Marshaler, ptr *NumListEval) (err error) {
	slot := NumListEval_Slot{ptr}
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

type NumListEval_Slice []NumListEval

func (op *NumListEval_Slice) GetType() string { return NumListEval_Type }

func (op *NumListEval_Slice) Marshal(m jsn.Marshaler) error {
	return NumListEval_Repeats_Marshal(m, (*[]NumListEval)(op))
}

func (op *NumListEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *NumListEval_Slice) SetSize(cnt int) {
	var els []NumListEval
	if cnt >= 0 {
		els = make(NumListEval_Slice, cnt)
	}
	(*op) = els
}

func (op *NumListEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return NumListEval_Marshal(m, &(*op)[i])
}

func NumListEval_Repeats_Marshal(m jsn.Marshaler, vals *[]NumListEval) error {
	return jsn.RepeatBlock(m, (*NumListEval_Slice)(vals))
}

func NumListEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]NumListEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = NumListEval_Repeats_Marshal(m, pv)
	}
	return
}

const NumberEval_Type = "number_eval"

var NumberEval_Optional_Marshal = NumberEval_Marshal

type NumberEval_Slot struct{ ptr *NumberEval }

func (at NumberEval_Slot) GetType() string              { return NumberEval_Type }
func (at NumberEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at NumberEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(NumberEval)
	return
}

func NumberEval_Marshal(m jsn.Marshaler, ptr *NumberEval) (err error) {
	slot := NumberEval_Slot{ptr}
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

type NumberEval_Slice []NumberEval

func (op *NumberEval_Slice) GetType() string { return NumberEval_Type }

func (op *NumberEval_Slice) Marshal(m jsn.Marshaler) error {
	return NumberEval_Repeats_Marshal(m, (*[]NumberEval)(op))
}

func (op *NumberEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *NumberEval_Slice) SetSize(cnt int) {
	var els []NumberEval
	if cnt >= 0 {
		els = make(NumberEval_Slice, cnt)
	}
	(*op) = els
}

func (op *NumberEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return NumberEval_Marshal(m, &(*op)[i])
}

func NumberEval_Repeats_Marshal(m jsn.Marshaler, vals *[]NumberEval) error {
	return jsn.RepeatBlock(m, (*NumberEval_Slice)(vals))
}

func NumberEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]NumberEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = NumberEval_Repeats_Marshal(m, pv)
	}
	return
}

const RecordEval_Type = "record_eval"

var RecordEval_Optional_Marshal = RecordEval_Marshal

type RecordEval_Slot struct{ ptr *RecordEval }

func (at RecordEval_Slot) GetType() string              { return RecordEval_Type }
func (at RecordEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at RecordEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(RecordEval)
	return
}

func RecordEval_Marshal(m jsn.Marshaler, ptr *RecordEval) (err error) {
	slot := RecordEval_Slot{ptr}
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

type RecordEval_Slice []RecordEval

func (op *RecordEval_Slice) GetType() string { return RecordEval_Type }

func (op *RecordEval_Slice) Marshal(m jsn.Marshaler) error {
	return RecordEval_Repeats_Marshal(m, (*[]RecordEval)(op))
}

func (op *RecordEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RecordEval_Slice) SetSize(cnt int) {
	var els []RecordEval
	if cnt >= 0 {
		els = make(RecordEval_Slice, cnt)
	}
	(*op) = els
}

func (op *RecordEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RecordEval_Marshal(m, &(*op)[i])
}

func RecordEval_Repeats_Marshal(m jsn.Marshaler, vals *[]RecordEval) error {
	return jsn.RepeatBlock(m, (*RecordEval_Slice)(vals))
}

func RecordEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RecordEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = RecordEval_Repeats_Marshal(m, pv)
	}
	return
}

const RecordListEval_Type = "record_list_eval"

var RecordListEval_Optional_Marshal = RecordListEval_Marshal

type RecordListEval_Slot struct{ ptr *RecordListEval }

func (at RecordListEval_Slot) GetType() string              { return RecordListEval_Type }
func (at RecordListEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at RecordListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(RecordListEval)
	return
}

func RecordListEval_Marshal(m jsn.Marshaler, ptr *RecordListEval) (err error) {
	slot := RecordListEval_Slot{ptr}
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

type RecordListEval_Slice []RecordListEval

func (op *RecordListEval_Slice) GetType() string { return RecordListEval_Type }

func (op *RecordListEval_Slice) Marshal(m jsn.Marshaler) error {
	return RecordListEval_Repeats_Marshal(m, (*[]RecordListEval)(op))
}

func (op *RecordListEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RecordListEval_Slice) SetSize(cnt int) {
	var els []RecordListEval
	if cnt >= 0 {
		els = make(RecordListEval_Slice, cnt)
	}
	(*op) = els
}

func (op *RecordListEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RecordListEval_Marshal(m, &(*op)[i])
}

func RecordListEval_Repeats_Marshal(m jsn.Marshaler, vals *[]RecordListEval) error {
	return jsn.RepeatBlock(m, (*RecordListEval_Slice)(vals))
}

func RecordListEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RecordListEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = RecordListEval_Repeats_Marshal(m, pv)
	}
	return
}

// Rule triggers a named series of statements when its filters and phase are satisfied.
type Rule struct {
	Name     string   `if:"label=_,type=text"`
	RawFlags float64  `if:"label=flags,type=number"`
	Filter   BoolEval `if:"label=when,optional"`
	Execute  Execute  `if:"label=do"`
}

func (*Rule) Compose() composer.Spec {
	return composer.Spec{
		Name: Rule_Type,
		Uses: composer.Type_Flow,
	}
}

const Rule_Type = "rule"

const Rule_Field_Name = "$NAME"
const Rule_Field_RawFlags = "$RAW_FLAGS"
const Rule_Field_Filter = "$FILTER"
const Rule_Field_Execute = "$EXECUTE"

func (op *Rule) Marshal(m jsn.Marshaler) error {
	return Rule_Marshal(m, op)
}

type Rule_Slice []Rule

func (op *Rule_Slice) GetType() string { return Rule_Type }

func (op *Rule_Slice) Marshal(m jsn.Marshaler) error {
	return Rule_Repeats_Marshal(m, (*[]Rule)(op))
}

func (op *Rule_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Rule_Slice) SetSize(cnt int) {
	var els []Rule
	if cnt >= 0 {
		els = make(Rule_Slice, cnt)
	}
	(*op) = els
}

func (op *Rule_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Rule_Marshal(m, &(*op)[i])
}

func Rule_Repeats_Marshal(m jsn.Marshaler, vals *[]Rule) error {
	return jsn.RepeatBlock(m, (*Rule_Slice)(vals))
}

func Rule_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Rule) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Rule_Repeats_Marshal(m, pv)
	}
	return
}

type Rule_Flow struct{ ptr *Rule }

func (n Rule_Flow) GetType() string      { return Rule_Type }
func (n Rule_Flow) GetLede() string      { return Rule_Type }
func (n Rule_Flow) GetFlow() interface{} { return n.ptr }
func (n Rule_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Rule); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Rule_Optional_Marshal(m jsn.Marshaler, pv **Rule) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Rule_Marshal(m, *pv)
	} else if !enc {
		var v Rule
		if err = Rule_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Rule_Marshal(m jsn.Marshaler, val *Rule) (err error) {
	if err = m.MarshalBlock(Rule_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Rule_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Rule_Field_Name))
		}
		e1 := m.MarshalKey("flags", Rule_Field_RawFlags)
		if e1 == nil {
			e1 = value.Number_Unboxed_Marshal(m, &val.RawFlags)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", Rule_Field_RawFlags))
		}
		e2 := m.MarshalKey("when", Rule_Field_Filter)
		if e2 == nil {
			e2 = BoolEval_Optional_Marshal(m, &val.Filter)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", Rule_Field_Filter))
		}
		e3 := m.MarshalKey("do", Rule_Field_Execute)
		if e3 == nil {
			e3 = Execute_Marshal(m, &val.Execute)
		}
		if e3 != nil && e3 != jsn.Missing {
			m.Error(errutil.New(e3, "in flow at", Rule_Field_Execute))
		}
		m.EndBlock()
	}
	return
}

const TextEval_Type = "text_eval"

var TextEval_Optional_Marshal = TextEval_Marshal

type TextEval_Slot struct{ ptr *TextEval }

func (at TextEval_Slot) GetType() string              { return TextEval_Type }
func (at TextEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at TextEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(TextEval)
	return
}

func TextEval_Marshal(m jsn.Marshaler, ptr *TextEval) (err error) {
	slot := TextEval_Slot{ptr}
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

type TextEval_Slice []TextEval

func (op *TextEval_Slice) GetType() string { return TextEval_Type }

func (op *TextEval_Slice) Marshal(m jsn.Marshaler) error {
	return TextEval_Repeats_Marshal(m, (*[]TextEval)(op))
}

func (op *TextEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TextEval_Slice) SetSize(cnt int) {
	var els []TextEval
	if cnt >= 0 {
		els = make(TextEval_Slice, cnt)
	}
	(*op) = els
}

func (op *TextEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TextEval_Marshal(m, &(*op)[i])
}

func TextEval_Repeats_Marshal(m jsn.Marshaler, vals *[]TextEval) error {
	return jsn.RepeatBlock(m, (*TextEval_Slice)(vals))
}

func TextEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TextEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = TextEval_Repeats_Marshal(m, pv)
	}
	return
}

const TextListEval_Type = "text_list_eval"

var TextListEval_Optional_Marshal = TextListEval_Marshal

type TextListEval_Slot struct{ ptr *TextListEval }

func (at TextListEval_Slot) GetType() string              { return TextListEval_Type }
func (at TextListEval_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at TextListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(TextListEval)
	return
}

func TextListEval_Marshal(m jsn.Marshaler, ptr *TextListEval) (err error) {
	slot := TextListEval_Slot{ptr}
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

type TextListEval_Slice []TextListEval

func (op *TextListEval_Slice) GetType() string { return TextListEval_Type }

func (op *TextListEval_Slice) Marshal(m jsn.Marshaler) error {
	return TextListEval_Repeats_Marshal(m, (*[]TextListEval)(op))
}

func (op *TextListEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TextListEval_Slice) SetSize(cnt int) {
	var els []TextListEval
	if cnt >= 0 {
		els = make(TextListEval_Slice, cnt)
	}
	(*op) = els
}

func (op *TextListEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TextListEval_Marshal(m, &(*op)[i])
}

func TextListEval_Repeats_Marshal(m jsn.Marshaler, vals *[]TextListEval) error {
	return jsn.RepeatBlock(m, (*TextListEval_Slice)(vals))
}

func TextListEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TextListEval) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = TextListEval_Repeats_Marshal(m, pv)
	}
	return
}

var Slots = []interface{}{
	(*Assignment)(nil),
	(*BoolEval)(nil),
	(*Execute)(nil),
	(*NumListEval)(nil),
	(*NumberEval)(nil),
	(*RecordEval)(nil),
	(*RecordListEval)(nil),
	(*TextEval)(nil),
	(*TextListEval)(nil),
}

var Slats = []composer.Composer{
	(*Fragment)(nil),
	(*Handler)(nil),
	(*Rule)(nil),
}

var Signatures = map[uint64]interface{}{
	2599336479591478530:  (*Fragment)(nil), /* Assign: */
	1974630509660097686:  (*Handler)(nil),  /* Handle do: */
	471702439204592930:   (*Handler)(nil),  /* Handle when:do: */
	14543757055917034449: (*Rule)(nil),     /* Rule:flags:do: */
	6115960915711515131:  (*Rule)(nil),     /* Rule:flags:when:do: */
}
