// Pattern calls and related helpers.
package call

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// trigger, a type of slot.
var Zt_Trigger = typeinfo.Slot{
	Name: "trigger",
	Markup: map[string]any{
		"comment":  "Helper for counting values.",
		"internal": true,
	},
}

// Holds a single slot.
type Trigger_Slot struct{ Value Trigger }

// Implements [typeinfo.Instance] for a single slot.
func (*Trigger_Slot) TypeInfo() typeinfo.T {
	return &Zt_Trigger
}

// Holds a slice of slots.
type Trigger_Slots []Trigger

// Implements [typeinfo.Instance] for a slice of slots.
func (*Trigger_Slots) TypeInfo() typeinfo.T {
	return &Zt_Trigger
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Trigger_Slots) Repeats() bool {
	return len(*op) > 0
}

// Determine whether a scene (aka domain) is active.
type ActiveScene struct {
	Name   string
	Markup map[string]any
}

// active_scene, a type of flow.
var Zt_ActiveScene typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ActiveScene) TypeInfo() typeinfo.T {
	return &Zt_ActiveScene
}

// Implements [typeinfo.Markup]
func (op *ActiveScene) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*ActiveScene)(nil)

// Holds a slice of type ActiveScene.
type ActiveScene_Slice []ActiveScene

// Implements [typeinfo.Instance] for a slice of ActiveScene.
func (*ActiveScene_Slice) TypeInfo() typeinfo.T {
	return &Zt_ActiveScene
}

// Implements [typeinfo.Repeats] for a slice of ActiveScene.
func (op *ActiveScene_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether a pattern is running.
type ActivePattern struct {
	PatternName string
	Markup      map[string]any
}

// active_pattern, a type of flow.
var Zt_ActivePattern typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ActivePattern) TypeInfo() typeinfo.T {
	return &Zt_ActivePattern
}

// Implements [typeinfo.Markup]
func (op *ActivePattern) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*ActivePattern)(nil)
var _ rtti.NumEval = (*ActivePattern)(nil)

// Holds a slice of type ActivePattern.
type ActivePattern_Slice []ActivePattern

// Implements [typeinfo.Instance] for a slice of ActivePattern.
func (*ActivePattern_Slice) TypeInfo() typeinfo.T {
	return &Zt_ActivePattern
}

// Implements [typeinfo.Repeats] for a slice of ActivePattern.
func (op *ActivePattern_Slice) Repeats() bool {
	return len(*op) > 0
}

// Pass a named value to a parameterized call.
type Arg struct {
	Name   string
	Value  rtti.Assignment
	Markup map[string]any
}

// arg, a type of flow.
var Zt_Arg typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Arg) TypeInfo() typeinfo.T {
	return &Zt_Arg
}

// Implements [typeinfo.Markup]
func (op *Arg) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Holds a slice of type Arg.
type Arg_Slice []Arg

// Implements [typeinfo.Instance] for a slice of Arg.
func (*Arg_Slice) TypeInfo() typeinfo.T {
	return &Zt_Arg
}

// Implements [typeinfo.Repeats] for a slice of Arg.
func (op *Arg_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide one or more execute commands for an assignment.
// Used internally for jess rules.
type FromExe struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// from_exe, a type of flow.
var Zt_FromExe typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromExe) TypeInfo() typeinfo.T {
	return &Zt_FromExe
}

// Implements [typeinfo.Markup]
func (op *FromExe) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromExe)(nil)

// Holds a slice of type FromExe.
type FromExe_Slice []FromExe

// Implements [typeinfo.Instance] for a slice of FromExe.
func (*FromExe_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromExe
}

// Implements [typeinfo.Repeats] for a slice of FromExe.
func (op *FromExe_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a stored value for an assignment.
type FromAddress struct {
	Value  rtti.Address
	Markup map[string]any
}

// from_address, a type of flow.
var Zt_FromAddress typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromAddress) TypeInfo() typeinfo.T {
	return &Zt_FromAddress
}

// Implements [typeinfo.Markup]
func (op *FromAddress) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromAddress)(nil)

// Holds a slice of type FromAddress.
type FromAddress_Slice []FromAddress

// Implements [typeinfo.Instance] for a slice of FromAddress.
func (*FromAddress_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromAddress
}

