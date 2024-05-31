// Common operations on objects, variables, and kinds.
package object

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// dot, a type of slot.
var Zt_Dot = typeinfo.Slot{
	Name: "dot",
	Markup: map[string]any{
		"comment": "Access values inside other values.",
	},
}

// Holds a single slot.
type Dot_Slot struct{ Value Dot }

// Implements [typeinfo.Instance] for a single slot.
func (*Dot_Slot) TypeInfo() typeinfo.T {
	return &Zt_Dot
}

// Holds a slice of slots.
type Dot_Slots []Dot

// Implements [typeinfo.Instance] for a slice of slots.
func (*Dot_Slots) TypeInfo() typeinfo.T {
	return &Zt_Dot
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Dot_Slots) Repeats() bool {
	return len(*op) > 0
}

// Store a value into a variable or object.
// Values are specified as a generic [Assignment].
// The various "From" commands exist to cast specific value types into an assignment.
//
// WARNING: This doesn't convert values from one type to another.
// For example:
//
//	Set:value:
//	- "@some_local_variable"
//	- FromText: "a piece of text to store."
//
// will only work if the local variable can store text. If the variable was declared as a number, the command will generate an error.
type SetValue struct {
	Target rtti.Address
	Value  rtti.Assignment
	Markup map[string]any
}

// set_value, a type of flow.
var Zt_SetValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*SetValue) TypeInfo() typeinfo.T {
	return &Zt_SetValue
}

// Implements [typeinfo.Markup]
func (op *SetValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*SetValue)(nil)

// Holds a slice of type SetValue.
type SetValue_Slice []SetValue

// Implements [typeinfo.Instance] for a slice of SetValue.
func (*SetValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_SetValue
}

// Implements [typeinfo.Repeats] for a slice of SetValue.
func (op *SetValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Set the state of an object or record.
// See also: story `Define state:names:`.
type SetState struct {
	Target rtti.Address
	Trait  rtti.TextEval
	Markup map[string]any
}

// set_state, a type of flow.
var Zt_SetState typeinfo.Flow

// Implements [typeinfo.Instance]
func (*SetState) TypeInfo() typeinfo.T {
	return &Zt_SetState
}

// Implements [typeinfo.Markup]
func (op *SetState) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*SetState)(nil)

// Holds a slice of type SetState.
type SetState_Slice []SetState

// Implements [typeinfo.Instance] for a slice of SetState.
func (*SetState_Slice) TypeInfo() typeinfo.T {
	return &Zt_SetState
}

// Implements [typeinfo.Repeats] for a slice of SetState.
func (op *SetState_Slice) Repeats() bool {
	return len(*op) > 0
}

// Read a value from an object. As a special case, if there are no dot parts, this will return the id of the object.
// In .tell files, this command is often specified with a shortcut. For example:
//
//	"#my_object.some_field"
//
// is a shorter way to say:
//
//	Object:dot:
//	- "my object"
//	- "some field"
//
// WARNING: This doesn't convert values from one type to another. For instance, if a field was declared as text, this will error if read as a boolean.
type ObjectDot struct {
	Name   rtti.TextEval
	Dot    []Dot
	Markup map[string]any
}

// object_dot, a type of flow.
var Zt_ObjectDot typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ObjectDot) TypeInfo() typeinfo.T {
	return &Zt_ObjectDot
}

// Implements [typeinfo.Markup]
func (op *ObjectDot) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Address = (*ObjectDot)(nil)
var _ rtti.BoolEval = (*ObjectDot)(nil)
var _ rtti.NumEval = (*ObjectDot)(nil)
var _ rtti.TextEval = (*ObjectDot)(nil)
var _ rtti.RecordEval = (*ObjectDot)(nil)
var _ rtti.NumListEval = (*ObjectDot)(nil)
var _ rtti.TextListEval = (*ObjectDot)(nil)
var _ rtti.RecordListEval = (*ObjectDot)(nil)

// Holds a slice of type ObjectDot.
type ObjectDot_Slice []ObjectDot

// Implements [typeinfo.Instance] for a slice of ObjectDot.
func (*ObjectDot_Slice) TypeInfo() typeinfo.T {
	return &Zt_ObjectDot
}

// Implements [typeinfo.Repeats] for a slice of ObjectDot.
func (op *ObjectDot_Slice) Repeats() bool {
	return len(*op) > 0
}

