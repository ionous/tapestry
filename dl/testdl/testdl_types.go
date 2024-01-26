// Code generated by Tapestry; edit at your own risk.
package testdl

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// test_slot, a type of slot.
var Zt_TestSlot = typeinfo.Slot{
	Name: "test_slot",
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_TestSlot_Slot struct{ Value TestSlot }

// implements typeinfo.Inspector for a single slot.
func (*FIX_TestSlot_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_TestSlot, false
}

// holds a slice of slots
type TestSlot_Slots []TestSlot

// implements typeinfo.Inspector for a series of slots.
func (*TestSlot_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_TestSlot, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_TestEmbed struct {
	TestFlow TestFlow
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*TestEmbed) Inspect() (typeinfo.T, bool) {
	return &Zt_TestEmbed, false
}

// return a valid markup map, creating it if necessary.
func (op *TestEmbed) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ TestSlot = (*TestEmbed)(nil)

// test_embed, a type of flow.
var Zt_TestEmbed = typeinfo.Flow{
	Name: "test_embed",
	Lede: "embed",
	Terms: []typeinfo.Term{{
		Name:  "test_flow",
		Label: "test_flow",
		Type:  &Zt_TestFlow,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_TestSlot,
	},
}

// holds a slice of type test_embed
// FIX: duplicates the spec decl.
type FIX_TestEmbed_Slice []TestEmbed

// implements typeinfo.Inspector
func (*TestEmbed_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_TestEmbed, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_TestFlow struct {
	Slot   TestSlot
	Txt    string
	Num    float64
	Bool   string
	Slots  TestSlot
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*TestFlow) Inspect() (typeinfo.T, bool) {
	return &Zt_TestFlow, false
}

// return a valid markup map, creating it if necessary.
func (op *TestFlow) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ TestSlot = (*TestFlow)(nil)

// test_flow, a type of flow.
var Zt_TestFlow = typeinfo.Flow{
	Name: "test_flow",
	Lede: "flow",
	Terms: []typeinfo.Term{{
		Name:     "slot",
		Label:    "slot",
		Optional: true,
		Type:     &Zt_TestSlot,
	}, {
		Name:     "txt",
		Label:    "txt",
		Optional: true,
		Type:     &Zt_TestTxt,
	}, {
		Name:     "num",
		Label:    "num",
		Optional: true,
		Type:     &Zt_TestNum,
	}, {
		Name:     "bool",
		Label:    "bool",
		Optional: true,
		Type:     &Zt_TestBool,
	}, {
		Name:     "slots",
		Label:    "slots",
		Optional: true,
		Repeats:  true,
		Type:     &Zt_TestSlot,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_TestSlot,
	},
}

// holds a slice of type test_flow
// FIX: duplicates the spec decl.
type FIX_TestFlow_Slice []TestFlow

// implements typeinfo.Inspector
func (*TestFlow_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_TestFlow, true
}

const (
	Zc_TestBool_True  = "true"
	Zc_TestBool_False = "false"
)

// test_bool, a type of str enum.
var Zt_TestBool = typeinfo.Str{
	Name: "test_bool",
	Options: []string{
		Zc_TestBool_True,
		Zc_TestBool_False,
	},
}

const (
	Zc_TestStr_One    = "one"
	Zc_TestStr_Other  = "other"
	Zc_TestStr_Option = "option"
)

// test_str, a type of str enum.
var Zt_TestStr = typeinfo.Str{
	Name: "test_str",
	Options: []string{
		Zc_TestStr_One,
		Zc_TestStr_Other,
		Zc_TestStr_Option,
	},
}
var Zt_TestTxt = typeinfo.Str{
	Name: "test_txt",
}

// test_num, a type of num.
var Zt_TestNum = typeinfo.Num{
	Name: "test_num",
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "testdl",
	Slot: z_slot_list,
	Flow: z_flow_list,
	Str:  z_str_list,
	Num:  z_num_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_TestSlot,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_TestEmbed,
	&Zt_TestFlow,
}

// a list of all strs in this this package
var z_str_list = []*typeinfo.Str{
	&Zt_TestBool,
	&Zt_TestStr,
	&Zt_TestTxt,
}

// a list of all nums in this this package
var z_num_list = []*typeinfo.Num{
	&Zt_TestNum,
}
