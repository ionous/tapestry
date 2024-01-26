// Code generated by Tapestry; edit at your own risk.
package rtti

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// assignment, a type of slot.
var Zt_Assignment = typeinfo.Slot{
	Name: "assignment",
	Markup: map[string]any{
		"comment": "Reads from evals in a uniform manner for common functions.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type Assignment_Slot struct{ Value Assignment }

// implements typeinfo.Inspector for a single slot.
func (*Assignment_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_Assignment, false
}

// holds a slice of slots
type Assignment_Slots []Assignment

// implements typeinfo.Inspector for a series of slots.
func (*Assignment_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_Assignment, true
}

// bool_eval, a type of slot.
var Zt_BoolEval = typeinfo.Slot{
	Name: "bool_eval",
	Markup: map[string]any{
		"blockly-color": "LOGIC_HUE",
		"comment":       "Statements which return true/false values.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type BoolEval_Slot struct{ Value BoolEval }

// implements typeinfo.Inspector for a single slot.
func (*BoolEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_BoolEval, false
}

// holds a slice of slots
type BoolEval_Slots []BoolEval

// implements typeinfo.Inspector for a series of slots.
func (*BoolEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_BoolEval, true
}

// execute, a type of slot.
var Zt_Execute = typeinfo.Slot{
	Name: "execute",
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
func (*Execute_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_Execute, false
}

// holds a slice of slots
type Execute_Slots []Execute

// implements typeinfo.Inspector for a series of slots.
func (*Execute_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_Execute, true
}

// num_list_eval, a type of slot.
var Zt_NumListEval = typeinfo.Slot{
	Name: "num_list_eval",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Statements which return a list of numbers.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type NumListEval_Slot struct{ Value NumListEval }

// implements typeinfo.Inspector for a single slot.
func (*NumListEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_NumListEval, false
}

// holds a slice of slots
type NumListEval_Slots []NumListEval

// implements typeinfo.Inspector for a series of slots.
func (*NumListEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_NumListEval, true
}

// number_eval, a type of slot.
var Zt_NumberEval = typeinfo.Slot{
	Name: "number_eval",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Statements which return a number.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type NumberEval_Slot struct{ Value NumberEval }

// implements typeinfo.Inspector for a single slot.
func (*NumberEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_NumberEval, false
}

// holds a slice of slots
type NumberEval_Slots []NumberEval

// implements typeinfo.Inspector for a series of slots.
func (*NumberEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_NumberEval, true
}

// text_eval, a type of slot.
var Zt_TextEval = typeinfo.Slot{
	Name: "text_eval",
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Statements which return text.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type TextEval_Slot struct{ Value TextEval }

// implements typeinfo.Inspector for a single slot.
func (*TextEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_TextEval, false
}

// holds a slice of slots
type TextEval_Slots []TextEval

// implements typeinfo.Inspector for a series of slots.
func (*TextEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_TextEval, true
}

// text_list_eval, a type of slot.
var Zt_TextListEval = typeinfo.Slot{
	Name: "text_list_eval",
	Markup: map[string]any{
		"blockly-color": "TEXTS_HUE",
		"comment":       "Statements which return a list of text.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type TextListEval_Slot struct{ Value TextListEval }

// implements typeinfo.Inspector for a single slot.
func (*TextListEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_TextListEval, false
}

// holds a slice of slots
type TextListEval_Slots []TextListEval

// implements typeinfo.Inspector for a series of slots.
func (*TextListEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_TextListEval, true
}

// record_eval, a type of slot.
var Zt_RecordEval = typeinfo.Slot{
	Name: "record_eval",
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Statements which return a record.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type RecordEval_Slot struct{ Value RecordEval }

// implements typeinfo.Inspector for a single slot.
func (*RecordEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_RecordEval, false
}

// holds a slice of slots
type RecordEval_Slots []RecordEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_RecordEval, true
}

// record_list_eval, a type of slot.
var Zt_RecordListEval = typeinfo.Slot{
	Name: "record_list_eval",
	Markup: map[string]any{
		"blockly-color": "LISTS_HUE",
		"comment":       "Statements which return a list of records.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type RecordListEval_Slot struct{ Value RecordListEval }

// implements typeinfo.Inspector for a single slot.
func (*RecordListEval_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_RecordListEval, false
}

// holds a slice of slots
type RecordListEval_Slots []RecordListEval

// implements typeinfo.Inspector for a series of slots.
func (*RecordListEval_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_RecordListEval, true
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "rtti",
	Slot: z_slot_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Assignment,
	&Zt_BoolEval,
	&Zt_Execute,
	&Zt_NumListEval,
	&Zt_NumberEval,
	&Zt_TextEval,
	&Zt_TextListEval,
	&Zt_RecordEval,
	&Zt_RecordListEval,
}
