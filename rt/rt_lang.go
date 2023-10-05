// Code generated by "makeops"; edit at your own risk.
package rt

import "git.sr.ht/~ionous/tapestry/jsn"

const Assignment_Type = "assignment"

var Assignment_Optional_Marshal = Assignment_Marshal

type Assignment_Slot struct{ Value *Assignment }

func (at Assignment_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at Assignment_Slot) GetType() string              { return Assignment_Type }
func (at Assignment_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at Assignment_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(Assignment)
	return
}

func Assignment_Marshal(m jsn.Marshaler, ptr *Assignment) (err error) {
	slot := Assignment_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Assignment_Repeats_Marshal(m, pv)
	}
	return
}

const BoolEval_Type = "bool_eval"

var BoolEval_Optional_Marshal = BoolEval_Marshal

type BoolEval_Slot struct{ Value *BoolEval }

func (at BoolEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at BoolEval_Slot) GetType() string              { return BoolEval_Type }
func (at BoolEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at BoolEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(BoolEval)
	return
}

func BoolEval_Marshal(m jsn.Marshaler, ptr *BoolEval) (err error) {
	slot := BoolEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = BoolEval_Repeats_Marshal(m, pv)
	}
	return
}

const Execute_Type = "execute"

var Execute_Optional_Marshal = Execute_Marshal

type Execute_Slot struct{ Value *Execute }

func (at Execute_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at Execute_Slot) GetType() string              { return Execute_Type }
func (at Execute_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at Execute_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(Execute)
	return
}

func Execute_Marshal(m jsn.Marshaler, ptr *Execute) (err error) {
	slot := Execute_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Execute_Repeats_Marshal(m, pv)
	}
	return
}

const NumListEval_Type = "num_list_eval"

var NumListEval_Optional_Marshal = NumListEval_Marshal

type NumListEval_Slot struct{ Value *NumListEval }

func (at NumListEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at NumListEval_Slot) GetType() string              { return NumListEval_Type }
func (at NumListEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at NumListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(NumListEval)
	return
}

func NumListEval_Marshal(m jsn.Marshaler, ptr *NumListEval) (err error) {
	slot := NumListEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = NumListEval_Repeats_Marshal(m, pv)
	}
	return
}

const NumberEval_Type = "number_eval"

var NumberEval_Optional_Marshal = NumberEval_Marshal

type NumberEval_Slot struct{ Value *NumberEval }

func (at NumberEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at NumberEval_Slot) GetType() string              { return NumberEval_Type }
func (at NumberEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at NumberEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(NumberEval)
	return
}

func NumberEval_Marshal(m jsn.Marshaler, ptr *NumberEval) (err error) {
	slot := NumberEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = NumberEval_Repeats_Marshal(m, pv)
	}
	return
}

const RecordEval_Type = "record_eval"

var RecordEval_Optional_Marshal = RecordEval_Marshal

type RecordEval_Slot struct{ Value *RecordEval }

func (at RecordEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at RecordEval_Slot) GetType() string              { return RecordEval_Type }
func (at RecordEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at RecordEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(RecordEval)
	return
}

func RecordEval_Marshal(m jsn.Marshaler, ptr *RecordEval) (err error) {
	slot := RecordEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RecordEval_Repeats_Marshal(m, pv)
	}
	return
}

const RecordListEval_Type = "record_list_eval"

var RecordListEval_Optional_Marshal = RecordListEval_Marshal

type RecordListEval_Slot struct{ Value *RecordListEval }

func (at RecordListEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at RecordListEval_Slot) GetType() string              { return RecordListEval_Type }
func (at RecordListEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at RecordListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(RecordListEval)
	return
}

func RecordListEval_Marshal(m jsn.Marshaler, ptr *RecordListEval) (err error) {
	slot := RecordListEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RecordListEval_Repeats_Marshal(m, pv)
	}
	return
}

const TextEval_Type = "text_eval"

var TextEval_Optional_Marshal = TextEval_Marshal

type TextEval_Slot struct{ Value *TextEval }

func (at TextEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at TextEval_Slot) GetType() string              { return TextEval_Type }
func (at TextEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at TextEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(TextEval)
	return
}

func TextEval_Marshal(m jsn.Marshaler, ptr *TextEval) (err error) {
	slot := TextEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TextEval_Repeats_Marshal(m, pv)
	}
	return
}

const TextListEval_Type = "text_list_eval"

var TextListEval_Optional_Marshal = TextListEval_Marshal

type TextListEval_Slot struct{ Value *TextListEval }

func (at TextListEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at TextListEval_Slot) GetType() string              { return TextListEval_Type }
func (at TextListEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at TextListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(TextListEval)
	return
}

func TextListEval_Marshal(m jsn.Marshaler, ptr *TextListEval) (err error) {
	slot := TextListEval_Slot{ptr}
	return slot.Marshal(m)
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
	if len(*pv) > 0 || !m.IsEncoding() {
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

var Signatures = map[uint64]interface{}{}
