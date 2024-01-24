// Code generated by Tapestry; edit at your own risk.
package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ReciprocalOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ReciprocalOf) Inspect() typeinfo.T {
	return &Z_ReciprocalOf_T
}

// return a valid markup map, creating it if necessary.
func (op *ReciprocalOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// reciprocal_of, a type of flow.
const Z_ReciprocalOf_Name = "reciprocal_of"

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*ReciprocalOf)(nil)

var Z_ReciprocalOf_T = typeinfo.Flow{
	Name: Z_ReciprocalOf_Name,
	Lede: "reciprocal",
	Terms: []typeinfo.Term{{
		Name:  "via",
		Label: "_",
		Type:  &Z_RelationName_T,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_TextEval_T,
	},
	Markup: map[string]any{
		"comment": "Returns the implied relative of a noun (ex. the source in a one-to-many relation.).",
	},
}

// holds a slice of type reciprocal_of
// FIX: duplicates the spec decl.
type FIX_ReciprocalOf_Slice []ReciprocalOf

// implements typeinfo.Inspector
func (*ReciprocalOf_Slice) Inspect() typeinfo.T {
	return &Z_ReciprocalOf_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ReciprocalsOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ReciprocalsOf) Inspect() typeinfo.T {
	return &Z_ReciprocalsOf_T
}

// return a valid markup map, creating it if necessary.
func (op *ReciprocalsOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// reciprocals_of, a type of flow.
const Z_ReciprocalsOf_Name = "reciprocals_of"

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*ReciprocalsOf)(nil)

var Z_ReciprocalsOf_T = typeinfo.Flow{
	Name: Z_ReciprocalsOf_Name,
	Lede: "reciprocals",
	Terms: []typeinfo.Term{{
		Name:  "via",
		Label: "_",
		Type:  &Z_RelationName_T,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_TextListEval_T,
	},
	Markup: map[string]any{
		"comment": "Returns the implied relative of a noun (ex. the sources of a many-to-many relation.).",
	},
}

// holds a slice of type reciprocals_of
// FIX: duplicates the spec decl.
type FIX_ReciprocalsOf_Slice []ReciprocalsOf

// implements typeinfo.Inspector
func (*ReciprocalsOf_Slice) Inspect() typeinfo.T {
	return &Z_ReciprocalsOf_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Relate struct {
	Object   rtti.TextEval
	ToObject rtti.TextEval
	Via      string
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*Relate) Inspect() typeinfo.T {
	return &Z_Relate_T
}

// return a valid markup map, creating it if necessary.
func (op *Relate) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// relate, a type of flow.
const Z_Relate_Name = "relate"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*Relate)(nil)

var Z_Relate_T = typeinfo.Flow{
	Name: Z_Relate_Name,
	Lede: "relate",
	Terms: []typeinfo.Term{{
		Name:  "object",
		Label: "_",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:  "to_object",
		Label: "to",
		Type:  &rtti.Z_TextEval_T,
	}, {
		Name:  "via",
		Label: "via",
		Type:  &Z_RelationName_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
	Markup: map[string]any{
		"comment": "Relate two nouns.",
	},
}

// holds a slice of type relate
// FIX: duplicates the spec decl.
type FIX_Relate_Slice []Relate

// implements typeinfo.Inspector
func (*Relate_Slice) Inspect() typeinfo.T {
	return &Z_Relate_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_RelativeOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RelativeOf) Inspect() typeinfo.T {
	return &Z_RelativeOf_T
}

// return a valid markup map, creating it if necessary.
func (op *RelativeOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// relative_of, a type of flow.
const Z_RelativeOf_Name = "relative_of"

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*RelativeOf)(nil)

var Z_RelativeOf_T = typeinfo.Flow{
	Name: Z_RelativeOf_Name,
	Lede: "relative",
	Terms: []typeinfo.Term{{
		Name:  "via",
		Label: "_",
		Type:  &Z_RelationName_T,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_TextEval_T,
	},
	Markup: map[string]any{
		"comment": "Returns the relative of a noun (ex. the target of a one-to-one relation.).",
	},
}

// holds a slice of type relative_of
// FIX: duplicates the spec decl.
type FIX_RelativeOf_Slice []RelativeOf

// implements typeinfo.Inspector
func (*RelativeOf_Slice) Inspect() typeinfo.T {
	return &Z_RelativeOf_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_RelativesOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RelativesOf) Inspect() typeinfo.T {
	return &Z_RelativesOf_T
}

// return a valid markup map, creating it if necessary.
func (op *RelativesOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// relatives_of, a type of flow.
const Z_RelativesOf_Name = "relatives_of"

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*RelativesOf)(nil)

var Z_RelativesOf_T = typeinfo.Flow{
	Name: Z_RelativesOf_Name,
	Lede: "relatives",
	Terms: []typeinfo.Term{{
		Name:  "via",
		Label: "_",
		Type:  &Z_RelationName_T,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Z_TextEval_T,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_TextListEval_T,
	},
	Markup: map[string]any{
		"comment": "Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).",
	},
}

// holds a slice of type relatives_of
// FIX: duplicates the spec decl.
type FIX_RelativesOf_Slice []RelativesOf

// implements typeinfo.Inspector
func (*RelativesOf_Slice) Inspect() typeinfo.T {
	return &Z_RelativesOf_T
}

// relation_name, a type of str.
const Z_RelationName_Name = "relation_name"

var Z_RelationName_T = typeinfo.Str{
	Name: Z_RelationName_Name,
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "rel",
	Flow: z_flow_list,
	Str:  z_str_list,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_ReciprocalOf_T,
	&Z_ReciprocalsOf_T,
	&Z_Relate_T,
	&Z_RelativeOf_T,
	&Z_RelativesOf_T,
}

// a list of all strs in this this package
var z_str_list = []*typeinfo.Str{
	&Z_RelationName_T,
}