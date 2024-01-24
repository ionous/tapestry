// Code generated by Tapestry; edit at your own risk.
package rtti

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// assignment, a type of slot.
const Z_Assignment_Name = "assignment"

var Z_Assignment_T = typeinfo.Slot{
	Name: Z_Assignment_Name,
	Markup: map[string]any{
		"comment": "Reads from evals in a uniform manner for common functions.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type Assignment_Slot struct{ Value Assignment }

// implements typeinfo.Inspector for a single slot.
func (*Assignment_Slot) Inspect() typeinfo.T {
	return &Z_Assignment_T
}

// holds a slice of slots
type Assignment_Slots []Assignment

// implements typeinfo.Inspector for a series of slots.
func (*Assignment_Slots) Inspect() typeinfo.T {
	return &Z_Assignment_T
}

// bool_eval, a type of slot.
const Z_BoolEval_Name = "bool_eval"

var Z_BoolEval_T = typeinfo.Slot{
	Name: Z_BoolEval_Name,
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
	return &Z_BoolEval_T
}

// holds a slice of slots
type BoolEval_Slots []BoolEval

// implements typeinfo.Inspector for a series of slots.
func (*BoolEval_Slots) Inspect() typeinfo.T {
	return &Z_BoolEval_T
}

// execute, a type of slot.
const Z_Execute_Name = "execute"

var Z_Execute_T = typeinfo.Slot{
	Name: Z_Execute_Name,
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
	return &Z_Execute_T
}

// holds a slice of slots
type Execute_Slots []Execute

// implements typeinfo.Inspector for a series of slots.
func (*Execute_Slots) Inspect() typeinfo.T {
	return &Z_Execute_T
}

// num_list_eval, a type of slot.
const Z_NumListEval_Name = "num_list_eval"

var Z_NumListEval_T = typeinfo.Slot{
	Name: Z_NumListEval_Name,
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
	return &Z_NumListEval_T
}

// holds a slice of slots
type NumListEval_Slots []NumListEval

// implements typeinfo.Inspector for a series of slots.
func (*NumListEval_Slots) Inspect() typeinfo.T {
	return &Z_NumListEval_T
}

// number_eval, a type of slot.
const Z_NumberEval_Name = "number_eval"

var Z_NumberEval_T = typeinfo.Slot{
	Name: Z_NumberEval_Name,
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
	return &Z_NumberEval_T
}

// holds a slice of slots
type NumberEval_Slots []NumberEval

// implements typeinfo.Inspector for a series of slots.
func (*NumberEval_Slots) Inspect() typeinfo.T {
	return &Z_NumberEval_T
}

// text_eval, a type of slot.
const Z_TextEval_Name = "text_eval"

var Z_TextEval_T = typeinfo.Slot{
	Name: Z_TextEval_Name,
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
	return &Z_TextEval_T
}

// holds a slice of slots
type TextEval_Slots []TextEval

// implements typeinfo.Inspector for a series of slots.
func (*TextEval_Slots) Inspect() typeinfo.T {
	return &Z_TextEval_T
}

// text_list_eval, a type of slot.
const Z_TextListEval_Name = "text_list_eval"

var Z_TextListEval_T = typeinfo.Slot{
	Name: Z_TextListEval_Name,
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
	return &Z_TextListEval_T
}

// holds a slice of slots
type TextListEval_Slots []TextListEval

// implements typeinfo.Inspector for a series of slots.
func (*TextListEval_Slots) Inspect() typeinfo.T {
	return &Z_TextListEval_T
}

// record_eval, a type of slot.
const Z_RecordEval_Name = "record_eval"

var Z_RecordEval_T = typeinfo.Slot{
	Name: Z_RecordEval_Name,
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
	return &Z_RecordEval_T
}

// holds a slice of slots
type RecordEval_Slots []RecordEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordEval_Slots) Inspect() typeinfo.T {
	return &Z_RecordEval_T
}

// record_list_eval, a type of slot.
const Z_RecordListEval_Name = "record_list_eval"

var Z_RecordListEval_T = typeinfo.Slot{
	Name: Z_RecordListEval_Name,
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
	return &Z_RecordListEval_T
}

// holds a slice of slots
type RecordListEval_Slots []RecordListEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordListEval_Slots) Inspect() typeinfo.T {
	return &Z_RecordListEval_T
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "rtti",
	Slot: z_slot_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_Assignment_T,
	&Z_BoolEval_T,
	&Z_Execute_T,
	&Z_NumListEval_T,
	&Z_NumberEval_T,
	&Z_TextEval_T,
	&Z_TextListEval_T,
	&Z_RecordEval_T,
	&Z_RecordListEval_T,
}