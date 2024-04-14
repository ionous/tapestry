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

// holds a single slot.
type Address_Slot struct{ Value Address }

// implements typeinfo.Instance for a single slot.
func (*Address_Slot) TypeInfo() typeinfo.T {
	return &Zt_Address
}

// holds a slice of slots.
type Address_Slots []Address

// implements typeinfo.Instance for a series of slots.
func (*Address_Slots) TypeInfo() typeinfo.T {
	return &Zt_Address
}

// implements typeinfo.Repeats
func (op *Address_Slots) Repeats() bool {
	return len(*op) > 0
}

// dot, a type of slot.
var Zt_Dot = typeinfo.Slot{
	Name: "dot",
	Markup: map[string]any{
		"blockly-color": "MATH_HUE",
		"comment":       "Picks values from types containing other values.",
	},
}

// holds a single slot.
type Dot_Slot struct{ Value Dot }

// implements typeinfo.Instance for a single slot.
func (*Dot_Slot) TypeInfo() typeinfo.T {
	return &Zt_Dot
}

// holds a slice of slots.
type Dot_Slots []Dot

// implements typeinfo.Instance for a series of slots.
func (*Dot_Slots) TypeInfo() typeinfo.T {
	return &Zt_Dot
}

// implements typeinfo.Repeats
func (op *Dot_Slots) Repeats() bool {
	return len(*op) > 0
}

// Store a value into a variable or object.
type SetValue struct {
	Target Address
	Value  rtti.Assignment
	Markup map[string]any
}

// set_value, a type of flow.
var Zt_SetValue typeinfo.Flow

// implements typeinfo.Instance
func (*SetValue) TypeInfo() typeinfo.T {
	return &Zt_SetValue
}

// implements typeinfo.Markup
func (op *SetValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetValue)(nil)

// holds a slice of type set_value
type SetValue_Slice []SetValue

// implements typeinfo.Instance
func (*SetValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_SetValue
}

// implements typeinfo.Repeats
func (op *SetValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Set the state of an object.
type SetTrait struct {
	Target rtti.TextEval
	Trait  rtti.TextEval
	Markup map[string]any
}

// set_trait, a type of flow.
var Zt_SetTrait typeinfo.Flow

// implements typeinfo.Instance
func (*SetTrait) TypeInfo() typeinfo.T {
	return &Zt_SetTrait
}

// implements typeinfo.Markup
func (op *SetTrait) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SetTrait)(nil)

// holds a slice of type set_trait
type SetTrait_Slice []SetTrait

// implements typeinfo.Instance
func (*SetTrait_Slice) TypeInfo() typeinfo.T {
	return &Zt_SetTrait
}