// Implements [typeinfo.Repeats] for a slice of FromAddress.
func (op *FromAddress_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a boolean value for an assignment.
type FromBool struct {
	Value  rtti.BoolEval
	Markup map[string]any
}

// from_bool, a type of flow.
var Zt_FromBool typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromBool) TypeInfo() typeinfo.T {
	return &Zt_FromBool
}

// Implements [typeinfo.Markup]
func (op *FromBool) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromBool)(nil)

// Holds a slice of type FromBool.
type FromBool_Slice []FromBool

// Implements [typeinfo.Instance] for a slice of FromBool.
func (*FromBool_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromBool
}

// Implements [typeinfo.Repeats] for a slice of FromBool.
func (op *FromBool_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a number for an assignment.
type FromNum struct {
	Value  rtti.NumEval
	Markup map[string]any
}

// from_num, a type of flow.
var Zt_FromNum typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromNum) TypeInfo() typeinfo.T {
	return &Zt_FromNum
}

// Implements [typeinfo.Markup]
func (op *FromNum) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromNum)(nil)

// Holds a slice of type FromNum.
type FromNum_Slice []FromNum

// Implements [typeinfo.Instance] for a slice of FromNum.
func (*FromNum_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromNum
}

// Implements [typeinfo.Repeats] for a slice of FromNum.
func (op *FromNum_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide some text for an assignment.
type FromText struct {
	Value  rtti.TextEval
	Markup map[string]any
}

// from_text, a type of flow.
var Zt_FromText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromText) TypeInfo() typeinfo.T {
	return &Zt_FromText
}

// Implements [typeinfo.Markup]
func (op *FromText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromText)(nil)

// Holds a slice of type FromText.
type FromText_Slice []FromText

// Implements [typeinfo.Instance] for a slice of FromText.
func (*FromText_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromText
}

// Implements [typeinfo.Repeats] for a slice of FromText.
func (op *FromText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a record for an assignment.
type FromRecord struct {
	Value  rtti.RecordEval
	Markup map[string]any
}

// from_record, a type of flow.
var Zt_FromRecord typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromRecord) TypeInfo() typeinfo.T {
	return &Zt_FromRecord
}

// Implements [typeinfo.Markup]
func (op *FromRecord) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromRecord)(nil)

// Holds a slice of type FromRecord.
type FromRecord_Slice []FromRecord

// Implements [typeinfo.Instance] for a slice of FromRecord.
func (*FromRecord_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromRecord
}

// Implements [typeinfo.Repeats] for a slice of FromRecord.
func (op *FromRecord_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a list of numbers for an assignment.
type FromNumList struct {
	Value  rtti.NumListEval
	Markup map[string]any
}

// from_num_list, a type of flow.
var Zt_FromNumList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromNumList) TypeInfo() typeinfo.T {
	return &Zt_FromNumList
}

// Implements [typeinfo.Markup]
func (op *FromNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromNumList)(nil)

// Holds a slice of type FromNumList.
type FromNumList_Slice []FromNumList

// Implements [typeinfo.Instance] for a slice of FromNumList.
func (*FromNumList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromNumList
}

// Implements [typeinfo.Repeats] for a slice of FromNumList.
func (op *FromNumList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a list of text values for an assignment.
type FromTextList struct {
	Value  rtti.TextListEval
	Markup map[string]any
}

// from_text_list, a type of flow.
var Zt_FromTextList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromTextList) TypeInfo() typeinfo.T {
	return &Zt_FromTextList
}

// Implements [typeinfo.Markup]
func (op *FromTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromTextList)(nil)

// Holds a slice of type FromTextList.
type FromTextList_Slice []FromTextList

// Implements [typeinfo.Instance] for a slice of FromTextList.
func (*FromTextList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromTextList
}

// Implements [typeinfo.Repeats] for a slice of FromTextList.
func (op *FromTextList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Provide a list of records for an assignment.
type FromRecordList struct {
	Value  rtti.RecordListEval
	Markup map[string]any
}

// from_record_list, a type of flow.
var Zt_FromRecordList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FromRecordList) TypeInfo() typeinfo.T {
	return &Zt_FromRecordList
}

