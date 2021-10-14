// Code generated by "makeops"; edit at your own risk.
package rt

import (
	"git.sr.ht/~ionous/iffy/jsn"
)

const Assignment_Type = "assignment"

var Assignment_Optional_Marshal = Assignment_Marshal

type Assignment_Slot struct{ ptr *Assignment }

func (at Assignment_Slot) HasSlot() bool { return at.ptr != nil }
func (at Assignment_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(Assignment)
	return
}

func Assignment_Marshal(n jsn.Marshaler, ptr *Assignment) {
	if ok := n.SlotValues(Assignment_Type, Assignment_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type Assignment_Slice []Assignment

func (op *Assignment_Slice) GetSize() int    { return len(*op) }
func (op *Assignment_Slice) SetSize(cnt int) { (*op) = make(Assignment_Slice, cnt) }

func Assignment_Repeats_Marshal(n jsn.Marshaler, vals *[]Assignment) {
	if n.RepeatValues(Assignment_Type, (*Assignment_Slice)(vals)) {
		for i := range *vals {
			Assignment_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const BoolEval_Type = "bool_eval"

var BoolEval_Optional_Marshal = BoolEval_Marshal

type BoolEval_Slot struct{ ptr *BoolEval }

func (at BoolEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at BoolEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(BoolEval)
	return
}

func BoolEval_Marshal(n jsn.Marshaler, ptr *BoolEval) {
	if ok := n.SlotValues(BoolEval_Type, BoolEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type BoolEval_Slice []BoolEval

func (op *BoolEval_Slice) GetSize() int    { return len(*op) }
func (op *BoolEval_Slice) SetSize(cnt int) { (*op) = make(BoolEval_Slice, cnt) }

func BoolEval_Repeats_Marshal(n jsn.Marshaler, vals *[]BoolEval) {
	if n.RepeatValues(BoolEval_Type, (*BoolEval_Slice)(vals)) {
		for i := range *vals {
			BoolEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const Execute_Type = "execute"

var Execute_Optional_Marshal = Execute_Marshal

type Execute_Slot struct{ ptr *Execute }

func (at Execute_Slot) HasSlot() bool { return at.ptr != nil }
func (at Execute_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(Execute)
	return
}

func Execute_Marshal(n jsn.Marshaler, ptr *Execute) {
	if ok := n.SlotValues(Execute_Type, Execute_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type Execute_Slice []Execute

func (op *Execute_Slice) GetSize() int    { return len(*op) }
func (op *Execute_Slice) SetSize(cnt int) { (*op) = make(Execute_Slice, cnt) }

func Execute_Repeats_Marshal(n jsn.Marshaler, vals *[]Execute) {
	if n.RepeatValues(Execute_Type, (*Execute_Slice)(vals)) {
		for i := range *vals {
			Execute_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const NumListEval_Type = "num_list_eval"

var NumListEval_Optional_Marshal = NumListEval_Marshal

type NumListEval_Slot struct{ ptr *NumListEval }

func (at NumListEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at NumListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(NumListEval)
	return
}

func NumListEval_Marshal(n jsn.Marshaler, ptr *NumListEval) {
	if ok := n.SlotValues(NumListEval_Type, NumListEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type NumListEval_Slice []NumListEval

func (op *NumListEval_Slice) GetSize() int    { return len(*op) }
func (op *NumListEval_Slice) SetSize(cnt int) { (*op) = make(NumListEval_Slice, cnt) }

func NumListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]NumListEval) {
	if n.RepeatValues(NumListEval_Type, (*NumListEval_Slice)(vals)) {
		for i := range *vals {
			NumListEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const NumberEval_Type = "number_eval"

var NumberEval_Optional_Marshal = NumberEval_Marshal

type NumberEval_Slot struct{ ptr *NumberEval }

func (at NumberEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at NumberEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(NumberEval)
	return
}

func NumberEval_Marshal(n jsn.Marshaler, ptr *NumberEval) {
	if ok := n.SlotValues(NumberEval_Type, NumberEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type NumberEval_Slice []NumberEval

func (op *NumberEval_Slice) GetSize() int    { return len(*op) }
func (op *NumberEval_Slice) SetSize(cnt int) { (*op) = make(NumberEval_Slice, cnt) }

func NumberEval_Repeats_Marshal(n jsn.Marshaler, vals *[]NumberEval) {
	if n.RepeatValues(NumberEval_Type, (*NumberEval_Slice)(vals)) {
		for i := range *vals {
			NumberEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const RecordEval_Type = "record_eval"

var RecordEval_Optional_Marshal = RecordEval_Marshal

type RecordEval_Slot struct{ ptr *RecordEval }

func (at RecordEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at RecordEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(RecordEval)
	return
}

func RecordEval_Marshal(n jsn.Marshaler, ptr *RecordEval) {
	if ok := n.SlotValues(RecordEval_Type, RecordEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type RecordEval_Slice []RecordEval

func (op *RecordEval_Slice) GetSize() int    { return len(*op) }
func (op *RecordEval_Slice) SetSize(cnt int) { (*op) = make(RecordEval_Slice, cnt) }

func RecordEval_Repeats_Marshal(n jsn.Marshaler, vals *[]RecordEval) {
	if n.RepeatValues(RecordEval_Type, (*RecordEval_Slice)(vals)) {
		for i := range *vals {
			RecordEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const RecordListEval_Type = "record_list_eval"

var RecordListEval_Optional_Marshal = RecordListEval_Marshal

type RecordListEval_Slot struct{ ptr *RecordListEval }

func (at RecordListEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at RecordListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(RecordListEval)
	return
}

func RecordListEval_Marshal(n jsn.Marshaler, ptr *RecordListEval) {
	if ok := n.SlotValues(RecordListEval_Type, RecordListEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type RecordListEval_Slice []RecordListEval

func (op *RecordListEval_Slice) GetSize() int    { return len(*op) }
func (op *RecordListEval_Slice) SetSize(cnt int) { (*op) = make(RecordListEval_Slice, cnt) }

func RecordListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]RecordListEval) {
	if n.RepeatValues(RecordListEval_Type, (*RecordListEval_Slice)(vals)) {
		for i := range *vals {
			RecordListEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const TextEval_Type = "text_eval"

var TextEval_Optional_Marshal = TextEval_Marshal

type TextEval_Slot struct{ ptr *TextEval }

func (at TextEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at TextEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(TextEval)
	return
}

func TextEval_Marshal(n jsn.Marshaler, ptr *TextEval) {
	if ok := n.SlotValues(TextEval_Type, TextEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type TextEval_Slice []TextEval

func (op *TextEval_Slice) GetSize() int    { return len(*op) }
func (op *TextEval_Slice) SetSize(cnt int) { (*op) = make(TextEval_Slice, cnt) }

func TextEval_Repeats_Marshal(n jsn.Marshaler, vals *[]TextEval) {
	if n.RepeatValues(TextEval_Type, (*TextEval_Slice)(vals)) {
		for i := range *vals {
			TextEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

const TextListEval_Type = "text_list_eval"

var TextListEval_Optional_Marshal = TextListEval_Marshal

type TextListEval_Slot struct{ ptr *TextListEval }

func (at TextListEval_Slot) HasSlot() bool { return at.ptr != nil }
func (at TextListEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(TextListEval)
	return
}

func TextListEval_Marshal(n jsn.Marshaler, ptr *TextListEval) {
	if ok := n.SlotValues(TextListEval_Type, TextListEval_Slot{ptr}); ok {
		(*ptr).(jsn.Marshalee).Marshal(n)
		n.EndValues()
	}
	return
}

type TextListEval_Slice []TextListEval

func (op *TextListEval_Slice) GetSize() int    { return len(*op) }
func (op *TextListEval_Slice) SetSize(cnt int) { (*op) = make(TextListEval_Slice, cnt) }

func TextListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]TextListEval) {
	if n.RepeatValues(TextListEval_Type, (*TextListEval_Slice)(vals)) {
		for i := range *vals {
			TextListEval_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
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

var Signatures = map[uint]interface{}{}
