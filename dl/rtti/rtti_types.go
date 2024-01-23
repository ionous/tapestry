// Code generated by Tapestry; edit at your own risk.
package rtti

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// assignment, a type of slot.
const Z_Assignment_Type = "assignment"

var Z_Assignment_Info = typeinfo.Slot{
	Name: Z_Assignment_Type,
	Markup: map[string]any{
		"comment": "Reads from evals in a uniform manner for common functions.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type Assignment_Slot struct{ Value Assignment }

// implements typeinfo.Inspector for a single slot.
func (*Assignment_Slot) Inspect() typeinfo.T {
	return &Z_Assignment_Info
}

// holds a slice of slots
type Assignment_Slots []Assignment

// implements typeinfo.Inspector for a series of slots.
func (*Assignment_Slots) Inspect() typeinfo.T {
	return &Z_Assignment_Info
}

// bool_eval, a type of slot.
const Z_BoolEval_Type = "bool_eval"

var Z_BoolEval_Info = typeinfo.Slot{
	Name: Z_BoolEval_Type,
	Markup: map[string]any{
		"blockly-color": "LOGIC_HUE",
		"comment":       "Statements which return true/false values.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type BoolEval_Slot struct{ Value BoolEval }

// implements typeinfo.Inspector for a single slot.
func (*BoolEval_Slot) Inspect() typeinfo.T {
	return &Z_BoolEval_Info
}

// holds a slice of slots
type BoolEval_Slots []BoolEval

// implements typeinfo.Inspector for a series of slots.
func (*BoolEval_Slots) Inspect() typeinfo.T {
	return &Z_BoolEval_Info
}

// execute, a type of slot.
const Z_Execute_Type = "execute"

var Z_Execute_Info = typeinfo.Slot{
	Name: Z_Execute_Type,
	Markup: map[string]any{
		"blockly-color": "PROCEDURES_HUE",
		"blockly-stack": true,
		"comment":       "Run a series of statements.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type Execute_Slot struct{ Value Execute }

// implements typeinfo.Inspector for a single slot.
func (*Execute_Slot) Inspect() typeinfo.T {
	return &Z_Execute_Info
}

// holds a slice of slots
type Execute_Slots []Execute

// implements typeinfo.Inspector for a series of slots.
func (*Execute_Slots) Inspect() typeinfo.T {
	return &Z_Execute_Info
}

// num_list_eval, a type of slot.
const Z_NumListEval_Type = "num_list_eval"

var Z_NumListEval_Info = typeinfo.Slot{
	Name: Z_NumListEval_Type,
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Statements which return a list of numbers.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type NumListEval_Slot struct{ Value NumListEval }

// implements typeinfo.Inspector for a single slot.
func (*NumListEval_Slot) Inspect() typeinfo.T {
	return &Z_NumListEval_Info
}

// holds a slice of slots
type NumListEval_Slots []NumListEval

// implements typeinfo.Inspector for a series of slots.
func (*NumListEval_Slots) Inspect() typeinfo.T {
	return &Z_NumListEval_Info
}

// number_eval, a type of slot.
const Z_NumberEval_Type = "number_eval"

var Z_NumberEval_Info = typeinfo.Slot{
	Name: Z_NumberEval_Type,
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Statements which return a number.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type NumberEval_Slot struct{ Value NumberEval }

// implements typeinfo.Inspector for a single slot.
func (*NumberEval_Slot) Inspect() typeinfo.T {
	return &Z_NumberEval_Info
}

// holds a slice of slots
type NumberEval_Slots []NumberEval

// implements typeinfo.Inspector for a series of slots.
func (*NumberEval_Slots) Inspect() typeinfo.T {
	return &Z_NumberEval_Info
}

// text_eval, a type of slot.
const Z_TextEval_Type = "text_eval"

var Z_TextEval_Info = typeinfo.Slot{
	Name: Z_TextEval_Type,
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Statements which return text.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type TextEval_Slot struct{ Value TextEval }

// implements typeinfo.Inspector for a single slot.
func (*TextEval_Slot) Inspect() typeinfo.T {
	return &Z_TextEval_Info
}

// holds a slice of slots
type TextEval_Slots []TextEval

// implements typeinfo.Inspector for a series of slots.
func (*TextEval_Slots) Inspect() typeinfo.T {
	return &Z_TextEval_Info
}

// text_list_eval, a type of slot.
const Z_TextListEval_Type = "text_list_eval"

var Z_TextListEval_Info = typeinfo.Slot{
	Name: Z_TextListEval_Type,
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Statements which return a list of text.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type TextListEval_Slot struct{ Value TextListEval }

// implements typeinfo.Inspector for a single slot.
func (*TextListEval_Slot) Inspect() typeinfo.T {
	return &Z_TextListEval_Info
}

// holds a slice of slots
type TextListEval_Slots []TextListEval

// implements typeinfo.Inspector for a series of slots.
func (*TextListEval_Slots) Inspect() typeinfo.T {
	return &Z_TextListEval_Info
}

// record_eval, a type of slot.
const Z_RecordEval_Type = "record_eval"

var Z_RecordEval_Info = typeinfo.Slot{
	Name: Z_RecordEval_Type,
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Statements which return a record.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type RecordEval_Slot struct{ Value RecordEval }

// implements typeinfo.Inspector for a single slot.
func (*RecordEval_Slot) Inspect() typeinfo.T {
	return &Z_RecordEval_Info
}

// holds a slice of slots
type RecordEval_Slots []RecordEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordEval_Slots) Inspect() typeinfo.T {
	return &Z_RecordEval_Info
}

// record_list_eval, a type of slot.
const Z_RecordListEval_Type = "record_list_eval"

var Z_RecordListEval_Info = typeinfo.Slot{
	Name: Z_RecordListEval_Type,
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Statements which return a list of records.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type RecordListEval_Slot struct{ Value RecordListEval }

// implements typeinfo.Inspector for a single slot.
func (*RecordListEval_Slot) Inspect() typeinfo.T {
	return &Z_RecordListEval_Info
}

// holds a slice of slots
type RecordListEval_Slots []RecordListEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordListEval_Slots) Inspect() typeinfo.T {
	return &Z_RecordListEval_Info
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "rtti",
	Slot: z_slot_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_Assignment_Info,
	&Z_BoolEval_Info,
	&Z_Execute_Info,
	&Z_NumListEval_Info,
	&Z_NumberEval_Info,
	&Z_TextEval_Info,
	&Z_TextListEval_Info,
	&Z_RecordEval_Info,
	&Z_RecordListEval_Info,
}