// Implements [typeinfo.Markup]
func (op *FromRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Assignment = (*FromRecordList)(nil)

// Holds a slice of type FromRecordList.
type FromRecordList_Slice []FromRecordList

// Implements [typeinfo.Instance] for a slice of FromRecordList.
func (*FromRecordList_Slice) TypeInfo() typeinfo.T {
	return &Zt_FromRecordList
}

// Implements [typeinfo.Repeats] for a slice of FromRecordList.
func (op *FromRecordList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Run a pattern, returning its result (if any).
// Tell files support calling patterns directly, so this is only needed when using the blockly editor.
type CallPattern struct {
	PatternName string
	Arguments   []Arg
	Markup      map[string]any
}

// call_pattern, a type of flow.
var Zt_CallPattern typeinfo.Flow

// Implements [typeinfo.Instance]
func (*CallPattern) TypeInfo() typeinfo.T {
	return &Zt_CallPattern
}

// Implements [typeinfo.Markup]
func (op *CallPattern) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*CallPattern)(nil)
var _ rtti.BoolEval = (*CallPattern)(nil)
var _ rtti.NumEval = (*CallPattern)(nil)
var _ rtti.TextEval = (*CallPattern)(nil)
var _ rtti.RecordEval = (*CallPattern)(nil)
var _ rtti.NumListEval = (*CallPattern)(nil)
var _ rtti.TextListEval = (*CallPattern)(nil)
var _ rtti.RecordListEval = (*CallPattern)(nil)

// Holds a slice of type CallPattern.
type CallPattern_Slice []CallPattern

// Implements [typeinfo.Instance] for a slice of CallPattern.
func (*CallPattern_Slice) TypeInfo() typeinfo.T {
	return &Zt_CallPattern
}

// Implements [typeinfo.Repeats] for a slice of CallPattern.
func (op *CallPattern_Slice) Repeats() bool {
	return len(*op) > 0
}

// Runtime version of count_of.
// A guard which returns true based on a counter.
type CallTrigger struct {
	Name    string
	Trigger Trigger
	Num     rtti.NumEval
	Markup  map[string]any
}

// call_trigger, a type of flow.
var Zt_CallTrigger typeinfo.Flow

// Implements [typeinfo.Instance]
func (*CallTrigger) TypeInfo() typeinfo.T {
	return &Zt_CallTrigger
}

// Implements [typeinfo.Markup]
func (op *CallTrigger) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*CallTrigger)(nil)

// Holds a slice of type CallTrigger.
type CallTrigger_Slice []CallTrigger

// Implements [typeinfo.Instance] for a slice of CallTrigger.
func (*CallTrigger_Slice) TypeInfo() typeinfo.T {
	return &Zt_CallTrigger
}

// Implements [typeinfo.Repeats] for a slice of CallTrigger.
func (op *CallTrigger_Slice) Repeats() bool {
	return len(*op) > 0
}

// call_trigger
type TriggerCycle struct {
	Markup map[string]any
}

// trigger_cycle, a type of flow.
var Zt_TriggerCycle typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TriggerCycle) TypeInfo() typeinfo.T {
	return &Zt_TriggerCycle
}

// Implements [typeinfo.Markup]
func (op *TriggerCycle) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Trigger = (*TriggerCycle)(nil)

// Holds a slice of type TriggerCycle.
type TriggerCycle_Slice []TriggerCycle

// Implements [typeinfo.Instance] for a slice of TriggerCycle.
func (*TriggerCycle_Slice) TypeInfo() typeinfo.T {
	return &Zt_TriggerCycle
}

// Implements [typeinfo.Repeats] for a slice of TriggerCycle.
func (op *TriggerCycle_Slice) Repeats() bool {
	return len(*op) > 0
}

// call_trigger
type TriggerOnce struct {
	Markup map[string]any
}

// trigger_once, a type of flow.
var Zt_TriggerOnce typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TriggerOnce) TypeInfo() typeinfo.T {
	return &Zt_TriggerOnce
}

// Implements [typeinfo.Markup]
func (op *TriggerOnce) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Trigger = (*TriggerOnce)(nil)

// Holds a slice of type TriggerOnce.
type TriggerOnce_Slice []TriggerOnce

// Implements [typeinfo.Instance] for a slice of TriggerOnce.
func (*TriggerOnce_Slice) TypeInfo() typeinfo.T {
	return &Zt_TriggerOnce
}

// Implements [typeinfo.Repeats] for a slice of TriggerOnce.
func (op *TriggerOnce_Slice) Repeats() bool {
	return len(*op) > 0
}

// call_trigger
type TriggerSwitch struct {
	Markup map[string]any
}

// trigger_switch, a type of flow.
var Zt_TriggerSwitch typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TriggerSwitch) TypeInfo() typeinfo.T {
	return &Zt_TriggerSwitch
}

