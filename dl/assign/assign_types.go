// Code generated by Tapestry; edit at your own risk.
package assign

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// address, a type of slot.
const Z_Address_Name = "address"

var Z_Address_T = typeinfo.Slot{
	Name: Z_Address_Name,
	Markup: map[string]any{
		"comment": []interface{}{"Identifies some particular object field, local variable, or pattern argument.", "Addresses can be read from or written to.", "That is to say, addresses implement all of the rt evals,", "and all commands which read from objects or variables should use the methods the address interface provides."},
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_Address_Slot struct{ Value Address }

// implements typeinfo.Inspector for a single slot.
func (*FIX_Address_Slot) Inspect() typeinfo.T {
	return &Z_Address_T
}

// holds a slice of slots
type Address_Slots []Address

// implements typeinfo.Inspector for a series of slots.
func (*Address_Slots) Inspect() typeinfo.T {
	return &Z_Address_T
}

// dot, a type of slot.
const Z_Dot_Name = "dot"

var Z_Dot_T = typeinfo.Slot{
	Name: Z_Dot_Name,
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Picks values from types containing other values.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_Dot_Slot struct{ Value Dot }

// implements typeinfo.Inspector for a single slot.
func (*FIX_Dot_Slot) Inspect() typeinfo.T {
	return &Z_Dot_T
}

// holds a slice of slots
type Dot_Slots []Dot

// implements typeinfo.Inspector for a series of slots.
func (*Dot_Slots) Inspect() typeinfo.T {
	return &Z_Dot_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SetValue struct {
	Target Address
	Value  rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SetValue) Inspect() typeinfo.T {
	return &Z_SetValue_T
}

// return a valid markup map, creating it if necessary.
func (op *SetValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// set_value, a type of flow.
const Z_SetValue_Name = "set_value"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetValue)(nil)

var Z_SetValue_T = typeinfo.Flow{
	Name: Z_SetValue_Name,
	Lede: "set",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &Z_Address_T,
	}, {
		Name:  "value",
		Label: "value",
		Type:  &rtti.Z_Assignment_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
	Markup: map[string]any{
		"comment": "Store a value into a variable or object.",
	},
}

// holds a slice of type set_value
// FIX: duplicates the spec decl.
type FIX_SetValue_Slice []SetValue

// implements typeinfo.Inspector
func (*SetValue_Slice) Inspect() typeinfo.T {
	return &Z_SetValue_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SetTrait struct {
	Target rtti.TextEval
	Trait  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SetTrait) Inspect() typeinfo.T {
	return &Z_SetTrait_T
}

// return a valid markup map, creating it if necessary.
func (op *SetTrait) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// set_trait, a type of flow.
const Z_SetTrait_Name = "set_trait"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetTrait)(nil)

var Z_SetTrait_T = typeinfo.Flow{
	Name: Z_SetTrait_Name,
	Lede: "set",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:  "trait",
		Label: "trait",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
	Markup: map[string]any{
		"comment": "Set the state of an object.",
	},
}

// holds a slice of type set_trait
// FIX: duplicates the spec decl.
type FIX_SetTrait_Slice []SetTrait

// implements typeinfo.Inspector
func (*SetTrait_Slice) Inspect() typeinfo.T {
	return &Z_SetTrait_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_CopyValue struct {
	Target Address
	Source Address
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*CopyValue) Inspect() typeinfo.T {
	return &Z_CopyValue_T
}

// return a valid markup map, creating it if necessary.
func (op *CopyValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// copy_value, a type of flow.
const Z_CopyValue_Name = "copy_value"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*CopyValue)(nil)

var Z_CopyValue_T = typeinfo.Flow{
	Name: Z_CopyValue_Name,
	Lede: "copy",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &Z_Address_T,
	}, {
		Name:  "source",
		Label: "from",
		Type:  &Z_Address_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Copy from one stored value to another.", "Requires that the type of the two values match exactly"},
	},
}

// holds a slice of type copy_value
// FIX: duplicates the spec decl.
type FIX_CopyValue_Slice []CopyValue

// implements typeinfo.Inspector
func (*CopyValue_Slice) Inspect() typeinfo.T {
	return &Z_CopyValue_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ObjectRef struct {
	Name   rtti.TextEval
	Field  rtti.TextEval
	Dot    Dot
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ObjectRef) Inspect() typeinfo.T {
	return &Z_ObjectRef_T
}

// return a valid markup map, creating it if necessary.
func (op *ObjectRef) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// object_ref, a type of flow.
const Z_ObjectRef_Name = "object_ref"

// ensure the command implements its specified slots:
var _ Address = (*ObjectRef)(nil)
var _ rtti.BoolEval = (*ObjectRef)(nil)
var _ rtti.NumberEval = (*ObjectRef)(nil)
var _ rtti.TextEval = (*ObjectRef)(nil)
var _ rtti.RecordEval = (*ObjectRef)(nil)
var _ rtti.NumListEval = (*ObjectRef)(nil)
var _ rtti.TextListEval = (*ObjectRef)(nil)
var _ rtti.RecordListEval = (*ObjectRef)(nil)

var Z_ObjectRef_T = typeinfo.Flow{
	Name: Z_ObjectRef_Name,
	Lede: "object",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:  "field",
		Label: "field",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:     "dot",
		Label:    "dot",
		Optional: true,
		Repeats:  true,
		Type:     &Z_Dot_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Address_T,
		&rtti.Z_BoolEval_T,
		&rtti.Z_NumberEval_T,
		&rtti.Z_TextEval_T,
		&rtti.Z_RecordEval_T,
		&rtti.Z_NumListEval_T,
		&rtti.Z_TextListEval_T,
		&rtti.Z_RecordListEval_T,
	},
}

// holds a slice of type object_ref
// FIX: duplicates the spec decl.
type FIX_ObjectRef_Slice []ObjectRef

// implements typeinfo.Inspector
func (*ObjectRef_Slice) Inspect() typeinfo.T {
	return &Z_ObjectRef_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_VariableRef struct {
	Name   rtti.TextEval
	Dot    Dot
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*VariableRef) Inspect() typeinfo.T {
	return &Z_VariableRef_T
}

// return a valid markup map, creating it if necessary.
func (op *VariableRef) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// variable_ref, a type of flow.
const Z_VariableRef_Name = "variable_ref"

// ensure the command implements its specified slots:
var _ Address = (*VariableRef)(nil)
var _ rtti.BoolEval = (*VariableRef)(nil)
var _ rtti.NumberEval = (*VariableRef)(nil)
var _ rtti.TextEval = (*VariableRef)(nil)
var _ rtti.RecordEval = (*VariableRef)(nil)
var _ rtti.NumListEval = (*VariableRef)(nil)
var _ rtti.TextListEval = (*VariableRef)(nil)
var _ rtti.RecordListEval = (*VariableRef)(nil)

var Z_VariableRef_T = typeinfo.Flow{
	Name: Z_VariableRef_Name,
	Lede: "variable",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:     "dot",
		Label:    "dot",
		Optional: true,
		Repeats:  true,
		Type:     &Z_Dot_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Address_T,
		&rtti.Z_BoolEval_T,
		&rtti.Z_NumberEval_T,
		&rtti.Z_TextEval_T,
		&rtti.Z_RecordEval_T,
		&rtti.Z_NumListEval_T,
		&rtti.Z_TextListEval_T,
		&rtti.Z_RecordListEval_T,
	},
}

// holds a slice of type variable_ref
// FIX: duplicates the spec decl.
type FIX_VariableRef_Slice []VariableRef

// implements typeinfo.Inspector
func (*VariableRef_Slice) Inspect() typeinfo.T {
	return &Z_VariableRef_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_AtField struct {
	Field  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*AtField) Inspect() typeinfo.T {
	return &Z_AtField_T
}

// return a valid markup map, creating it if necessary.
func (op *AtField) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// at_field, a type of flow.
const Z_AtField_Name = "at_field"

// ensure the command implements its specified slots:
var _ Dot = (*AtField)(nil)

var Z_AtField_T = typeinfo.Flow{
	Name: Z_AtField_Name,
	Lede: "at_field",
	Terms: []typeinfo.Term{{
		Name:  "field",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Dot_T,
	},
}

// holds a slice of type at_field
// FIX: duplicates the spec decl.
type FIX_AtField_Slice []AtField

// implements typeinfo.Inspector
func (*AtField_Slice) Inspect() typeinfo.T {
	return &Z_AtField_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_AtIndex struct {
	Index  rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*AtIndex) Inspect() typeinfo.T {
	return &Z_AtIndex_T
}

// return a valid markup map, creating it if necessary.
func (op *AtIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// at_index, a type of flow.
const Z_AtIndex_Name = "at_index"

// ensure the command implements its specified slots:
var _ Dot = (*AtIndex)(nil)

var Z_AtIndex_T = typeinfo.Flow{
	Name: Z_AtIndex_Name,
	Lede: "at_index",
	Terms: []typeinfo.Term{{
		Name:  "index",
		Label: "_",
		Type:  &rtti.Z_NumberEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Dot_T,
	},
}

// holds a slice of type at_index
// FIX: duplicates the spec decl.
type FIX_AtIndex_Slice []AtIndex

// implements typeinfo.Inspector
func (*AtIndex_Slice) Inspect() typeinfo.T {
	return &Z_AtIndex_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_CallPattern struct {
	PatternName string
	Arguments   Arg
	Markup      map[string]any
}

// implements typeinfo.Inspector
func (*CallPattern) Inspect() typeinfo.T {
	return &Z_CallPattern_T
}

// return a valid markup map, creating it if necessary.
func (op *CallPattern) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// call_pattern, a type of flow.
const Z_CallPattern_Name = "call_pattern"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*CallPattern)(nil)
var _ rtti.BoolEval = (*CallPattern)(nil)
var _ rtti.NumberEval = (*CallPattern)(nil)
var _ rtti.TextEval = (*CallPattern)(nil)
var _ rtti.RecordEval = (*CallPattern)(nil)
var _ rtti.NumListEval = (*CallPattern)(nil)
var _ rtti.TextListEval = (*CallPattern)(nil)
var _ rtti.RecordListEval = (*CallPattern)(nil)

var Z_CallPattern_T = typeinfo.Flow{
	Name: Z_CallPattern_Name,
	Lede: "determine",
	Terms: []typeinfo.Term{{
		Name:  "pattern_name",
		Label: "_",
		Type:  &prim.Z_Text_T,
	}, {
		Name:    "arguments",
		Label:   "args",
		Repeats: true,
		Type:    &Z_Arg_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
		&rtti.Z_BoolEval_T,
		&rtti.Z_NumberEval_T,
		&rtti.Z_TextEval_T,
		&rtti.Z_RecordEval_T,
		&rtti.Z_NumListEval_T,
		&rtti.Z_TextListEval_T,
		&rtti.Z_RecordListEval_T,
	},
	Markup: map[string]any{
		"comment": "Executes a pattern, and potentially returns a value.",
	},
}

// holds a slice of type call_pattern
// FIX: duplicates the spec decl.
type FIX_CallPattern_Slice []CallPattern

// implements typeinfo.Inspector
func (*CallPattern_Slice) Inspect() typeinfo.T {
	return &Z_CallPattern_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Arg struct {
	Name   string
	Value  rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Arg) Inspect() typeinfo.T {
	return &Z_Arg_T
}

// return a valid markup map, creating it if necessary.
func (op *Arg) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// arg, a type of flow.
const Z_Arg_Name = "arg"

var Z_Arg_T = typeinfo.Flow{
	Name: Z_Arg_Name,
	Lede: "arg",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "value",
		Label: "from",
		Type:  &rtti.Z_Assignment_T,
	}},
	Markup: map[string]any{
		"comment": "Runtime version of argument.",
	},
}

// holds a slice of type arg
// FIX: duplicates the spec decl.
type FIX_Arg_Slice []Arg

// implements typeinfo.Inspector
func (*Arg_Slice) Inspect() typeinfo.T {
	return &Z_Arg_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromExe struct {
	Exe    rtti.Execute
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromExe) Inspect() typeinfo.T {
	return &Z_FromExe_T
}

// return a valid markup map, creating it if necessary.
func (op *FromExe) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_exe, a type of flow.
const Z_FromExe_Name = "from_exe"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromExe)(nil)

var Z_FromExe_T = typeinfo.Flow{
	Name: Z_FromExe_Name,
	Lede: "from_exe",
	Terms: []typeinfo.Term{{
		Name:  "exe",
		Label: "_",
		Type:  &rtti.Z_Execute_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Adapts an execute statement to an assignment.", "Used internally for package shuttle."},
	},
}

// holds a slice of type from_exe
// FIX: duplicates the spec decl.
type FIX_FromExe_Slice []FromExe

// implements typeinfo.Inspector
func (*FromExe_Slice) Inspect() typeinfo.T {
	return &Z_FromExe_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromBool struct {
	Value  rtti.BoolEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromBool) Inspect() typeinfo.T {
	return &Z_FromBool_T
}

// return a valid markup map, creating it if necessary.
func (op *FromBool) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_bool, a type of flow.
const Z_FromBool_Name = "from_bool"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromBool)(nil)

var Z_FromBool_T = typeinfo.Flow{
	Name: Z_FromBool_Name,
	Lede: "from_bool",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_BoolEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a boolean value.",
	},
}

// holds a slice of type from_bool
// FIX: duplicates the spec decl.
type FIX_FromBool_Slice []FromBool

// implements typeinfo.Inspector
func (*FromBool_Slice) Inspect() typeinfo.T {
	return &Z_FromBool_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromNumber struct {
	Value  rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromNumber) Inspect() typeinfo.T {
	return &Z_FromNumber_T
}

// return a valid markup map, creating it if necessary.
func (op *FromNumber) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_number, a type of flow.
const Z_FromNumber_Name = "from_number"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumber)(nil)

var Z_FromNumber_T = typeinfo.Flow{
	Name: Z_FromNumber_Name,
	Lede: "from_number",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_NumberEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a number.",
	},
}

// holds a slice of type from_number
// FIX: duplicates the spec decl.
type FIX_FromNumber_Slice []FromNumber

// implements typeinfo.Inspector
func (*FromNumber_Slice) Inspect() typeinfo.T {
	return &Z_FromNumber_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromText struct {
	Value  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromText) Inspect() typeinfo.T {
	return &Z_FromText_T
}

// return a valid markup map, creating it if necessary.
func (op *FromText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_text, a type of flow.
const Z_FromText_Name = "from_text"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromText)(nil)

var Z_FromText_T = typeinfo.Flow{
	Name: Z_FromText_Name,
	Lede: "from_text",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a text string.",
	},
}

// holds a slice of type from_text
// FIX: duplicates the spec decl.
type FIX_FromText_Slice []FromText

// implements typeinfo.Inspector
func (*FromText_Slice) Inspect() typeinfo.T {
	return &Z_FromText_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromRecord struct {
	Value  rtti.RecordEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromRecord) Inspect() typeinfo.T {
	return &Z_FromRecord_T
}

// return a valid markup map, creating it if necessary.
func (op *FromRecord) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_record, a type of flow.
const Z_FromRecord_Name = "from_record"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecord)(nil)

var Z_FromRecord_T = typeinfo.Flow{
	Name: Z_FromRecord_Name,
	Lede: "from_record",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_RecordEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a record.",
	},
}

// holds a slice of type from_record
// FIX: duplicates the spec decl.
type FIX_FromRecord_Slice []FromRecord

// implements typeinfo.Inspector
func (*FromRecord_Slice) Inspect() typeinfo.T {
	return &Z_FromRecord_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromNumList struct {
	Value  rtti.NumListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromNumList) Inspect() typeinfo.T {
	return &Z_FromNumList_T
}

// return a valid markup map, creating it if necessary.
func (op *FromNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_num_list, a type of flow.
const Z_FromNumList_Name = "from_num_list"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumList)(nil)

var Z_FromNumList_T = typeinfo.Flow{
	Name: Z_FromNumList_Name,
	Lede: "from_num_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_NumListEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of numbers.",
	},
}

// holds a slice of type from_num_list
// FIX: duplicates the spec decl.
type FIX_FromNumList_Slice []FromNumList

// implements typeinfo.Inspector
func (*FromNumList_Slice) Inspect() typeinfo.T {
	return &Z_FromNumList_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromTextList struct {
	Value  rtti.TextListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromTextList) Inspect() typeinfo.T {
	return &Z_FromTextList_T
}

// return a valid markup map, creating it if necessary.
func (op *FromTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_text_list, a type of flow.
const Z_FromTextList_Name = "from_text_list"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromTextList)(nil)

var Z_FromTextList_T = typeinfo.Flow{
	Name: Z_FromTextList_Name,
	Lede: "from_text_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_TextListEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of text strings.",
	},
}

// holds a slice of type from_text_list
// FIX: duplicates the spec decl.
type FIX_FromTextList_Slice []FromTextList

// implements typeinfo.Inspector
func (*FromTextList_Slice) Inspect() typeinfo.T {
	return &Z_FromTextList_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromRecordList struct {
	Value  rtti.RecordListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromRecordList) Inspect() typeinfo.T {
	return &Z_FromRecordList_T
}

// return a valid markup map, creating it if necessary.
func (op *FromRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// from_record_list, a type of flow.
const Z_FromRecordList_Name = "from_record_list"

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecordList)(nil)

var Z_FromRecordList_T = typeinfo.Flow{
	Name: Z_FromRecordList_Name,
	Lede: "from_record_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_RecordListEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Assignment_T,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of records.",
	},
}

// holds a slice of type from_record_list
// FIX: duplicates the spec decl.
type FIX_FromRecordList_Slice []FromRecordList

// implements typeinfo.Inspector
func (*FromRecordList_Slice) Inspect() typeinfo.T {
	return &Z_FromRecordList_T
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "assign",
	Slot: z_slot_list,
	Flow: z_flow_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_Address_T,
	&Z_Dot_T,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_SetValue_T,
	&Z_SetTrait_T,
	&Z_CopyValue_T,
	&Z_ObjectRef_T,
	&Z_VariableRef_T,
	&Z_AtField_T,
	&Z_AtIndex_T,
	&Z_CallPattern_T,
	&Z_Arg_T,
	&Z_FromExe_T,
	&Z_FromBool_T,
	&Z_FromNumber_T,
	&Z_FromText_T,
	&Z_FromRecord_T,
	&Z_FromNumList_T,
	&Z_FromTextList_T,
	&Z_FromRecordList_T,
}