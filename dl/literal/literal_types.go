// Supply constant values to runtime evaluations.
// ( ie. a specific number when a [NumEval] is required. )
//
// Tell files support "shortcuts" which turn primitive values into literal commands. For instance, the number '5' in a .tell file is automatically transformed into the command "Num value: 5" whenever that's needed.
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
		"comment":       "A slot to identify constant values.",
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
	Markup map[string]any `json:",omitempty"`
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

// Specify a particular number.
type NumValue struct {
	Value  float64
	Markup map[string]any `json:",omitempty"`
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

// Specify a list of literal numbers.
type NumList struct {
	Values []float64
	Markup map[string]any `json:",omitempty"`
}

// num_list, a type of flow.
var Zt_NumList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NumList) TypeInfo() typeinfo.T {
	return &Zt_NumList
}

// Implements [typeinfo.Markup]
func (op *NumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumListEval = (*NumList)(nil)
var _ LiteralValue = (*NumList)(nil)

// Holds a slice of type NumList.
type NumList_Slice []NumList

// Implements [typeinfo.Instance] for a slice of NumList.
func (*NumList_Slice) TypeInfo() typeinfo.T {
	return &Zt_NumList
}

// Implements [typeinfo.Repeats] for a slice of NumList.
func (op *NumList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify some constant text.
type TextValue struct {
	KindName string
	Value    string
	Markup   map[string]any `json:",omitempty"`
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

// Specify a list of literal text values.
type TextList struct {
	KindName string
	Values   []string
	Markup   map[string]any `json:",omitempty"`
}

// text_list, a type of flow.
var Zt_TextList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextList) TypeInfo() typeinfo.T {
	return &Zt_TextList
}

// Implements [typeinfo.Markup]
func (op *TextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*TextList)(nil)
var _ LiteralValue = (*TextList)(nil)

// Holds a slice of type TextList.
type TextList_Slice []TextList

// Implements [typeinfo.Instance] for a slice of TextList.
func (*TextList_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextList
}

// Implements [typeinfo.Repeats] for a slice of TextList.
func (op *TextList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Specify a record composed of literal values.
type RecordValue struct {
	KindName string
	Fields   []FieldValue
	cache    RecordCache
	Markup   map[string]any `json:",omitempty"`
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
	KindName string
	Records  []FieldList
	cache    RecordsCache
	Markup   map[string]any `json:",omitempty"`
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

// A series of values used to build a record.
type FieldList struct {
	Fields []FieldValue
	Markup map[string]any `json:",omitempty"`
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

// The name and value of a field used for initializing a literal record.
type FieldValue struct {
	FieldName string
	Value     LiteralValue
	Markup    map[string]any `json:",omitempty"`
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

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_BoolValue = typeinfo.Flow{
		Name: "bool_value",
		Lede: "bool",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"comment": "The true or false value.",
			},
			Type: &prim.Zt_Bool,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify an explicit true or false.",
		},
	}
	Zt_NumValue = typeinfo.Flow{
		Name: "num_value",
		Lede: "num",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"comment": "A literal number.",
			},
			Type: &prim.Zt_Num,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a particular number.",
		},
	}
	Zt_NumList = typeinfo.Flow{
		Name: "num_list",
		Lede: "num",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "values",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Zero or more literal numbers.",
			},
			Type: &prim.Zt_Num,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumListEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a list of literal numbers.",
		},
	}
	Zt_TextValue = typeinfo.Flow{
		Name: "text_value",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:     "kind_name",
			Label:    "kind",
			Optional: true,
			Markup: map[string]any{
				"comment": []string{"Optionally, when the text represents the name of an (existing) object,", "the kind of the object in question."},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"comment": "Some literal text.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify some constant text.",
		},
	}
	Zt_TextList = typeinfo.Flow{
		Name: "text_list",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:     "kind_name",
			Label:    "kind",
			Optional: true,
			Markup: map[string]any{
				"comment": []string{"Optionally, when the text represents the names of (existing) objects,", "the kind of the objects in question."},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "values",
			Label:   "values",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Zero or more text literals.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "Specify a list of literal text values.",
		},
	}
	Zt_RecordValue = typeinfo.Flow{
		Name: "record_value",
		Lede: "record",
		Terms: []typeinfo.Term{{
			Name: "kind_name",
			Markup: map[string]any{
				"comment": []string{"The kind of the record being constructed.", "All kinds must be pre-declared ( ex. via [DefineKind] or via jess. )"},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "fields",
			Label:   "value",
			Repeats: true,
			Markup: map[string]any{
				"comment": []string{"A set of literal values for the fields of the record.", "Any fields of the record which are not specified here,", "are \"zero initialized.\""},
			},
			Type: &Zt_FieldValue,
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
			Name: "kind_name",
			Markup: map[string]any{
				"comment": []string{"The kind of the records being constructed.", "All of the records in the list must be of the same kind.", "All kinds must be pre-declared ( ex. via [DefineKind] or via jess. )"},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "records",
			Label:   "values",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Zero or more record literals.",
			},
			Type: &Zt_FieldList,
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
	Zt_FieldList = typeinfo.Flow{
		Name: "field_list",
		Lede: "field",
		Terms: []typeinfo.Term{{
			Name:    "fields",
			Label:   "list",
			Repeats: true,
			Markup: map[string]any{
				"comment": []string{"A set of literal values for the fields of the record.", "Any fields of the record which are not specified here,", "are \"zero initialized.\""},
			},
			Type: &Zt_FieldValue,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_LiteralValue,
		},
		Markup: map[string]any{
			"comment": "A series of values used to build a record.",
		},
	}
	Zt_FieldValue = typeinfo.Flow{
		Name: "field_value",
		Lede: "field",
		Terms: []typeinfo.Term{{
			Name: "field_name",
			Markup: map[string]any{
				"comment": []string{"The name of a field in a record to initialize.", "New field names cannot be added to records at runtime;", "the field names must be part of the original declaration of the kind."},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"comment": []string{"The literal value of the field.", "The type of value must match the original declaration of the field.", "( ex. If the field was declared as a number, only a number can be used to initialize it. )"},
			},
			Type: &Zt_LiteralValue,
		}},
		Markup: map[string]any{
			"comment": "The name and value of a field used for initializing a literal record.",
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "literal",
	Comment: []string{
		"Supply constant values to runtime evaluations.",
		"( ie. a specific number when a [NumEval] is required. )",
		"",
		"Tell files support \"shortcuts\" which turn primitive values into literal commands. For instance, the number '5' in a .tell file is automatically transformed into the command \"Num value: 5\" whenever that's needed.",
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
	&Zt_NumValue,
	&Zt_NumList,
	&Zt_TextValue,
	&Zt_TextList,
	&Zt_RecordValue,
	&Zt_RecordList,
	&Zt_FieldList,
	&Zt_FieldValue,
}

// gob like registration
func Register(reg func(any)) {
	reg((*BoolValue)(nil))
	reg((*NumValue)(nil))
	reg((*NumList)(nil))
	reg((*TextValue)(nil))
	reg((*TextList)(nil))
	reg((*RecordValue)(nil))
	reg((*RecordList)(nil))
	reg((*FieldList)(nil))
	reg((*FieldValue)(nil))
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	17656638186047966738: (*FieldValue)(nil),  /* Field:value: */
	2028829358589965004:  (*BoolValue)(nil),   /* bool_eval=Bool value: */
	11511029631426206694: (*BoolValue)(nil),   /* literal_value=Bool value: */
	1426627792852548653:  (*FieldList)(nil),   /* literal_value=Field list: */
	15362209855253663632: (*NumValue)(nil),    /* literal_value=Num value: */
	16565175635984030252: (*NumValue)(nil),    /* num_eval=Num value: */
	12282038377752822419: (*NumList)(nil),     /* literal_value=Num values: */
	8089072108541894314:  (*NumList)(nil),     /* num_list_eval=Num values: */
	5076557270712812679:  (*RecordValue)(nil), /* literal_value=Record:value: */
	6692708173911561442:  (*RecordValue)(nil), /* record_eval=Record:value: */
	8711768526197034738:  (*RecordList)(nil),  /* literal_value=Record:values: */
	14652198550804167624: (*RecordList)(nil),  /* record_list_eval=Record:values: */
	9199430333197009343:  (*TextValue)(nil),   /* literal_value=Text kind:value: */
	4296855323747417954:  (*TextValue)(nil),   /* text_eval=Text kind:value: */
	1727066802786539834:  (*TextList)(nil),    /* literal_value=Text kind:values: */
	17727986415105280230: (*TextList)(nil),    /* text_list_eval=Text kind:values: */
	13114183353368545439: (*TextValue)(nil),   /* literal_value=Text value: */
	4705033170011872932:  (*TextValue)(nil),   /* text_eval=Text value: */
	2231933745037898906:  (*TextList)(nil),    /* literal_value=Text values: */
	5151885117815687006:  (*TextList)(nil),    /* text_list_eval=Text values: */
}