// Read a value from a variable.
// In .tell files, this command is often specified with a shortcut. For example:
//
//	"@some_local_variable"
//
// is a shorter way to say:
//
//	Variable:dot: "some local variable"
//
// WARNING: This doesn't convert values from one type to another. For instance, if a field was declared as text, this will error if read as a boolean.
type VariableDot struct {
	Name   rtti.TextEval
	Dot    []Dot
	Markup map[string]any
}

// variable_dot, a type of flow.
var Zt_VariableDot typeinfo.Flow

// Implements [typeinfo.Instance]
func (*VariableDot) TypeInfo() typeinfo.T {
	return &Zt_VariableDot
}

// Implements [typeinfo.Markup]
func (op *VariableDot) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Address = (*VariableDot)(nil)
var _ rtti.BoolEval = (*VariableDot)(nil)
var _ rtti.NumEval = (*VariableDot)(nil)
var _ rtti.TextEval = (*VariableDot)(nil)
var _ rtti.RecordEval = (*VariableDot)(nil)
var _ rtti.NumListEval = (*VariableDot)(nil)
var _ rtti.TextListEval = (*VariableDot)(nil)
var _ rtti.RecordListEval = (*VariableDot)(nil)

// Holds a slice of type VariableDot.
type VariableDot_Slice []VariableDot

// Implements [typeinfo.Instance] for a slice of VariableDot.
func (*VariableDot_Slice) TypeInfo() typeinfo.T {
	return &Zt_VariableDot
}

// Implements [typeinfo.Repeats] for a slice of VariableDot.
func (op *VariableDot_Slice) Repeats() bool {
	return len(*op) > 0
}

// Select a named field from a record, or a named property from an object.
type AtField struct {
	Field  rtti.TextEval
	Markup map[string]any
}

// at_field, a type of flow.
var Zt_AtField typeinfo.Flow

// Implements [typeinfo.Instance]
func (*AtField) TypeInfo() typeinfo.T {
	return &Zt_AtField
}

// Implements [typeinfo.Markup]
func (op *AtField) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Dot = (*AtField)(nil)

// Holds a slice of type AtField.
type AtField_Slice []AtField

// Implements [typeinfo.Instance] for a slice of AtField.
func (*AtField_Slice) TypeInfo() typeinfo.T {
	return &Zt_AtField
}

// Implements [typeinfo.Repeats] for a slice of AtField.
func (op *AtField_Slice) Repeats() bool {
	return len(*op) > 0
}

// Select a value from a list of values.
type AtIndex struct {
	Index  rtti.NumEval
	Markup map[string]any
}

// at_index, a type of flow.
var Zt_AtIndex typeinfo.Flow

// Implements [typeinfo.Instance]
func (*AtIndex) TypeInfo() typeinfo.T {
	return &Zt_AtIndex
}

// Implements [typeinfo.Markup]
func (op *AtIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Dot = (*AtIndex)(nil)

// Holds a slice of type AtIndex.
type AtIndex_Slice []AtIndex

// Implements [typeinfo.Instance] for a slice of AtIndex.
func (*AtIndex_Slice) TypeInfo() typeinfo.T {
	return &Zt_AtIndex
}

// Implements [typeinfo.Repeats] for a slice of AtIndex.
func (op *AtIndex_Slice) Repeats() bool {
	return len(*op) > 0
}

// Full name of the object.
type NameOf struct {
	Object rtti.TextEval
	Markup map[string]any
}

// name_of, a type of flow.
var Zt_NameOf typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NameOf) TypeInfo() typeinfo.T {
	return &Zt_NameOf
}

// Implements [typeinfo.Markup]
func (op *NameOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*NameOf)(nil)

// Holds a slice of type NameOf.
type NameOf_Slice []NameOf

// Implements [typeinfo.Instance] for a slice of NameOf.
func (*NameOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_NameOf
}

// Implements [typeinfo.Repeats] for a slice of NameOf.
func (op *NameOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns all of the object's current traits as a list of text.
type ObjectTraits struct {
	Object rtti.TextEval
	Markup map[string]any
}

// object_traits, a type of flow.
var Zt_ObjectTraits typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ObjectTraits) TypeInfo() typeinfo.T {
	return &Zt_ObjectTraits
}

