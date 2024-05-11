// Code generated by Tapestry; edit at your own risk.
package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// Returns the implied relative of a noun (ex. the source in a one-to-many relation.).
type ReciprocalOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// reciprocal_of, a type of flow.
var Zt_ReciprocalOf typeinfo.Flow

// implements typeinfo.Instance
func (*ReciprocalOf) TypeInfo() typeinfo.T {
	return &Zt_ReciprocalOf
}

// implements typeinfo.Markup
func (op *ReciprocalOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*ReciprocalOf)(nil)

// holds a slice of type reciprocal_of
type ReciprocalOf_Slice []ReciprocalOf

// implements typeinfo.Instance
func (*ReciprocalOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_ReciprocalOf
}

// implements typeinfo.Repeats
func (op *ReciprocalOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns the implied relative of a noun (ex. the sources of a many-to-many relation.).
type ReciprocalsOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// reciprocals_of, a type of flow.
var Zt_ReciprocalsOf typeinfo.Flow

// implements typeinfo.Instance
func (*ReciprocalsOf) TypeInfo() typeinfo.T {
	return &Zt_ReciprocalsOf
}

// implements typeinfo.Markup
func (op *ReciprocalsOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*ReciprocalsOf)(nil)

// holds a slice of type reciprocals_of
type ReciprocalsOf_Slice []ReciprocalsOf

// implements typeinfo.Instance
func (*ReciprocalsOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_ReciprocalsOf
}

// implements typeinfo.Repeats
func (op *ReciprocalsOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// Relate two nouns.
type Relate struct {
	Object   rtti.TextEval
	ToObject rtti.TextEval
	Via      string
	Markup   map[string]any
}

// relate, a type of flow.
var Zt_Relate typeinfo.Flow

// implements typeinfo.Instance
func (*Relate) TypeInfo() typeinfo.T {
	return &Zt_Relate
}

// implements typeinfo.Markup
func (op *Relate) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*Relate)(nil)

// holds a slice of type relate
type Relate_Slice []Relate

// implements typeinfo.Instance
func (*Relate_Slice) TypeInfo() typeinfo.T {
	return &Zt_Relate
}

// implements typeinfo.Repeats
func (op *Relate_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns the relative of a noun (ex. the target of a one-to-one relation.).
type RelativeOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// relative_of, a type of flow.
var Zt_RelativeOf typeinfo.Flow

// implements typeinfo.Instance
func (*RelativeOf) TypeInfo() typeinfo.T {
	return &Zt_RelativeOf
}

// implements typeinfo.Markup
func (op *RelativeOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextEval = (*RelativeOf)(nil)

// holds a slice of type relative_of
type RelativeOf_Slice []RelativeOf

// implements typeinfo.Instance
func (*RelativeOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_RelativeOf
}

// implements typeinfo.Repeats
func (op *RelativeOf_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).
type RelativesOf struct {
	Via    string
	Object rtti.TextEval
	Markup map[string]any
}

// relatives_of, a type of flow.
var Zt_RelativesOf typeinfo.Flow

// implements typeinfo.Instance
func (*RelativesOf) TypeInfo() typeinfo.T {
	return &Zt_RelativesOf
}

// implements typeinfo.Markup
func (op *RelativesOf) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*RelativesOf)(nil)

// holds a slice of type relatives_of
type RelativesOf_Slice []RelativesOf

// implements typeinfo.Instance
func (*RelativesOf_Slice) TypeInfo() typeinfo.T {
	return &Zt_RelativesOf
}

// implements typeinfo.Repeats
func (op *RelativesOf_Slice) Repeats() bool {
	return len(*op) > 0
}

var Zt_RelationName = typeinfo.Str{
	Name: "relation_name",
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_ReciprocalOf = typeinfo.Flow{
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
	Zt_ReciprocalsOf = typeinfo.Flow{
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
	Zt_Relate = typeinfo.Flow{
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
	Zt_RelativeOf = typeinfo.Flow{
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
	Zt_RelativesOf = typeinfo.Flow{
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
var z_signatures = map[uint64]typeinfo.Instance{
	6987621383789599381:  (*ReciprocalOf)(nil),  /* text_eval=Reciprocal:object: */
	16170704865359856399: (*ReciprocalsOf)(nil), /* text_list_eval=Reciprocals:object: */
	15160920709871392391: (*Relate)(nil),        /* execute=Relate:to:via: */
	14535552277213572673: (*RelativeOf)(nil),    /* text_eval=Relative:object: */
	13180339401044333799: (*RelativesOf)(nil),   /* text_list_eval=Relatives:object: */
}
