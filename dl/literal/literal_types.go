// literal
package literal

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// literal_value, a type of slot.
var Zt_LiteralValue = typeinfo.Slot{
	Name: "literal_value",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Slot for constant values.",
	},
}

// Holds a single slot.
type LiteralValue_Slot struct{ Value LiteralValue }

// Implements [typeinfo.Instance] for a single slot.
func (*LiteralValue_Slot) TypeInfo() typeinfo.T {
	return &Zt_LiteralValue
}

// Holds a slice of slots.
type LiteralValue_Slots []LiteralValue

// Implements [typeinfo.Instance] for a slice of slots.
func (*LiteralValue_Slots) TypeInfo() typeinfo.T {
	return &Zt_LiteralValue
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *LiteralValue_Slots) Repeats() bool {
	return len(*op) > 0
}

// Specify an explicit true or false.
type BoolValue struct {
	Value  bool
	Kind   string
	Markup map[string]any
}

// bool_value, a type of flow.
var Zt_BoolValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*BoolValue) TypeInfo() typeinfo.T {
	return &Zt_BoolValue
}

// Implements [typeinfo.Markup]
func (op *BoolValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*BoolValue)(nil)
var _ LiteralValue = (*BoolValue)(nil)

// Holds a slice of type BoolValue.
type BoolValue_Slice []BoolValue

// Implements [typeinfo.Instance] for a slice of BoolValue.
func (*BoolValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_BoolValue
}

// Implements [typeinfo.Repeats] for a slice of BoolValue.
func (op *BoolValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// A fixed value of a record.
type FieldValue struct {
	Field  string
	Value  LiteralValue
	Markup map[string]any
}

// field_value, a type of flow.
var Zt_FieldValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FieldValue) TypeInfo() typeinfo.T {
	return &Zt_FieldValue
}