// Implements [typeinfo.Markup]
func (op *TriggerSwitch) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Trigger = (*TriggerSwitch)(nil)

// Holds a slice of type TriggerSwitch.
type TriggerSwitch_Slice []TriggerSwitch

// Implements [typeinfo.Instance] for a slice of TriggerSwitch.
func (*TriggerSwitch_Slice) TypeInfo() typeinfo.T {
	return &Zt_TriggerSwitch
}

// Implements [typeinfo.Repeats] for a slice of TriggerSwitch.
func (op *TriggerSwitch_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_ActiveScene = typeinfo.Flow{
		Name: "active_scene",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "name",
			Label: "scene",
			Markup: map[string]any{
				"comment": "The name of the scene to check.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": "Determine whether a scene (aka domain) is active.",
		},
	}
	Zt_ActivePattern = typeinfo.Flow{
		Name: "active_pattern",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "pattern_name",
			Label: "pattern",
			Markup: map[string]any{
				"comment": "The name of the pattern to check.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumEval,
		},
		Markup: map[string]any{
			"comment": "Determine whether a pattern is running.",
		},
	}
	Zt_Arg = typeinfo.Flow{
		Name: "arg",
		Lede: "arg",
		Terms: []typeinfo.Term{{
			Name: "name",
			Markup: map[string]any{
				"comment": "Name of the parameter. An empty string is treated as an unnamed parameter.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "value",
			Label: "from",
			Markup: map[string]any{
				"comment": "Value to assign to the parameter.",
			},
			Type: &rtti.Zt_Assignment,
		}},
		Markup: map[string]any{
			"comment": "Pass a named value to a parameterized call.",
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
			"comment":  []interface{}{"Provide one or more execute commands for an assignment.", "Used internally for jess rules."},
			"internal": true,
		},
	}
	Zt_FromAddress = typeinfo.Flow{
		Name: "from_address",
		Lede: "from_address",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Address to read from.",
			},
			Type: &rtti.Zt_Address,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a stored value for an assignment.",
		},
	}
	Zt_FromBool = typeinfo.Flow{
		Name: "from_bool",
		Lede: "from_bool",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Boolean value for the assignment.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a boolean value for an assignment.",
		},
	}
	Zt_FromNum = typeinfo.Flow{
		Name: "from_num",
		Lede: "from_num",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Number for the assignment.",
			},
			Type: &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a number for an assignment.",
		},
	}
	Zt_FromText = typeinfo.Flow{
		Name: "from_text",
		Lede: "from_text",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Text for the assignment.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide some text for an assignment.",
		},
	}
	Zt_FromRecord = typeinfo.Flow{
		Name: "from_record",
		Lede: "from_record",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Record for the assignment.",
			},
			Type: &rtti.Zt_RecordEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a record for an assignment.",
		},
	}
	Zt_FromNumList = typeinfo.Flow{
		Name: "from_num_list",
		Lede: "from_num_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Numbers for the assignment.",
			},
			Type: &rtti.Zt_NumListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a list of numbers for an assignment.",
		},
	}
	Zt_FromTextList = typeinfo.Flow{
		Name: "from_text_list",
		Lede: "from_text_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Text values for the assignment.",
			},
			Type: &rtti.Zt_TextListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a list of text values for an assignment.",
		},
	}
	Zt_FromRecordList = typeinfo.Flow{
		Name: "from_record_list",
		Lede: "from_record_list",
		Terms: []typeinfo.Term{{
			Name: "value",
			Markup: map[string]any{
				"comment": "Record values for the assignment.",
			},
			Type: &rtti.Zt_RecordListEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Assignment,
		},
		Markup: map[string]any{
			"comment": "Provide a list of records for an assignment.",
		},
	}
	Zt_CallPattern = typeinfo.Flow{
		Name: "call_pattern",
		Lede: "determine",
		Terms: []typeinfo.Term{{
			Name: "pattern_name",
			Markup: map[string]any{
				"comment": "The name of the pattern to run.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "arguments",
			Label:   "args",
			Repeats: true,
			Markup: map[string]any{
				"comment": []interface{}{"Arguments to pass to the pattern.", "Any unnamed arguments must proceed all named arguments. Unnamed arguments are assigned to parameters in the order the parameters were declared. It's considered an error to assign the same parameter multiple times."},
			},
			Type: &Zt_Arg,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumEval,
			&rtti.Zt_TextEval,
			&rtti.Zt_RecordEval,
			&rtti.Zt_NumListEval,
			&rtti.Zt_TextListEval,
			&rtti.Zt_RecordListEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Run a pattern, returning its result (if any).", "Tell files support calling patterns directly, so this is only needed when using the blockly editor."},
		},
	}
	Zt_CallTrigger = typeinfo.Flow{
		Name: "call_trigger",
		Lede: "trigger",
		Terms: []typeinfo.Term{{
			Name: "name",
			Type: &prim.Zt_Text,
		}, {
			Name:  "trigger",
			Label: "on",
			Type:  &Zt_Trigger,
		}, {
			Name:  "num",
			Label: "num",
			Type:  &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment":  []interface{}{"Runtime version of count_of.", "A guard which returns true based on a counter."},
			"internal": true,
		},
	}
	Zt_TriggerCycle = typeinfo.Flow{
		Name:  "trigger_cycle",
		Lede:  "every",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&Zt_Trigger,
		},
		Markup: map[string]any{
			"comment":  "call_trigger",
			"internal": true,
		},
	}
	Zt_TriggerOnce = typeinfo.Flow{
		Name:  "trigger_once",
		Lede:  "at",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&Zt_Trigger,
		},
		Markup: map[string]any{
			"comment":  "call_trigger",
			"internal": true,
		},
	}
	Zt_TriggerSwitch = typeinfo.Flow{
		Name:  "trigger_switch",
		Lede:  "after",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&Zt_Trigger,
		},
		Markup: map[string]any{
			"comment":  "call_trigger",
			"internal": true,
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "call",
	Comment: []string{
		"Pattern calls and related helpers.",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Trigger,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_ActiveScene,
	&Zt_ActivePattern,
	&Zt_Arg,
	&Zt_FromExe,
	&Zt_FromAddress,
	&Zt_FromBool,
	&Zt_FromNum,
	&Zt_FromText,
	&Zt_FromRecord,
	&Zt_FromNumList,
	&Zt_FromTextList,
	&Zt_FromRecordList,
	&Zt_CallPattern,
	&Zt_CallTrigger,
	&Zt_TriggerCycle,
	&Zt_TriggerOnce,
	&Zt_TriggerSwitch,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	6291103735245333139:  (*Arg)(nil),            /* Arg:from: */
	9392469773844077696:  (*TriggerSwitch)(nil),  /* trigger=After */
	2233111806717201007:  (*TriggerOnce)(nil),    /* trigger=At */
	5430006510328108403:  (*CallPattern)(nil),    /* bool_eval=Determine:args: */
	11666175118824200195: (*CallPattern)(nil),    /* execute=Determine:args: */
	9675109928599400849:  (*CallPattern)(nil),    /* num_eval=Determine:args: */
	16219448703619493492: (*CallPattern)(nil),    /* num_list_eval=Determine:args: */
	13992013847750998452: (*CallPattern)(nil),    /* record_eval=Determine:args: */
	352268441608212603:   (*CallPattern)(nil),    /* record_list_eval=Determine:args: */
	5079530186593846942:  (*CallPattern)(nil),    /* text_eval=Determine:args: */
	13938609641525654217: (*CallPattern)(nil),    /* text_list_eval=Determine:args: */
	1457631626735043065:  (*TriggerCycle)(nil),   /* trigger=Every */
	9651737781749814793:  (*FromAddress)(nil),    /* assignment=FromAddress: */
	16065241269206568079: (*FromBool)(nil),       /* assignment=FromBool: */
	9721304908210135401:  (*FromExe)(nil),        /* assignment=FromExe: */
	13937541344191718121: (*FromNum)(nil),        /* assignment=FromNum: */
	15276643347016776669: (*FromNumList)(nil),    /* assignment=FromNumList: */
	8445595699766392240:  (*FromRecord)(nil),     /* assignment=FromRecord: */
	17510952281883199828: (*FromRecordList)(nil), /* assignment=FromRecordList: */
	9783457335751138546:  (*FromText)(nil),       /* assignment=FromText: */
	3267530751198060154:  (*FromTextList)(nil),   /* assignment=FromTextList: */
	10847423070654993213: (*ActivePattern)(nil),  /* bool_eval=Is pattern: */
	15097434718788250679: (*ActivePattern)(nil),  /* num_eval=Is pattern: */
	2257319580031922583:  (*ActiveScene)(nil),    /* bool_eval=Is scene: */
	2711869841453509536:  (*CallTrigger)(nil),    /* bool_eval=Trigger:on:num: */
}
