// Slots used to produce and consume values.
package rtti

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// address, a type of slot.
var Zt_Address = typeinfo.Slot{
	Name: "address",
	Markup: map[string]any{
		"comment": "Identifies an object field, local variable, or pattern parameter. Addresses can be read from or written to.",
	},
}

// Holds a single slot.
type Address_Slot struct{ Value Address }

// Implements [typeinfo.Instance] for a single slot.
func (*Address_Slot) TypeInfo() typeinfo.T {
	return &Zt_Address
}

// Holds a slice of slots.
type Address_Slots []Address

// Implements [typeinfo.Instance] for a slice of slots.
func (*Address_Slots) TypeInfo() typeinfo.T {
	return &Zt_Address
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Address_Slots) Repeats() bool {
	return len(*op) > 0
}

// assignment, a type of slot.
var Zt_Assignment = typeinfo.Slot{
	Name: "assignment",
	Markup: map[string]any{
		"comment": []interface{}{"Provides access to values in a generic way.", "For example, when the type of a value isn't known in advance.", "See also package assign. ( ex. [assign.FromBool] )"},
	},
}

// Holds a single slot.
type Assignment_Slot struct{ Value Assignment }

// Implements [typeinfo.Instance] for a single slot.
func (*Assignment_Slot) TypeInfo() typeinfo.T {
	return &Zt_Assignment
}

// Holds a slice of slots.
type Assignment_Slots []Assignment

// Implements [typeinfo.Instance] for a slice of slots.
func (*Assignment_Slots) TypeInfo() typeinfo.T {
	return &Zt_Assignment
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Assignment_Slots) Repeats() bool {
	return len(*op) > 0
}

// bool_eval, a type of slot.
var Zt_BoolEval = typeinfo.Slot{
	Name: "bool_eval",
	Markup: map[string]any{
		"blockly-color": "LOGIC_HUE",
		"comment":       "Commands which return true/false values.",
	},
}

// Holds a single slot.
type BoolEval_Slot struct{ Value BoolEval }

// Implements [typeinfo.Instance] for a single slot.
func (*BoolEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_BoolEval
}

// Holds a slice of slots.
type BoolEval_Slots []BoolEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*BoolEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_BoolEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *BoolEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// execute, a type of slot.
var Zt_Execute = typeinfo.Slot{
	Name: "execute",
	Markup: map[string]any{
		"blockly-color": "PROCEDURES_HUE",
		"blockly-stack": true,
		"comment":       "Commands which don't return a value.",
	},
}

// Holds a single slot.
type Execute_Slot struct{ Value Execute }

// Implements [typeinfo.Instance] for a single slot.
func (*Execute_Slot) TypeInfo() typeinfo.T {
	return &Zt_Execute
}

// Holds a slice of slots.
type Execute_Slots []Execute

// Implements [typeinfo.Instance] for a slice of slots.
func (*Execute_Slots) TypeInfo() typeinfo.T {
	return &Zt_Execute
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Execute_Slots) Repeats() bool {
	return len(*op) > 0
}

// num_list_eval, a type of slot.
var Zt_NumListEval = typeinfo.Slot{
	Name: "num_list_eval",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Commands which return a list of numbers.",
	},
}

// Holds a single slot.
type NumListEval_Slot struct{ Value NumListEval }

// Implements [typeinfo.Instance] for a single slot.
func (*NumListEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_NumListEval
}

// Holds a slice of slots.
type NumListEval_Slots []NumListEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*NumListEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_NumListEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *NumListEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// num_eval, a type of slot.
var Zt_NumEval = typeinfo.Slot{
	Name: "num_eval",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Commands which return a number.",
	},
}

// Holds a single slot.
type NumEval_Slot struct{ Value NumEval }

// Implements [typeinfo.Instance] for a single slot.
func (*NumEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_NumEval
}

// Holds a slice of slots.
type NumEval_Slots []NumEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*NumEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_NumEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *NumEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// text_eval, a type of slot.
var Zt_TextEval = typeinfo.Slot{
	Name: "text_eval",
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Commands which return text.",
	},
}

// Holds a single slot.
type TextEval_Slot struct{ Value TextEval }

// Implements [typeinfo.Instance] for a single slot.
func (*TextEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_TextEval
}

// Holds a slice of slots.
type TextEval_Slots []TextEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*TextEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_TextEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *TextEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// text_list_eval, a type of slot.
var Zt_TextListEval = typeinfo.Slot{
	Name: "text_list_eval",
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Commands which return a list of text.",
	},
}

// Holds a single slot.
type TextListEval_Slot struct{ Value TextListEval }

// Implements [typeinfo.Instance] for a single slot.
func (*TextListEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_TextListEval
}

// Holds a slice of slots.
type TextListEval_Slots []TextListEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*TextListEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_TextListEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *TextListEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// record_eval, a type of slot.
var Zt_RecordEval = typeinfo.Slot{
	Name: "record_eval",
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Commands which return a record.",
	},
}

// Holds a single slot.
type RecordEval_Slot struct{ Value RecordEval }

// Implements [typeinfo.Instance] for a single slot.
func (*RecordEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_RecordEval
}

// Holds a slice of slots.
type RecordEval_Slots []RecordEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*RecordEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_RecordEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *RecordEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// record_list_eval, a type of slot.
var Zt_RecordListEval = typeinfo.Slot{
	Name: "record_list_eval",
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Commands which return a list of records.",
	},
}

// Holds a single slot.
type RecordListEval_Slot struct{ Value RecordListEval }

// Implements [typeinfo.Instance] for a single slot.
func (*RecordListEval_Slot) TypeInfo() typeinfo.T {
	return &Zt_RecordListEval
}

// Holds a slice of slots.
type RecordListEval_Slots []RecordListEval

// Implements [typeinfo.Instance] for a slice of slots.
func (*RecordListEval_Slots) TypeInfo() typeinfo.T {
	return &Zt_RecordListEval
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *RecordListEval_Slots) Repeats() bool {
	return len(*op) > 0
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "rtti",
	Comment: []string{
		"Slots used to produce and consume values.",
	},

	Slot: z_slot_list,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Address,
	&Zt_Assignment,
	&Zt_BoolEval,
	&Zt_Execute,
	&Zt_NumListEval,
	&Zt_NumEval,
	&Zt_TextEval,
	&Zt_TextListEval,
	&Zt_RecordEval,
	&Zt_RecordListEval,
}
