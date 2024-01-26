// Code generated by Tapestry; edit at your own risk.
package assign

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// address, a type of slot.
var Zt_Address = typeinfo.Slot{
	Name: "address",
	Markup: map[string]any{
		"comment": []interface{}{"Identifies some particular object field, local variable, or pattern argument.", "Addresses can be read from or written to.", "That is to say, addresses implement all of the rt evals,", "and all commands which read from objects or variables should use the methods the address interface provides."},
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_Address_Slot struct{ Value Address }

// implements typeinfo.Inspector for a single slot.
func (*FIX_Address_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_Address, false
}

// holds a slice of slots
type Address_Slots []Address

// implements typeinfo.Inspector for a series of slots.
func (*Address_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_Address, true
}

// dot, a type of slot.
var Zt_Dot = typeinfo.Slot{
	Name: "dot",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Picks values from types containing other values.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_Dot_Slot struct{ Value Dot }

// implements typeinfo.Inspector for a single slot.
func (*FIX_Dot_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_Dot, false
}

// holds a slice of slots
type Dot_Slots []Dot

// implements typeinfo.Inspector for a series of slots.
func (*Dot_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_Dot, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SetValue struct {
	Target Address
	Value  rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SetValue) Inspect() (typeinfo.T, bool) {
	return &Zt_SetValue, false
}

// return a valid markup map, creating it if necessary.
func (op *SetValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetValue)(nil)

// set_value, a type of flow.
var Zt_SetValue = typeinfo.Flow{
	Name: "set_value",
	Lede: "set",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &Zt_Address,
	}, {
		Name:  "value",
		Label: "value",
		Type:  &rtti.Zt_Assignment,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
	Markup: map[string]any{
		"comment": "Store a value into a variable or object.",
	},
}

// holds a slice of type set_value
// FIX: duplicates the spec decl.
type FIX_SetValue_Slice []SetValue

// implements typeinfo.Inspector
func (*SetValue_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_SetValue, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SetTrait struct {
	Target rtti.TextEval
	Trait  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SetTrait) Inspect() (typeinfo.T, bool) {
	return &Zt_SetTrait, false
}

// return a valid markup map, creating it if necessary.
func (op *SetTrait) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetTrait)(nil)

// set_trait, a type of flow.
var Zt_SetTrait = typeinfo.Flow{
	Name: "set_trait",
	Lede: "set",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &rtti.Zt_TextEval,
	}, {
		Name:  "trait",
		Label: "trait",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
	Markup: map[string]any{
		"comment": "Set the state of an object.",
	},
}

// holds a slice of type set_trait
// FIX: duplicates the spec decl.
type FIX_SetTrait_Slice []SetTrait

// implements typeinfo.Inspector
func (*SetTrait_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_SetTrait, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_CopyValue struct {
	Target Address
	Source Address
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*CopyValue) Inspect() (typeinfo.T, bool) {
	return &Zt_CopyValue, false
}

// return a valid markup map, creating it if necessary.
func (op *CopyValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*CopyValue)(nil)

// copy_value, a type of flow.
var Zt_CopyValue = typeinfo.Flow{
	Name: "copy_value",
	Lede: "copy",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &Zt_Address,
	}, {
		Name:  "source",
		Label: "from",
		Type:  &Zt_Address,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Copy from one stored value to another.", "Requires that the type of the two values match exactly"},
	},
}

// holds a slice of type copy_value
// FIX: duplicates the spec decl.
type FIX_CopyValue_Slice []CopyValue

// implements typeinfo.Inspector
func (*CopyValue_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_CopyValue, true
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
func (*ObjectRef) Inspect() (typeinfo.T, bool) {
	return &Zt_ObjectRef, false
}

// return a valid markup map, creating it if necessary.
func (op *ObjectRef) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Address = (*ObjectRef)(nil)
var _ rtti.BoolEval = (*ObjectRef)(nil)
var _ rtti.NumberEval = (*ObjectRef)(nil)
var _ rtti.TextEval = (*ObjectRef)(nil)
var _ rtti.RecordEval = (*ObjectRef)(nil)
var _ rtti.NumListEval = (*ObjectRef)(nil)
var _ rtti.TextListEval = (*ObjectRef)(nil)
var _ rtti.RecordListEval = (*ObjectRef)(nil)

// object_ref, a type of flow.
var Zt_ObjectRef = typeinfo.Flow{
	Name: "object_ref",
	Lede: "object",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &rtti.Zt_TextEval,
	}, {
		Name:  "field",
		Label: "field",
		Type:  &rtti.Zt_TextEval,
	}, {
		Name:     "dot",
		Label:    "dot",
		Optional: true,
		Repeats:  true,
		Type:     &Zt_Dot,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Address,
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumberEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordListEval,
	},
}

// holds a slice of type object_ref
// FIX: duplicates the spec decl.
type FIX_ObjectRef_Slice []ObjectRef

// implements typeinfo.Inspector
func (*ObjectRef_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_ObjectRef, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_VariableRef struct {
	Name   rtti.TextEval
	Dot    Dot
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*VariableRef) Inspect() (typeinfo.T, bool) {
	return &Zt_VariableRef, false
}

// return a valid markup map, creating it if necessary.
func (op *VariableRef) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Address = (*VariableRef)(nil)
var _ rtti.BoolEval = (*VariableRef)(nil)
var _ rtti.NumberEval = (*VariableRef)(nil)
var _ rtti.TextEval = (*VariableRef)(nil)
var _ rtti.RecordEval = (*VariableRef)(nil)
var _ rtti.NumListEval = (*VariableRef)(nil)
var _ rtti.TextListEval = (*VariableRef)(nil)
var _ rtti.RecordListEval = (*VariableRef)(nil)

// variable_ref, a type of flow.
var Zt_VariableRef = typeinfo.Flow{
	Name: "variable_ref",
	Lede: "variable",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &rtti.Zt_TextEval,
	}, {
		Name:     "dot",
		Label:    "dot",
		Optional: true,
		Repeats:  true,
		Type:     &Zt_Dot,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Address,
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumberEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordListEval,
	},
}

// holds a slice of type variable_ref
// FIX: duplicates the spec decl.
type FIX_VariableRef_Slice []VariableRef

// implements typeinfo.Inspector
func (*VariableRef_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_VariableRef, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_AtField struct {
	Field  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*AtField) Inspect() (typeinfo.T, bool) {
	return &Zt_AtField, false
}

// return a valid markup map, creating it if necessary.
func (op *AtField) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Dot = (*AtField)(nil)

// at_field, a type of flow.
var Zt_AtField = typeinfo.Flow{
	Name: "at_field",
	Lede: "at_field",
	Terms: []typeinfo.Term{{
		Name:  "field",
		Label: "_",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Dot,
	},
}

// holds a slice of type at_field
// FIX: duplicates the spec decl.
type FIX_AtField_Slice []AtField

// implements typeinfo.Inspector
func (*AtField_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_AtField, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_AtIndex struct {
	Index  rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*AtIndex) Inspect() (typeinfo.T, bool) {
	return &Zt_AtIndex, false
}

// return a valid markup map, creating it if necessary.
func (op *AtIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Dot = (*AtIndex)(nil)

// at_index, a type of flow.
var Zt_AtIndex = typeinfo.Flow{
	Name: "at_index",
	Lede: "at_index",
	Terms: []typeinfo.Term{{
		Name:  "index",
		Label: "_",
		Type:  &rtti.Zt_NumberEval,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Dot,
	},
}

// holds a slice of type at_index
// FIX: duplicates the spec decl.
type FIX_AtIndex_Slice []AtIndex

// implements typeinfo.Inspector
func (*AtIndex_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_AtIndex, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_CallPattern struct {
	PatternName string
	Arguments   Arg
	Markup      map[string]any
}

// implements typeinfo.Inspector
func (*CallPattern) Inspect() (typeinfo.T, bool) {
	return &Zt_CallPattern, false
}

// return a valid markup map, creating it if necessary.
func (op *CallPattern) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*CallPattern)(nil)
var _ rtti.BoolEval = (*CallPattern)(nil)
var _ rtti.NumberEval = (*CallPattern)(nil)
var _ rtti.TextEval = (*CallPattern)(nil)
var _ rtti.RecordEval = (*CallPattern)(nil)
var _ rtti.NumListEval = (*CallPattern)(nil)
var _ rtti.TextListEval = (*CallPattern)(nil)
var _ rtti.RecordListEval = (*CallPattern)(nil)

// call_pattern, a type of flow.
var Zt_CallPattern = typeinfo.Flow{
	Name: "call_pattern",
	Lede: "determine",
	Terms: []typeinfo.Term{{
		Name:  "pattern_name",
		Label: "_",
		Type:  &prim.Zt_Text,
	}, {
		Name:    "arguments",
		Label:   "args",
		Repeats: true,
		Type:    &Zt_Arg,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumberEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordListEval,
	},
	Markup: map[string]any{
		"comment": "Executes a pattern, and potentially returns a value.",
	},
}

// holds a slice of type call_pattern
// FIX: duplicates the spec decl.
type FIX_CallPattern_Slice []CallPattern

// implements typeinfo.Inspector
func (*CallPattern_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_CallPattern, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Arg struct {
	Name   string
	Value  rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Arg) Inspect() (typeinfo.T, bool) {
	return &Zt_Arg, false
}

// return a valid markup map, creating it if necessary.
func (op *Arg) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// arg, a type of flow.
var Zt_Arg = typeinfo.Flow{
	Name: "arg",
	Lede: "arg",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "_",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "value",
		Label: "from",
		Type:  &rtti.Zt_Assignment,
	}},
	Markup: map[string]any{
		"comment": "Runtime version of argument.",
	},
}

// holds a slice of type arg
// FIX: duplicates the spec decl.
type FIX_Arg_Slice []Arg

// implements typeinfo.Inspector
func (*Arg_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_Arg, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromExe struct {
	Exe    rtti.Execute
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromExe) Inspect() (typeinfo.T, bool) {
	return &Zt_FromExe, false
}

// return a valid markup map, creating it if necessary.
func (op *FromExe) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromExe)(nil)

// from_exe, a type of flow.
var Zt_FromExe = typeinfo.Flow{
	Name: "from_exe",
	Lede: "from_exe",
	Terms: []typeinfo.Term{{
		Name:  "exe",
		Label: "_",
		Type:  &rtti.Zt_Execute,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Adapts an execute statement to an assignment.", "Used internally for package shuttle."},
	},
}

// holds a slice of type from_exe
// FIX: duplicates the spec decl.
type FIX_FromExe_Slice []FromExe

// implements typeinfo.Inspector
func (*FromExe_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromExe, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromBool struct {
	Value  rtti.BoolEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromBool) Inspect() (typeinfo.T, bool) {
	return &Zt_FromBool, false
}

// return a valid markup map, creating it if necessary.
func (op *FromBool) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromBool)(nil)

// from_bool, a type of flow.
var Zt_FromBool = typeinfo.Flow{
	Name: "from_bool",
	Lede: "from_bool",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_BoolEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a boolean value.",
	},
}

// holds a slice of type from_bool
// FIX: duplicates the spec decl.
type FIX_FromBool_Slice []FromBool

// implements typeinfo.Inspector
func (*FromBool_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromBool, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromNumber struct {
	Value  rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromNumber) Inspect() (typeinfo.T, bool) {
	return &Zt_FromNumber, false
}

// return a valid markup map, creating it if necessary.
func (op *FromNumber) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumber)(nil)

// from_number, a type of flow.
var Zt_FromNumber = typeinfo.Flow{
	Name: "from_number",
	Lede: "from_number",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_NumberEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a number.",
	},
}

// holds a slice of type from_number
// FIX: duplicates the spec decl.
type FIX_FromNumber_Slice []FromNumber

// implements typeinfo.Inspector
func (*FromNumber_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromNumber, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromText struct {
	Value  rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromText) Inspect() (typeinfo.T, bool) {
	return &Zt_FromText, false
}

// return a valid markup map, creating it if necessary.
func (op *FromText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromText)(nil)

// from_text, a type of flow.
var Zt_FromText = typeinfo.Flow{
	Name: "from_text",
	Lede: "from_text",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a text string.",
	},
}

// holds a slice of type from_text
// FIX: duplicates the spec decl.
type FIX_FromText_Slice []FromText

// implements typeinfo.Inspector
func (*FromText_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromText, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromRecord struct {
	Value  rtti.RecordEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromRecord) Inspect() (typeinfo.T, bool) {
	return &Zt_FromRecord, false
}

// return a valid markup map, creating it if necessary.
func (op *FromRecord) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecord)(nil)

// from_record, a type of flow.
var Zt_FromRecord = typeinfo.Flow{
	Name: "from_record",
	Lede: "from_record",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_RecordEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a record.",
	},
}

// holds a slice of type from_record
// FIX: duplicates the spec decl.
type FIX_FromRecord_Slice []FromRecord

// implements typeinfo.Inspector
func (*FromRecord_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromRecord, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromNumList struct {
	Value  rtti.NumListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromNumList) Inspect() (typeinfo.T, bool) {
	return &Zt_FromNumList, false
}

// return a valid markup map, creating it if necessary.
func (op *FromNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumList)(nil)

// from_num_list, a type of flow.
var Zt_FromNumList = typeinfo.Flow{
	Name: "from_num_list",
	Lede: "from_num_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_NumListEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of numbers.",
	},
}

// holds a slice of type from_num_list
// FIX: duplicates the spec decl.
type FIX_FromNumList_Slice []FromNumList

// implements typeinfo.Inspector
func (*FromNumList_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromNumList, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromTextList struct {
	Value  rtti.TextListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromTextList) Inspect() (typeinfo.T, bool) {
	return &Zt_FromTextList, false
}

// return a valid markup map, creating it if necessary.
func (op *FromTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromTextList)(nil)

// from_text_list, a type of flow.
var Zt_FromTextList = typeinfo.Flow{
	Name: "from_text_list",
	Lede: "from_text_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_TextListEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of text strings.",
	},
}

// holds a slice of type from_text_list
// FIX: duplicates the spec decl.
type FIX_FromTextList_Slice []FromTextList

// implements typeinfo.Inspector
func (*FromTextList_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromTextList, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FromRecordList struct {
	Value  rtti.RecordListEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FromRecordList) Inspect() (typeinfo.T, bool) {
	return &Zt_FromRecordList, false
}

// return a valid markup map, creating it if necessary.
func (op *FromRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecordList)(nil)

// from_record_list, a type of flow.
var Zt_FromRecordList = typeinfo.Flow{
	Name: "from_record_list",
	Lede: "from_record_list",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Zt_RecordListEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Assignment,
	},
	Markup: map[string]any{
		"comment": "Calculates a list of records.",
	},
}

// holds a slice of type from_record_list
// FIX: duplicates the spec decl.
type FIX_FromRecordList_Slice []FromRecordList

// implements typeinfo.Inspector
func (*FromRecordList_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FromRecordList, true
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name:       "assign",
	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Address,
	&Zt_Dot,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_SetValue,
	&Zt_SetTrait,
	&Zt_CopyValue,
	&Zt_ObjectRef,
	&Zt_VariableRef,
	&Zt_AtField,
	&Zt_AtIndex,
	&Zt_CallPattern,
	&Zt_Arg,
	&Zt_FromExe,
	&Zt_FromBool,
	&Zt_FromNumber,
	&Zt_FromText,
	&Zt_FromRecord,
	&Zt_FromNumList,
	&Zt_FromTextList,
	&Zt_FromRecordList,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]any{
	6291103735245333139:  (*Arg)(nil),            /* Arg:from: */
	1683104564853176068:  (*AtField)(nil),        /* dot=AtField: */
	17908840355303216180: (*AtIndex)(nil),        /* dot=AtIndex: */
	12187184211547847098: (*CopyValue)(nil),      /* execute=Copy:from: */
	5430006510328108403:  (*CallPattern)(nil),    /* bool_eval=Determine:args: */
	11666175118824200195: (*CallPattern)(nil),    /* execute=Determine:args: */
	16219448703619493492: (*CallPattern)(nil),    /* num_list_eval=Determine:args: */
	15584772020364696136: (*CallPattern)(nil),    /* number_eval=Determine:args: */
	13992013847750998452: (*CallPattern)(nil),    /* record_eval=Determine:args: */
	352268441608212603:   (*CallPattern)(nil),    /* record_list_eval=Determine:args: */
	5079530186593846942:  (*CallPattern)(nil),    /* text_eval=Determine:args: */
	13938609641525654217: (*CallPattern)(nil),    /* text_list_eval=Determine:args: */
	16065241269206568079: (*FromBool)(nil),       /* assignment=FromBool: */
	9721304908210135401:  (*FromExe)(nil),        /* assignment=FromExe: */
	15276643347016776669: (*FromNumList)(nil),    /* assignment=FromNumList: */
	10386192108847008240: (*FromNumber)(nil),     /* assignment=FromNumber: */
	8445595699766392240:  (*FromRecord)(nil),     /* assignment=FromRecord: */
	17510952281883199828: (*FromRecordList)(nil), /* assignment=FromRecordList: */
	9783457335751138546:  (*FromText)(nil),       /* assignment=FromText: */
	3267530751198060154:  (*FromTextList)(nil),   /* assignment=FromTextList: */
	683773550166455203:   (*ObjectRef)(nil),      /* address=Object:field: */
	1942271780557121620:  (*ObjectRef)(nil),      /* bool_eval=Object:field: */
	8839776639979820731:  (*ObjectRef)(nil),      /* num_list_eval=Object:field: */
	10918337914011251575: (*ObjectRef)(nil),      /* number_eval=Object:field: */
	2347663618411162107:  (*ObjectRef)(nil),      /* record_eval=Object:field: */
	11613264323388154988: (*ObjectRef)(nil),      /* record_list_eval=Object:field: */
	16935348020531425213: (*ObjectRef)(nil),      /* text_eval=Object:field: */
	7207525564346341058:  (*ObjectRef)(nil),      /* text_list_eval=Object:field: */
	2801199650842020300:  (*ObjectRef)(nil),      /* address=Object:field:dot: */
	5711121365333637715:  (*ObjectRef)(nil),      /* bool_eval=Object:field:dot: */
	1214997628858983108:  (*ObjectRef)(nil),      /* num_list_eval=Object:field:dot: */
	11071357156742037304: (*ObjectRef)(nil),      /* number_eval=Object:field:dot: */
	1517965638051539844:  (*ObjectRef)(nil),      /* record_eval=Object:field:dot: */
	13722223890291796107: (*ObjectRef)(nil),      /* record_list_eval=Object:field:dot: */
	15784348372409109382: (*ObjectRef)(nil),      /* text_eval=Object:field:dot: */
	11516059561048599401: (*ObjectRef)(nil),      /* text_list_eval=Object:field:dot: */
	3109912816783629323:  (*SetTrait)(nil),       /* execute=Set:trait: */
	3912570011939708664:  (*SetValue)(nil),       /* execute=Set:value: */
	13692207992970428220: (*VariableRef)(nil),    /* address=Variable: */
	17908519799628660539: (*VariableRef)(nil),    /* bool_eval=Variable: */
	11022385456290008164: (*VariableRef)(nil),    /* num_list_eval=Variable: */
	14722688844418158720: (*VariableRef)(nil),    /* number_eval=Variable: */
	15906653930217516836: (*VariableRef)(nil),    /* record_eval=Variable: */
	16032903663975260899: (*VariableRef)(nil),    /* record_list_eval=Variable: */
	11181798416019134386: (*VariableRef)(nil),    /* text_eval=Variable: */
	14769776891888769773: (*VariableRef)(nil),    /* text_list_eval=Variable: */
	15966558056732701531: (*VariableRef)(nil),    /* address=Variable:dot: */
	7739360284898038596:  (*VariableRef)(nil),    /* bool_eval=Variable:dot: */
	14012826006150347811: (*VariableRef)(nil),    /* num_list_eval=Variable:dot: */
	2218494529839714071:  (*VariableRef)(nil),    /* number_eval=Variable:dot: */
	3479001804857346403:  (*VariableRef)(nil),    /* record_eval=Variable:dot: */
	11938488787528882828: (*VariableRef)(nil),    /* record_list_eval=Variable:dot: */
	4798713833623285465:  (*VariableRef)(nil),    /* text_eval=Variable:dot: */
	12039638244497140214: (*VariableRef)(nil),    /* text_list_eval=Variable:dot: */
}