// Implements [typeinfo.Markup]
func (op *ObjectTraits) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*ObjectTraits)(nil)

// Holds a slice of type ObjectTraits.
type ObjectTraits_Slice []ObjectTraits

// Implements [typeinfo.Instance] for a slice of ObjectTraits.
func (*ObjectTraits_Slice) TypeInfo() typeinfo.T {
	return &Zt_ObjectTraits
}

// Implements [typeinfo.Repeats] for a slice of ObjectTraits.
func (op *ObjectTraits_Slice) Repeats() bool {
	return len(*op) > 0
}

// True if the object is exactly the named kind.
type IsExactKindOf struct {
	Object rtti.TextEval
	Kind   string
	Markup map[string]any
}

// is_exact_kind_of, a type of flow.
var Zt_IsExactKindOf typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsExactKindOf) TypeInfo() typeinfo.T {
	return &Zt_IsExactKindOf
}

// Implements [typeinfo.Markup]
func (op *IsExactKindOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsExactKindOf)(nil)

// Holds a slice of type IsExactKindOf.
type IsExactKindOf_Slice []IsExactKindOf

// Implements [typeinfo.Instance] for a slice of IsExactKindOf.
func (*IsExactKindOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsExactKindOf
}

// Implements [typeinfo.Repeats] for a slice of IsExactKindOf.
func (op *IsExactKindOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// True if the object is compatible with the named kind.
type IsKindOf struct {
	Object  rtti.TextEval
	Kind    string
	Nothing bool
	Markup  map[string]any
}

// is_kind_of, a type of flow.
var Zt_IsKindOf typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsKindOf) TypeInfo() typeinfo.T {
	return &Zt_IsKindOf
}

// Implements [typeinfo.Markup]
func (op *IsKindOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsKindOf)(nil)

// Holds a slice of type IsKindOf.
type IsKindOf_Slice []IsKindOf

// Implements [typeinfo.Instance] for a slice of IsKindOf.
func (*IsKindOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsKindOf
}

// Implements [typeinfo.Repeats] for a slice of IsKindOf.
func (op *IsKindOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// Friendly name of the object's kind.
type KindOf struct {
	Object  rtti.TextEval
	Nothing bool
	Markup  map[string]any
}

// kind_of, a type of flow.
var Zt_KindOf typeinfo.Flow

// Implements [typeinfo.Instance]
func (*KindOf) TypeInfo() typeinfo.T {
	return &Zt_KindOf
}

// Implements [typeinfo.Markup]
func (op *KindOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*KindOf)(nil)

// Holds a slice of type KindOf.
type KindOf_Slice []KindOf

// Implements [typeinfo.Instance] for a slice of KindOf.
func (*KindOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_KindOf
}

// Implements [typeinfo.Repeats] for a slice of KindOf.
func (op *KindOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// A list of compatible kinds.
type KindsOf struct {
	Kind   string
	Markup map[string]any
}

// kinds_of, a type of flow.
var Zt_KindsOf typeinfo.Flow

// Implements [typeinfo.Instance]
func (*KindsOf) TypeInfo() typeinfo.T {
	return &Zt_KindsOf
}

// Implements [typeinfo.Markup]
func (op *KindsOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*KindsOf)(nil)

// Holds a slice of type KindsOf.
type KindsOf_Slice []KindsOf

// Implements [typeinfo.Instance] for a slice of KindsOf.
func (*KindsOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_KindsOf
}

// Implements [typeinfo.Repeats] for a slice of KindsOf.
func (op *KindsOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// List of the field names of a kind.
type FieldsOfKind struct {
	KindName rtti.TextEval
	Markup   map[string]any
}

// fields_of_kind, a type of flow.
var Zt_FieldsOfKind typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FieldsOfKind) TypeInfo() typeinfo.T {
	return &Zt_FieldsOfKind
}

// Implements [typeinfo.Markup]
func (op *FieldsOfKind) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*FieldsOfKind)(nil)

// Holds a slice of type FieldsOfKind.
type FieldsOfKind_Slice []FieldsOfKind

// Implements [typeinfo.Instance] for a slice of FieldsOfKind.
func (*FieldsOfKind_Slice) TypeInfo() typeinfo.T {
	return &Zt_FieldsOfKind
}

// Implements [typeinfo.Repeats] for a slice of FieldsOfKind.
func (op *FieldsOfKind_Slice) Repeats() bool {
	return len(*op) > 0
}

// Increases the value of a trait held by an object aspect.
// Returns the new value of the trait.
type IncrementAspect struct {
	Target rtti.TextEval
	Aspect rtti.TextEval
	Step   rtti.NumEval
	Clamp  rtti.BoolEval
	Markup map[string]any
}

// increment_aspect, a type of flow.
var Zt_IncrementAspect typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IncrementAspect) TypeInfo() typeinfo.T {
	return &Zt_IncrementAspect
}