// Implements [typeinfo.Markup]
func (op *FieldValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Holds a slice of type FieldValue.
type FieldValue_Slice []FieldValue

// Implements [typeinfo.Instance] for a slice of FieldValue.
func (*FieldValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_FieldValue
}

// Implements [typeinfo.Repeats] for a slice of FieldValue.
func (op *FieldValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// A series of values all for the same record.
// While it can be specified wherever a literal value can, it only has meaning when the record type is known.
type FieldList struct {
	Fields []FieldValue
	Markup map[string]any
}

// field_list, a type of flow.
var Zt_FieldList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FieldList) TypeInfo() typeinfo.T {
	return &Zt_FieldList
}

// Implements [typeinfo.Markup]
func (op *FieldList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ LiteralValue = (*FieldList)(nil)

// Holds a slice of type FieldList.
type FieldList_Slice []FieldList

// Implements [typeinfo.Instance] for a slice of FieldList.
func (*FieldList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FieldList
}

// Implements [typeinfo.Repeats] for a slice of FieldList.
func (op *FieldList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify a particular number.
type NumValue struct {
	Value  float64
	Kind   string
	Markup map[string]any
}

// num_value, a type of flow.
var Zt_NumValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NumValue) TypeInfo() typeinfo.T {
	return &Zt_NumValue
}

// Implements [typeinfo.Markup]
func (op *NumValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumEval = (*NumValue)(nil)
var _ LiteralValue = (*NumValue)(nil)

// Holds a slice of type NumValue.
type NumValue_Slice []NumValue

// Implements [typeinfo.Instance] for a slice of NumValue.
func (*NumValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_NumValue
}

// Implements [typeinfo.Repeats] for a slice of NumValue.
func (op *NumValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Number List: Specify a list of numbers.
type NumValues struct {
	Values []float64
	Kind   string
	Markup map[string]any
}

// num_values, a type of flow.
var Zt_NumValues typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NumValues) TypeInfo() typeinfo.T {
	return &Zt_NumValues
}

// Implements [typeinfo.Markup]
func (op *NumValues) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumListEval = (*NumValues)(nil)
var _ LiteralValue = (*NumValues)(nil)

// Holds a slice of type NumValues.
type NumValues_Slice []NumValues

// Implements [typeinfo.Instance] for a slice of NumValues.
func (*NumValues_Slice) TypeInfo() typeinfo.T {
	return &Zt_NumValues
}

// Implements [typeinfo.Repeats] for a slice of NumValues.
func (op *NumValues_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify a record composed of literal values.
type RecordValue struct {
	Kind   string
	Fields []FieldValue
	Cache  RecordCache
	Markup map[string]any
}

// record_value, a type of flow.
var Zt_RecordValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*RecordValue) TypeInfo() typeinfo.T {
	return &Zt_RecordValue
}

// Implements [typeinfo.Markup]
func (op *RecordValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.RecordEval = (*RecordValue)(nil)
var _ LiteralValue = (*RecordValue)(nil)

// Holds a slice of type RecordValue.
type RecordValue_Slice []RecordValue

// Implements [typeinfo.Instance] for a slice of RecordValue.
func (*RecordValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_RecordValue
}

// Implements [typeinfo.Repeats] for a slice of RecordValue.
func (op *RecordValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify a series of records, all of the same kind.
type RecordList struct {
	Kind    string
	Records []FieldList
	Cache   RecordsCache
	Markup  map[string]any
}

// record_list, a type of flow.
var Zt_RecordList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*RecordList) TypeInfo() typeinfo.T {
	return &Zt_RecordList
}

// Implements [typeinfo.Markup]
func (op *RecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.RecordListEval = (*RecordList)(nil)
var _ LiteralValue = (*RecordList)(nil)

// Holds a slice of type RecordList.
type RecordList_Slice []RecordList

// Implements [typeinfo.Instance] for a slice of RecordList.
func (*RecordList_Slice) TypeInfo() typeinfo.T {
	return &Zt_RecordList
}

// Implements [typeinfo.Repeats] for a slice of RecordList.
func (op *RecordList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify a small bit of text.
type TextValue struct {
	Value  string
	Kind   string
	Markup map[string]any
}

// text_value, a type of flow.
var Zt_TextValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextValue) TypeInfo() typeinfo.T {
	return &Zt_TextValue
}

// Implements [typeinfo.Markup]
func (op *TextValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*TextValue)(nil)
var _ LiteralValue = (*TextValue)(nil)

// Holds a slice of type TextValue.
type TextValue_Slice []TextValue

// Implements [typeinfo.Instance] for a slice of TextValue.
func (*TextValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextValue
}

// Implements [typeinfo.Repeats] for a slice of TextValue.
func (op *TextValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Text List: Specifies a set of text values.
type TextValues struct {
	Values []string
	Kind   string
	Markup map[string]any
}

// text_values, a type of flow.
var Zt_TextValues typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextValues) TypeInfo() typeinfo.T {
	return &Zt_TextValues
}

// Implements [typeinfo.Markup]
func (op *TextValues) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*TextValues)(nil)
var _ LiteralValue = (*TextValues)(nil)

// Holds a slice of type TextValues.
type TextValues_Slice []TextValues

// Implements [typeinfo.Instance] for a slice of TextValues.
func (*TextValues_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextValues
}

// Implements [typeinfo.Repeats] for a slice of TextValues.
func (op *TextValues_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_BoolValue = typeinfo.Flow{
		Name: "bool_value",
		Lede: "bool",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Type:  &prim.Zt_Bool,
		}, {
			Name:     "kind",
			Label:    "kind",
			Optional: true,
			Type:     &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify an explicit true or false.",
		},
	}
	Zt_FieldValue = typeinfo.Flow{
		Name: "field_value",
		Lede: "field",
		Terms: []typeinfo.Term{{
			Name: "field",
			Type: &prim.Zt_Text,
		}, {
			Name:  "value",
			Label: "value",
			Type:  &Zt_LiteralValue,
		}},
		Markup: map[string]any{
			"comment": "A fixed value of a record.",
		},
	}
	Zt_FieldList = typeinfo.Flow{
		Name: "field_list",
		Lede: "field_list",
		Terms: []typeinfo.Term{{
			Name:    "fields",
			Repeats: true,
			Type:    &Zt_FieldValue,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": []interface{}{"A series of values all for the same record.", "While it can be specified wherever a literal value can, it only has meaning when the record type is known."},
		},
	}
	Zt_NumValue = typeinfo.Flow{
		Name: "num_value",
		Lede: "num",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Type:  &prim.Zt_Num,
		}, {
			Name:     "kind",
			Label:    "kind",
			Optional: true,
			Type:     &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a particular number.",
		},
	}
	Zt_NumValues = typeinfo.Flow{
		Name: "num_values",
		Lede: "num",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "values",
			Repeats: true,
			Type:    &prim.Zt_Num,
		}, {
			Name:     "kind",
			Label:    "kind",
			Optional: true,
			Type:     &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumListEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Number List: Specify a list of numbers.",
		},
	}
	Zt_RecordValue = typeinfo.Flow{
		Name: "record_value",
		Lede: "record",
		Terms: []typeinfo.Term{{
			Name: "kind",
			Type: &prim.Zt_Text,
		}, {
			Name:    "fields",
			Label:   "fields",
			Repeats: true,
			Type:    &Zt_FieldValue,
		}, {
			Name:    "cache",
			Label:   "cache",
			Private: true,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_RecordEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a record composed of literal values.",
		},
	}
	Zt_RecordList = typeinfo.Flow{
		Name: "record_list",
		Lede: "record",
		Terms: []typeinfo.Term{{
			Name: "kind",
			Type: &prim.Zt_Text,
		}, {
			Name:    "records",
			Label:   "values",
			Repeats: true,
			Type:    &Zt_FieldList,
		}, {
			Name:    "cache",
			Label:   "cache",
			Private: true,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_RecordListEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a series of records, all of the same kind.",
		},
	}
	Zt_TextValue = typeinfo.Flow{
		Name: "text_value",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Type:  &prim.Zt_Text,
		}, {
			Name:     "kind",
			Label:    "kind",
			Optional: true,
			Type:     &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a small bit of text.",
		},
	}
	Zt_TextValues = typeinfo.Flow{
		Name: "text_values",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "values",
			Repeats: true,
			Type:    &prim.Zt_Text,
		}, {
			Name:     "kind",
			Label:    "kind",
			Optional: true,
			Type:     &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Text List: Specifies a set of text values.",
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "literal",
	Comment: []string{
		"literal",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_LiteralValue,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_BoolValue,
	&Zt_FieldValue,
	&Zt_FieldList,
	&Zt_NumValue,
	&Zt_NumValues,
	&Zt_RecordValue,
	&Zt_RecordList,
	&Zt_TextValue,
	&Zt_TextValues,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	17656638186047966738: (*FieldValue)(nil),  /* Field:value: */
	2028829358589965004:  (*BoolValue)(nil),   /* bool_eval=Bool value: */
	11511029631426206694: (*BoolValue)(nil),   /* literal_value=Bool value: */
	10808478223495627740: (*BoolValue)(nil),   /* bool_eval=Bool value:kind: */
	3205100557739257174:  (*BoolValue)(nil),   /* literal_value=Bool value:kind: */
	3071550758741756995:  (*FieldList)(nil),   /* literal_value=FieldList: */
	15362209855253663632: (*NumValue)(nil),    /* literal_value=Num value: */
	16565175635984030252: (*NumValue)(nil),    /* num_eval=Num value: */
	607468628506983640:   (*NumValue)(nil),    /* literal_value=Num value:kind: */
	11555965194136863548: (*NumValue)(nil),    /* num_eval=Num value:kind: */
	12282038377752822419: (*NumValues)(nil),   /* literal_value=Num values: */
	8089072108541894314:  (*NumValues)(nil),   /* num_list_eval=Num values: */
	16844579494806292121: (*NumValues)(nil),   /* literal_value=Num values:kind: */
	18166562587031464546: (*NumValues)(nil),   /* num_list_eval=Num values:kind: */
	5942123174065535899:  (*RecordValue)(nil), /* literal_value=Record:fields: */
	5794725022419893180:  (*RecordValue)(nil), /* record_eval=Record:fields: */
	8711768526197034738:  (*RecordList)(nil),  /* literal_value=Record:values: */
	14652198550804167624: (*RecordList)(nil),  /* record_list_eval=Record:values: */
	13114183353368545439: (*TextValue)(nil),   /* literal_value=Text value: */
	4705033170011872932:  (*TextValue)(nil),   /* text_eval=Text value: */
	6339203747835692413:  (*TextValue)(nil),   /* literal_value=Text value:kind: */
	18213962910681037476: (*TextValue)(nil),   /* text_eval=Text value:kind: */
	2231933745037898906:  (*TextValues)(nil),  /* literal_value=Text values: */
	5151885117815687006:  (*TextValues)(nil),  /* text_list_eval=Text values: */
	4866602258857929938:  (*TextValues)(nil),  /* literal_value=Text values:kind: */
	10847058762172166526: (*TextValues)(nil),  /* text_list_eval=Text values:kind: */
}
