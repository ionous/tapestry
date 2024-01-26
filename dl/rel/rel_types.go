// Code generated by Tapestry; edit at your own risk.
package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type ReciprocalOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ReciprocalOf) Inspect() (typeinfo.T, bool) {
	return &Zt_ReciprocalOf, false
}

// return a valid markup map, creating it if necessary.
func (op *ReciprocalOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*ReciprocalOf)(nil)

// reciprocal_of, a type of flow.
var Zt_ReciprocalOf = typeinfo.Flow{
	Name: "reciprocal_of",
	Lede: "reciprocal",
	Terms: []typeinfo.Term{{
		Name: "via",
		Type: &Zt_RelationName,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_TextEval,
	},
	Markup: map[string]any{
		"comment": "Returns the implied relative of a noun (ex. the source in a one-to-many relation.).",
	},
}

// holds a slice of type reciprocal_of
// FIX: duplicates the spec decl.
type ReciprocalOf_Slice []ReciprocalOf

// implements typeinfo.Inspector
func (*ReciprocalOf_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_ReciprocalOf, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type ReciprocalsOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ReciprocalsOf) Inspect() (typeinfo.T, bool) {
	return &Zt_ReciprocalsOf, false
}

// return a valid markup map, creating it if necessary.
func (op *ReciprocalsOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*ReciprocalsOf)(nil)

// reciprocals_of, a type of flow.
var Zt_ReciprocalsOf = typeinfo.Flow{
	Name: "reciprocals_of",
	Lede: "reciprocals",
	Terms: []typeinfo.Term{{
		Name: "via",
		Type: &Zt_RelationName,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_TextListEval,
	},
	Markup: map[string]any{
		"comment": "Returns the implied relative of a noun (ex. the sources of a many-to-many relation.).",
	},
}

// holds a slice of type reciprocals_of
// FIX: duplicates the spec decl.
type ReciprocalsOf_Slice []ReciprocalsOf

// implements typeinfo.Inspector
func (*ReciprocalsOf_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_ReciprocalsOf, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type Relate struct {
	Object   rtti.TextEval
	ToObject rtti.TextEval
	Via      string
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*Relate) Inspect() (typeinfo.T, bool) {
	return &Zt_Relate, false
}

// return a valid markup map, creating it if necessary.
func (op *Relate) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*Relate)(nil)

// relate, a type of flow.
var Zt_Relate = typeinfo.Flow{
	Name: "relate",
	Lede: "relate",
	Terms: []typeinfo.Term{{
		Name: "object",
		Type: &rtti.Zt_TextEval,
	}, {
		Name:  "to_object",
		Label: "to",
		Type:  &rtti.Zt_TextEval,
	}, {
		Name:  "via",
		Label: "via",
		Type:  &Zt_RelationName,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
	Markup: map[string]any{
		"comment": "Relate two nouns.",
	},
}

// holds a slice of type relate
// FIX: duplicates the spec decl.
type Relate_Slice []Relate

// implements typeinfo.Inspector
func (*Relate_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_Relate, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type RelativeOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RelativeOf) Inspect() (typeinfo.T, bool) {
	return &Zt_RelativeOf, false
}

// return a valid markup map, creating it if necessary.
func (op *RelativeOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*RelativeOf)(nil)

// relative_of, a type of flow.
var Zt_RelativeOf = typeinfo.Flow{
	Name: "relative_of",
	Lede: "relative",
	Terms: []typeinfo.Term{{
		Name: "via",
		Type: &Zt_RelationName,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_TextEval,
	},
	Markup: map[string]any{
		"comment": "Returns the relative of a noun (ex. the target of a one-to-one relation.).",
	},
}

// holds a slice of type relative_of
// FIX: duplicates the spec decl.
type RelativeOf_Slice []RelativeOf

// implements typeinfo.Inspector
func (*RelativeOf_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_RelativeOf, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type RelativesOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RelativesOf) Inspect() (typeinfo.T, bool) {
	return &Zt_RelativesOf, false
}

// return a valid markup map, creating it if necessary.
func (op *RelativesOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*RelativesOf)(nil)

// relatives_of, a type of flow.
var Zt_RelativesOf = typeinfo.Flow{
	Name: "relatives_of",
	Lede: "relatives",
	Terms: []typeinfo.Term{{
		Name: "via",
		Type: &Zt_RelationName,
	}, {
		Name:  "object",
		Label: "object",
		Type:  &rtti.Zt_TextEval,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_TextListEval,
	},
	Markup: map[string]any{
		"comment": "Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).",
	},
}

// holds a slice of type relatives_of
// FIX: duplicates the spec decl.
type RelativesOf_Slice []RelativesOf

// implements typeinfo.Inspector
func (*RelativesOf_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_RelativesOf, true
}

var Zt_RelationName = typeinfo.Str{
	Name: "relation_name",
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name:       "rel",
	Flow:       z_flow_list,
	Str:        z_str_list,
	Signatures: z_signatures,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_ReciprocalOf,
	&Zt_ReciprocalsOf,
	&Zt_Relate,
	&Zt_RelativeOf,
	&Zt_RelativesOf,
}

// a list of all strs in this this package
var z_str_list = []*typeinfo.Str{
	&Zt_RelationName,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]any{
	6987621383789599381:  (*ReciprocalOf)(nil),  /* text_eval=Reciprocal:object: */
	16170704865359856399: (*ReciprocalsOf)(nil), /* text_list_eval=Reciprocals:object: */
	15160920709871392391: (*Relate)(nil),        /* execute=Relate:to:via: */
	14535552277213572673: (*RelativeOf)(nil),    /* text_eval=Relative:object: */
	13180339401044333799: (*RelativesOf)(nil),   /* text_list_eval=Relatives:object: */
}