// Implements [typeinfo.Markup]
func (op *IncrementAspect) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*IncrementAspect)(nil)
var _ rtti.Execute = (*IncrementAspect)(nil)

// Holds a slice of type IncrementAspect.
type IncrementAspect_Slice []IncrementAspect

// Implements [typeinfo.Instance] for a slice of IncrementAspect.
func (*IncrementAspect_Slice) TypeInfo() typeinfo.T {
	return &Zt_IncrementAspect
}

// Implements [typeinfo.Repeats] for a slice of IncrementAspect.
func (op *IncrementAspect_Slice) Repeats() bool {
	return len(*op) > 0
}

// Increases the value of a trait held by an object aspect.
// Returns the new value of the trait.
type DecrementAspect struct {
	Target rtti.TextEval
	Aspect rtti.TextEval
	Step   rtti.NumEval
	Clamp  rtti.BoolEval
	Markup map[string]any
}

// decrement_aspect, a type of flow.
var Zt_DecrementAspect typeinfo.Flow

// Implements [typeinfo.Instance]
func (*DecrementAspect) TypeInfo() typeinfo.T {
	return &Zt_DecrementAspect
}

// Implements [typeinfo.Markup]
func (op *DecrementAspect) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*DecrementAspect)(nil)
var _ rtti.Execute = (*DecrementAspect)(nil)

// Holds a slice of type DecrementAspect.
type DecrementAspect_Slice []DecrementAspect

// Implements [typeinfo.Instance] for a slice of DecrementAspect.
func (*DecrementAspect_Slice) TypeInfo() typeinfo.T {
	return &Zt_DecrementAspect
}