// implements typeinfo.Repeats
func (op *SetTrait_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy from one stored value to another.
// Requires that the type of the two values match exactly
type CopyValue struct {
	Target Address
	Source Address
	Markup map[string]any
}

// copy_value, a type of flow.
var Zt_CopyValue typeinfo.Flow

// implements typeinfo.Instance
func (*CopyValue) TypeInfo() typeinfo.T {
	return &Zt_CopyValue
}

// implements typeinfo.Markup
func (op *CopyValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*CopyValue)(nil)

// holds a slice of type copy_value
type CopyValue_Slice []CopyValue

// implements typeinfo.Instance
func (*CopyValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_CopyValue
}

// implements typeinfo.Repeats
func (op *CopyValue_Slice) Repeats() bool {
	return len(*op) > 0
}

type ObjectRef struct {
	Name   rtti.TextEval
	Field  rtti.TextEval
	Dot    []Dot
	Markup map[string]any
}

// object_ref, a type of flow.
var Zt_ObjectRef typeinfo.Flow

// implements typeinfo.Instance
func (*ObjectRef) TypeInfo() typeinfo.T {
	return &Zt_ObjectRef
}

// implements typeinfo.Markup
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

// holds a slice of type object_ref
type ObjectRef_Slice []ObjectRef

// implements typeinfo.Instance
func (*ObjectRef_Slice) TypeInfo() typeinfo.T {
	return &Zt_ObjectRef
}

// implements typeinfo.Repeats
func (op *ObjectRef_Slice) Repeats() bool {
	return len(*op) > 0
}

type VariableRef struct {
	Name   rtti.TextEval
	Dot    []Dot
	Markup map[string]any
}

// variable_ref, a type of flow.
var Zt_VariableRef typeinfo.Flow

// implements typeinfo.Instance
func (*VariableRef) TypeInfo() typeinfo.T {
	return &Zt_VariableRef
}

// implements typeinfo.Markup
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

// holds a slice of type variable_ref
type VariableRef_Slice []VariableRef

// implements typeinfo.Instance
func (*VariableRef_Slice) TypeInfo() typeinfo.T {
	return &Zt_VariableRef
}

// implements typeinfo.Repeats
func (op *VariableRef_Slice) Repeats() bool {
	return len(*op) > 0
}

type AtField struct {
	Field  rtti.TextEval
	Markup map[string]any
}

// at_field, a type of flow.
var Zt_AtField typeinfo.Flow

// implements typeinfo.Instance
func (*AtField) TypeInfo() typeinfo.T {
	return &Zt_AtField
}

// implements typeinfo.Markup
func (op *AtField) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Dot = (*AtField)(nil)

// holds a slice of type at_field
type AtField_Slice []AtField

// implements typeinfo.Instance
func (*AtField_Slice) TypeInfo() typeinfo.T {
	return &Zt_AtField
}

// implements typeinfo.Repeats
func (op *AtField_Slice) Repeats() bool {
	return len(*op) > 0
}

type AtIndex struct {
	Index  rtti.NumberEval
	Markup map[string]any
}

// at_index, a type of flow.
var Zt_AtIndex typeinfo.Flow

// implements typeinfo.Instance
func (*AtIndex) TypeInfo() typeinfo.T {
	return &Zt_AtIndex
}

// implements typeinfo.Markup
func (op *AtIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Dot = (*AtIndex)(nil)

// holds a slice of type at_index
type AtIndex_Slice []AtIndex

// implements typeinfo.Instance
func (*AtIndex_Slice) TypeInfo() typeinfo.T {
	return &Zt_AtIndex
}

// implements typeinfo.Repeats
func (op *AtIndex_Slice) Repeats() bool {
	return len(*op) > 0
}

// Executes a pattern, and potentially returns a value.
type CallPattern struct {
	PatternName string
	Arguments   []Arg
	Markup      map[string]any
}

// call_pattern, a type of flow.
var Zt_CallPattern typeinfo.Flow

// implements typeinfo.Instance
func (*CallPattern) TypeInfo() typeinfo.T {
	return &Zt_CallPattern
}

// implements typeinfo.Markup
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

// holds a slice of type call_pattern
type CallPattern_Slice []CallPattern

// implements typeinfo.Instance
func (*CallPattern_Slice) TypeInfo() typeinfo.T {
	return &Zt_CallPattern
}

// implements typeinfo.Repeats
func (op *CallPattern_Slice) Repeats() bool {
	return len(*op) > 0
}

// Runtime version of argument.
type Arg struct {
	Name   string
	Value  rtti.Assignment
	Markup map[string]any
}

// arg, a type of flow.
var Zt_Arg typeinfo.Flow

// implements typeinfo.Instance
func (*Arg) TypeInfo() typeinfo.T {
	return &Zt_Arg
}

// implements typeinfo.Markup
func (op *Arg) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type arg
type Arg_Slice []Arg

// implements typeinfo.Instance
func (*Arg_Slice) TypeInfo() typeinfo.T {
	return &Zt_Arg
}

// implements typeinfo.Repeats
func (op *Arg_Slice) Repeats() bool {
	return len(*op) > 0
}

// Adapts an execute statement to an assignment.
// Used internally for package shuttle.
type FromExe struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// from_exe, a type of flow.
var Zt_FromExe typeinfo.Flow

// implements typeinfo.Instance
func (*FromExe) TypeInfo() typeinfo.T {
	return &Zt_FromExe
}

// implements typeinfo.Markup
func (op *FromExe) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromExe)(nil)

// holds a slice of type from_exe
type FromExe_Slice []FromExe

// implements typeinfo.Instance
func (*FromExe_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromExe
}

// implements typeinfo.Repeats
func (op *FromExe_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a boolean value.
type FromBool struct {
	Value  rtti.BoolEval
	Markup map[string]any
}

// from_bool, a type of flow.
var Zt_FromBool typeinfo.Flow

// implements typeinfo.Instance
func (*FromBool) TypeInfo() typeinfo.T {
	return &Zt_FromBool
}

// implements typeinfo.Markup
func (op *FromBool) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromBool)(nil)

// holds a slice of type from_bool
type FromBool_Slice []FromBool

// implements typeinfo.Instance
func (*FromBool_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromBool
}

// implements typeinfo.Repeats
func (op *FromBool_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a number.
type FromNumber struct {
	Value  rtti.NumberEval
	Markup map[string]any
}

// from_number, a type of flow.
var Zt_FromNumber typeinfo.Flow

// implements typeinfo.Instance
func (*FromNumber) TypeInfo() typeinfo.T {
	return &Zt_FromNumber
}

// implements typeinfo.Markup
func (op *FromNumber) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumber)(nil)

// holds a slice of type from_number
type FromNumber_Slice []FromNumber

// implements typeinfo.Instance
func (*FromNumber_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromNumber
}

// implements typeinfo.Repeats
func (op *FromNumber_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a text string.
type FromText struct {
	Value  rtti.TextEval
	Markup map[string]any
}

// from_text, a type of flow.
var Zt_FromText typeinfo.Flow

// implements typeinfo.Instance
func (*FromText) TypeInfo() typeinfo.T {
	return &Zt_FromText
}

// implements typeinfo.Markup
func (op *FromText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromText)(nil)

// holds a slice of type from_text
type FromText_Slice []FromText

// implements typeinfo.Instance
func (*FromText_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromText
}

// implements typeinfo.Repeats
func (op *FromText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a record.
type FromRecord struct {
	Value  rtti.RecordEval
	Markup map[string]any
}

// from_record, a type of flow.
var Zt_FromRecord typeinfo.Flow

// implements typeinfo.Instance
func (*FromRecord) TypeInfo() typeinfo.T {
	return &Zt_FromRecord
}

// implements typeinfo.Markup
func (op *FromRecord) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecord)(nil)

// holds a slice of type from_record
type FromRecord_Slice []FromRecord

// implements typeinfo.Instance
func (*FromRecord_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromRecord
}

// implements typeinfo.Repeats
func (op *FromRecord_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a list of numbers.
type FromNumList struct {
	Value  rtti.NumListEval
	Markup map[string]any
}

// from_num_list, a type of flow.
var Zt_FromNumList typeinfo.Flow

// implements typeinfo.Instance
func (*FromNumList) TypeInfo() typeinfo.T {
	return &Zt_FromNumList
}

// implements typeinfo.Markup
func (op *FromNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromNumList)(nil)

// holds a slice of type from_num_list
type FromNumList_Slice []FromNumList

// implements typeinfo.Instance
func (*FromNumList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromNumList
}

// implements typeinfo.Repeats
func (op *FromNumList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a list of text strings.
type FromTextList struct {
	Value  rtti.TextListEval
	Markup map[string]any
}

// from_text_list, a type of flow.
var Zt_FromTextList typeinfo.Flow

// implements typeinfo.Instance
func (*FromTextList) TypeInfo() typeinfo.T {
	return &Zt_FromTextList
}

// implements typeinfo.Markup
func (op *FromTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromTextList)(nil)

// holds a slice of type from_text_list
type FromTextList_Slice []FromTextList

// implements typeinfo.Instance
func (*FromTextList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromTextList
}

// implements typeinfo.Repeats
func (op *FromTextList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Calculates a list of records.
type FromRecordList struct {
	Value  rtti.RecordListEval
	Markup map[string]any
}

// from_record_list, a type of flow.
var Zt_FromRecordList typeinfo.Flow

// implements typeinfo.Instance
func (*FromRecordList) TypeInfo() typeinfo.T {
	return &Zt_FromRecordList
}

// implements typeinfo.Markup
func (op *FromRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Assignment = (*FromRecordList)(nil)

// holds a slice of type from_record_list
type FromRecordList_Slice []FromRecordList

// implements typeinfo.Instance
func (*FromRecordList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromRecordList
}

// implements typeinfo.Repeats
func (op *FromRecordList_Slice) Repeats() bool {
	return len(*op) > 0
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
var z_signatures = map[uint64]typeinfo.Instance{
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

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_SetValue = typeinfo.Flow{
		Name: "set_value",
		Lede: "set",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &Zt_Address,
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
	Zt_SetTrait = typeinfo.Flow{
		Name: "set_trait",
		Lede: "set",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &rtti.Zt_TextEval,
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
	Zt_CopyValue = typeinfo.Flow{
		Name: "copy_value",
		Lede: "copy",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &Zt_Address,
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
	Zt_ObjectRef = typeinfo.Flow{
		Name: "object_ref",
		Lede: "object",
		Terms: []typeinfo.Term{{
			Name: "name",
			Type: &rtti.Zt_TextEval,
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
	Zt_VariableRef = typeinfo.Flow{
		Name: "variable_ref",
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
	Zt_AtField = typeinfo.Flow{
		Name: "at_field",
		Lede: "at_field",
		Terms: []typeinfo.Term{{
			Name: "field",
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Dot,
		},
	}
	Zt_AtIndex = typeinfo.Flow{
		Name: "at_index",
		Lede: "at_index",
		Terms: []typeinfo.Term{{
			Name: "index",
			Type: &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Dot,
		},
	}
	Zt_CallPattern = typeinfo.Flow{
		Name: "call_pattern",
		Lede: "determine",
		Terms: []typeinfo.Term{{
			Name: "pattern_name",
			Type: &prim.Zt_Text,
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
	Zt_Arg = typeinfo.Flow{
		Name: "arg",
		Lede: "arg",
		Terms: []typeinfo.Term{{
			Name: "name",
			Type: &prim.Zt_Text,
		}, {
			Name:  "value",
			Label: "from",
			Type:  &rtti.Zt_Assignment,
		}},
		Markup: map[string]any{
			"comment": "Runtime version of argument.",
		},
	}
	Zt_FromExe = typeinfo.Flow{
		Name: "from_exe",
		Lede: "from_exe",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Repeats: true,
			Type:    &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Adapts an execute statement to an assignment.", "Used internally for package shuttle."},
		},
	}
	Zt_FromBool = typeinfo.Flow{
		Name: "from_bool",
		Lede: "from_bool",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a boolean value.",
		},
	}
	Zt_FromNumber = typeinfo.Flow{
		Name: "from_number",
		Lede: "from_number",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a number.",
		},
	}
	Zt_FromText = typeinfo.Flow{
		Name: "from_text",
		Lede: "from_text",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a text string.",
		},
	}
	Zt_FromRecord = typeinfo.Flow{
		Name: "from_record",
		Lede: "from_record",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_RecordEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a record.",
		},
	}
	Zt_FromNumList = typeinfo.Flow{
		Name: "from_num_list",
		Lede: "from_num_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_NumListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a list of numbers.",
		},
	}
	Zt_FromTextList = typeinfo.Flow{
		Name: "from_text_list",
		Lede: "from_text_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_TextListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a list of text strings.",
		},
	}
	Zt_FromRecordList = typeinfo.Flow{
		Name: "from_record_list",
		Lede: "from_record_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_RecordListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Calculates a list of records.",
		},
	}
}
