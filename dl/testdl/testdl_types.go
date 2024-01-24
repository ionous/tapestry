// Code generated by Tapestry; edit at your own risk.
package testdl

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// test_slot, a type of slot.
const Z_TestSlot_Name = "test_slot"

var Z_TestSlot_T = typeinfo.Slot{
	Name: Z_TestSlot_Name,
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_TestSlot_Slot struct{ Value TestSlot }

// implements typeinfo.Inspector for a single slot.
func (*FIX_TestSlot_Slot) Inspect() typeinfo.T {
	return &Z_TestSlot_T
}

// holds a slice of slots
type TestSlot_Slots []TestSlot

// implements typeinfo.Inspector for a series of slots.
func (*TestSlot_Slots) Inspect() typeinfo.T {
	return &Z_TestSlot_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_TestEmbed struct {
	TestFlow TestFlow
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*TestEmbed) Inspect() typeinfo.T {
	return &Z_TestEmbed_T
}

// return a valid markup map, creating it if necessary.
func (op *TestEmbed) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// test_embed, a type of flow.
const Z_TestEmbed_Name = "test_embed"

// ensure the command implements its specified slots:
var _ TestSlot = (*TestEmbed)(nil)

var Z_TestEmbed_T = typeinfo.Flow{
	Name: Z_TestEmbed_Name,
	Lede: "embed",
	Terms: []typeinfo.Term{{
		Name:  "test_flow",
		Label: "test_flow",
		Type:  &Z_TestFlow_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_TestSlot_T,
	},
}

// holds a slice of type test_embed
// FIX: duplicates the spec decl.
type FIX_TestEmbed_Slice []TestEmbed

// implements typeinfo.Inspector
func (*TestEmbed_Slice) Inspect() typeinfo.T {
	return &Z_TestEmbed_T
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
func (*TestFlow) Inspect() typeinfo.T {
	return &Z_TestFlow_T
}

// return a valid markup map, creating it if necessary.
func (op *TestFlow) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// test_flow, a type of flow.
const Z_TestFlow_Name = "test_flow"

// ensure the command implements its specified slots:
var _ TestSlot = (*TestFlow)(nil)

var Z_TestFlow_T = typeinfo.Flow{
	Name: Z_TestFlow_Name,
	Lede: "flow",
	Terms: []typeinfo.Term{{
		Name:     "slot",
		Label:    "slot",
		Optional: true,
		Type:     &Z_TestSlot_T,
	}, {
		Name:     "txt",
		Label:    "txt",
		Optional: true,
		Type:     &Z_TestTxt_T,
	}, {
		Name:     "num",
		Label:    "num",
		Optional: true,
		Type:     &Z_TestNum_T,
	}, {
		Name:     "bool",
		Label:    "bool",
		Optional: true,
		Type:     &Z_TestBool_T,
	}, {
		Name:     "slots",
		Label:    "slots",
		Optional: true,
		Repeats:  true,
		Type:     &Z_TestSlot_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_TestSlot_T,
	},
}

// holds a slice of type test_flow
// FIX: duplicates the spec decl.
type FIX_TestFlow_Slice []TestFlow

// implements typeinfo.Inspector
func (*TestFlow_Slice) Inspect() typeinfo.T {
	return &Z_TestFlow_T
}

// test_bool, a type of str enum.
const Z_TestBool_Name = "test_bool"

const (
	W_TestBool_True  = "true"
	W_TestBool_False = "false"
)

var Z_TestBool_T = typeinfo.Str{
	Name: Z_TestBool_Name,
	Options: []string{
		W_TestBool_True,
		W_TestBool_False,
	},
}

// test_str, a type of str enum.
const Z_TestStr_Name = "test_str"

const (
	W_TestStr_One    = "one"
	W_TestStr_Other  = "other"
	W_TestStr_Option = "option"
)

var Z_TestStr_T = typeinfo.Str{
	Name: Z_TestStr_Name,
	Options: []string{
		W_TestStr_One,
		W_TestStr_Other,
		W_TestStr_Option,
	},
}

// test_txt, a type of str.
const Z_TestTxt_Name = "test_txt"

var Z_TestTxt_T = typeinfo.Str{
	Name: Z_TestTxt_Name,
}

// test_num, a type of num.
const Z_TestNum_Name = "test_num"

var Z_TestNum_T = typeinfo.Num{
	Name: Z_TestNum_Name,
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
	&Z_TestSlot_T,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_TestEmbed_T,
	&Z_TestFlow_T,
}

// a list of all strs in this this package
var z_str_list = []*typeinfo.Str{
	&Z_TestBool_T,
	&Z_TestStr_T,
	&Z_TestTxt_T,
}

// a list of all nums in this this package
var z_num_list = []*typeinfo.Num{
	&Z_TestNum_T,
}