// Implements [typeinfo.Repeats] for a slice of DecrementAspect.
func (op *DecrementAspect_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_SetValue = typeinfo.Flow{
		Name: "set_value",
		Lede: "set",
		Terms: []typeinfo.Term{{
			Name: "target",
			Markup: map[string]any{
				"comment": "Object property or variable into which to write the value.",
			},
			Type: &rtti.Zt_Address,
		}, {
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"comment": "The value to copy into the destination.",
			},
			Type: &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Store a value into a variable or object.", "Values are specified as a generic [Assignment].", "The various \"From\" commands exist to cast specific value types into an assignment.", "", "WARNING: This doesn't convert values from one type to another.", "For example:", "  Set:value:", "  - \"@some_local_variable\"", "  - FromText: \"a piece of text to store.\"", "will only work if the local variable can store text. If the variable was declared as a number, the command will generate an error."},
		},
	}
	Zt_SetState = typeinfo.Flow{
		Name: "set_state",
		Lede: "set",
		Terms: []typeinfo.Term{{
			Name: "target",
			Markup: map[string]any{
				"comment": "Object or record to change.",
			},
			Type: &rtti.Zt_Address,
		}, {
			Name:  "trait",
			Label: "state",
			Markup: map[string]any{
				"comment": []interface{}{"Name of the state to set.", "Only one state in a state set is considered active at a time so this implicitly deactivates the other states in its set.", "Errors if the state wasn't declared as part of the object's kind."},
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Set the state of an object or record.", "See also: story `Define state:names:`."},
		},
	}
	Zt_ObjectDot = typeinfo.Flow{
		Name: "object_dot",
		Lede: "object",
		Terms: []typeinfo.Term{{
			Name: "name",
			Markup: map[string]any{
				"comment": "Id or friendly name of the object.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:     "dot",
			Label:    "dot",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"comment": "A field or path within the object to read from.",
			},
			Type: &Zt_Dot,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Address,
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumEval,
			&rtti.Zt_TextEval,
			&rtti.Zt_RecordEval,
			&rtti.Zt_NumListEval,
			&rtti.Zt_TextListEval,
			&rtti.Zt_RecordListEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Read a value from an object. As a special case, if there are no dot parts, this will return the id of the object.", "In .tell files, this command is often specified with a shortcut. For example:", "  \"#my_object.some_field\"", "is a shorter way to say:", "  Object:dot:", "  - \"my object\"", "  - \"some field\"", "WARNING: This doesn't convert values from one type to another. For instance, if a field was declared as text, this will error if read as a boolean."},
		},
	}
	Zt_VariableDot = typeinfo.Flow{
		Name: "variable_dot",
		Lede: "variable",
		Terms: []typeinfo.Term{{
			Name: "name",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:     "dot",
			Label:    "dot",
			Optional: true,
			Repeats:  true,
			Type:     &Zt_Dot,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Address,
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumEval,
			&rtti.Zt_TextEval,
			&rtti.Zt_RecordEval,
			&rtti.Zt_NumListEval,
			&rtti.Zt_TextListEval,
			&rtti.Zt_RecordListEval,
		},
		Markup: map[string]any{
			"blockly-color": "MATH_HUE",
			"comment":       []interface{}{"Read a value from a variable.", "In .tell files, this command is often specified with a shortcut. For example:", "  \"@some_local_variable\"", "is a shorter way to say:", "  Variable:dot: \"some local variable\"", "WARNING: This doesn't convert values from one type to another. For instance, if a field was declared as text, this will error if read as a boolean."},
		},
	}
	Zt_AtField = typeinfo.Flow{
		Name: "at_field",
		Lede: "at_field",
		Terms: []typeinfo.Term{{
			Name: "field",
			Markup: map[string]any{
				"comment": []interface{}{"The name of the field to read or write.", "The field must exist in the object or record being accessed."},
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Dot,
		},
		Markup: map[string]any{
			"comment": "Select a named field from a record, or a named property from an object.",
		},
	}
	Zt_AtIndex = typeinfo.Flow{
		Name: "at_index",
		Lede: "at_index",
		Terms: []typeinfo.Term{{
			Name: "index",
			Markup: map[string]any{
				"comment": []interface{}{"The zero-based index to read or write.", "The index must exist within the list being targeted."},
			},
			Type: &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Dot,
		},
		Markup: map[string]any{
			"comment": "Select a value from a list of values.",
		},
	}
	Zt_NameOf = typeinfo.Flow{
		Name: "name_of",
		Lede: "name_of",
		Terms: []typeinfo.Term{{
			Name: "object",
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Full name of the object.",
		},
	}
	Zt_ObjectTraits = typeinfo.Flow{
		Name: "object_traits",
		Lede: "object",
		Terms: []typeinfo.Term{{
			Name:  "object",
			Label: "traits",
			Type:  &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
		},
		Markup: map[string]any{
			"comment": "Returns all of the object's current traits as a list of text.",
		},
	}
	Zt_IsExactKindOf = typeinfo.Flow{
		Name: "is_exact_kind_of",
		Lede: "kind_of",
		Terms: []typeinfo.Term{{
			Name: "object",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "kind",
			Label: "is_exactly",
			Type:  &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": "True if the object is exactly the named kind.",
		},
	}
	Zt_IsKindOf = typeinfo.Flow{
		Name: "is_kind_of",
		Lede: "kind_of",
		Terms: []typeinfo.Term{{
			Name: "object",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "kind",
			Label: "is",
			Type:  &prim.Zt_Text,
		}, {
			Name:     "nothing",
			Label:    "nothing",
			Optional: true,
			Markup: map[string]any{
				"comment": []interface{}{"try to check the type of nothing objects?", "normally, nothing objects have no kind."},
			},
			Type: &prim.Zt_Bool,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": "True if the object is compatible with the named kind.",
		},
	}
	Zt_KindOf = typeinfo.Flow{
		Name: "kind_of",
		Lede: "kind_of",
		Terms: []typeinfo.Term{{
			Name: "object",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:     "nothing",
			Label:    "nothing",
			Optional: true,
			Markup: map[string]any{
				"comment": []interface{}{"try to check the type of nothing objects?", "normally, nothing objects have no kind."},
			},
			Type: &prim.Zt_Bool,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Friendly name of the object's kind.",
		},
	}
	Zt_KindsOf = typeinfo.Flow{
		Name: "kinds_of",
		Lede: "kinds_of",
		Terms: []typeinfo.Term{{
			Name: "kind",
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
		},
		Markup: map[string]any{
			"comment": "A list of compatible kinds.",
		},
	}
	Zt_FieldsOfKind = typeinfo.Flow{
		Name: "fields_of_kind",
		Lede: "fields",
		Terms: []typeinfo.Term{{
			Name:  "kind_name",
			Label: "of",
			Type:  &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
		},
		Markup: map[string]any{
			"comment": "List of the field names of a kind.",
		},
	}
	Zt_IncrementAspect = typeinfo.Flow{
		Name: "increment_aspect",
		Lede: "increase",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "aspect",
			Label: "aspect",
			Type:  &rtti.Zt_TextEval,
		}, {
			Name:     "step",
			Label:    "by",
			Optional: true,
			Markup: map[string]any{
				"comment": "if not specified, increments by a single step.",
			},
			Type: &rtti.Zt_NumEval,
		}, {
			Name:     "clamp",
			Label:    "clamp",
			Optional: true,
			Markup: map[string]any{
				"comment": "if not specified, wraps.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Increases the value of a trait held by an object aspect.", "Returns the new value of the trait."},
		},
	}
	Zt_DecrementAspect = typeinfo.Flow{
		Name: "decrement_aspect",
		Lede: "decrease",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "aspect",
			Label: "aspect",
			Type:  &rtti.Zt_TextEval,
		}, {
			Name:     "step",
			Label:    "by",
			Optional: true,
			Markup: map[string]any{
				"comment": "if not specified, increments by a single step.",
			},
			Type: &rtti.Zt_NumEval,
		}, {
			Name:     "clamp",
			Label:    "clamp",
			Optional: true,
			Markup: map[string]any{
				"comment": "if not specified, wraps.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Increases the value of a trait held by an object aspect.", "Returns the new value of the trait."},
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "object",
	Comment: []string{
		"Common operations on objects, variables, and kinds.",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Dot,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_SetValue,
	&Zt_SetState,
	&Zt_ObjectDot,
	&Zt_VariableDot,
	&Zt_AtField,
	&Zt_AtIndex,
	&Zt_NameOf,
	&Zt_ObjectTraits,
	&Zt_IsExactKindOf,
	&Zt_IsKindOf,
	&Zt_KindOf,
	&Zt_KindsOf,
	&Zt_FieldsOfKind,
	&Zt_IncrementAspect,
	&Zt_DecrementAspect,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	1683104564853176068:  (*AtField)(nil),         /* dot=AtField: */
	17908840355303216180: (*AtIndex)(nil),         /* dot=AtIndex: */
	13259725831972112539: (*DecrementAspect)(nil), /* execute=Decrease:aspect: */
	9604047801594713852:  (*DecrementAspect)(nil), /* text_eval=Decrease:aspect: */
	11515881376122775668: (*DecrementAspect)(nil), /* execute=Decrease:aspect:by: */
	1589765377795283065:  (*DecrementAspect)(nil), /* text_eval=Decrease:aspect:by: */
	10691394634979399555: (*DecrementAspect)(nil), /* execute=Decrease:aspect:by:clamp: */
	16351892255943407142: (*DecrementAspect)(nil), /* text_eval=Decrease:aspect:by:clamp: */
	16567257087826189312: (*DecrementAspect)(nil), /* execute=Decrease:aspect:clamp: */
	7498537354592687963:  (*DecrementAspect)(nil), /* text_eval=Decrease:aspect:clamp: */
	2224842870997259213:  (*FieldsOfKind)(nil),    /* text_list_eval=Fields of: */
	11043224857467493683: (*IncrementAspect)(nil), /* execute=Increase:aspect: */
	1296309673842091672:  (*IncrementAspect)(nil), /* text_eval=Increase:aspect: */
	4473637830475551932:  (*IncrementAspect)(nil), /* execute=Increase:aspect:by: */
	18328024260427443133: (*IncrementAspect)(nil), /* text_eval=Increase:aspect:by: */
	1150598923989934235:  (*IncrementAspect)(nil), /* execute=Increase:aspect:by:clamp: */
	16465259325356451354: (*IncrementAspect)(nil), /* text_eval=Increase:aspect:by:clamp: */
	4522630356185077352:  (*IncrementAspect)(nil), /* execute=Increase:aspect:clamp: */
	705264554644415287:   (*IncrementAspect)(nil), /* text_eval=Increase:aspect:clamp: */
	16305715626122315047: (*KindOf)(nil),          /* text_eval=KindOf: */
	16744881049704292640: (*IsKindOf)(nil),        /* bool_eval=KindOf:is: */
	210805642732508805:   (*IsKindOf)(nil),        /* bool_eval=KindOf:is:nothing: */
	7296079450764183372:  (*IsExactKindOf)(nil),   /* bool_eval=KindOf:isExactly: */
	4254622167054960918:  (*KindOf)(nil),          /* text_eval=KindOf:nothing: */
	6869420318733086481:  (*KindsOf)(nil),         /* text_list_eval=KindsOf: */
	15519818243985955688: (*NameOf)(nil),          /* text_eval=NameOf: */
	15933580486837544843: (*ObjectTraits)(nil),    /* text_list_eval=Object traits: */
	8656684385605626625:  (*ObjectDot)(nil),       /* address=Object: */
	6106842879255343810:  (*ObjectDot)(nil),       /* bool_eval=Object: */
	14709650427635515944: (*ObjectDot)(nil),       /* num_eval=Object: */
	3322847371150895433:  (*ObjectDot)(nil),       /* num_list_eval=Object: */
	1988642049281593865:  (*ObjectDot)(nil),       /* record_eval=Object: */
	9599721143262547914:  (*ObjectDot)(nil),       /* record_list_eval=Object: */
	16083123907778192555: (*ObjectDot)(nil),       /* text_eval=Object: */
	15780956574897965792: (*ObjectDot)(nil),       /* text_list_eval=Object: */
	8121157847033684962:  (*ObjectDot)(nil),       /* address=Object:dot: */
	5205171710741514089:  (*ObjectDot)(nil),       /* bool_eval=Object:dot: */
	13854256590934503743: (*ObjectDot)(nil),       /* num_eval=Object:dot: */
	3914994200631113354:  (*ObjectDot)(nil),       /* num_list_eval=Object:dot: */
	1364775634664390090:  (*ObjectDot)(nil),       /* record_eval=Object:dot: */
	16877508779303594737: (*ObjectDot)(nil),       /* record_list_eval=Object:dot: */
	17663678026468030644: (*ObjectDot)(nil),       /* text_eval=Object:dot: */
	725008522959645559:   (*ObjectDot)(nil),       /* text_list_eval=Object:dot: */
	9616350989753725148:  (*SetState)(nil),        /* execute=Set:state: */
	3912570011939708664:  (*SetValue)(nil),        /* execute=Set:value: */
	13692207992970428220: (*VariableDot)(nil),     /* address=Variable: */
	17908519799628660539: (*VariableDot)(nil),     /* bool_eval=Variable: */
	17658028528032582325: (*VariableDot)(nil),     /* num_eval=Variable: */
	11022385456290008164: (*VariableDot)(nil),     /* num_list_eval=Variable: */
	15906653930217516836: (*VariableDot)(nil),     /* record_eval=Variable: */
	16032903663975260899: (*VariableDot)(nil),     /* record_list_eval=Variable: */
	11181798416019134386: (*VariableDot)(nil),     /* text_eval=Variable: */
	14769776891888769773: (*VariableDot)(nil),     /* text_list_eval=Variable: */
	15966558056732701531: (*VariableDot)(nil),     /* address=Variable:dot: */
	7739360284898038596:  (*VariableDot)(nil),     /* bool_eval=Variable:dot: */
	4938834444414070846:  (*VariableDot)(nil),     /* num_eval=Variable:dot: */
	14012826006150347811: (*VariableDot)(nil),     /* num_list_eval=Variable:dot: */
	3479001804857346403:  (*VariableDot)(nil),     /* record_eval=Variable:dot: */
	11938488787528882828: (*VariableDot)(nil),     /* record_list_eval=Variable:dot: */
	4798713833623285465:  (*VariableDot)(nil),     /* text_eval=Variable:dot: */
	12039638244497140214: (*VariableDot)(nil),     /* text_list_eval=Variable:dot: */
}
