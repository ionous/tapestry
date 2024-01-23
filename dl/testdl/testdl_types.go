// Code generated by Tapestry; edit at your own risk.
package testdl

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// test_slot, a type of slot.
const Z_TestSlot_Type = "test_slot"

var Z_TestSlot_Info = typeinfo.Slot{
	Name: Z_TestSlot_Type,
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_TestSlot_Slot struct{ Value TestSlot }

// implements typeinfo.Inspector for a single slot.
func (*FIX_TestSlot_Slot) Inspect() typeinfo.T {
	return &Z_TestSlot_Info
}

// holds a slice of slots
type TestSlot_Slots []TestSlot

// implements typeinfo.Inspector for a series of slots.
func (*TestSlot_Slots) Inspect() typeinfo.T {
	return &Z_TestSlot_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_TestEmbed struct {
	TestFlow TestFlow
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*TestEmbed) Inspect() typeinfo.T {
	return &Z_TestEmbed_Info
}

// return a valid markup map, creating it if necessary.
func (op *TestEmbed) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// test_embed, a type of flow.
const Z_TestEmbed_Type = "test_embed"

// ensure the command implements its specified slots:
var _ TestSlot = (*TestEmbed)(nil)

var Z_TestEmbed_Info = typeinfo.Flow{
	Name: Z_TestEmbed_Type,
	Lede: "embed",
	Terms: []typeinfo.Term{{
		Name:  "test_flow",
		Label: "test_flow",
		Type:  &Z_TestFlow_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_TestSlot_Info,
	},
}

// holds a slice of type test_embed
// FIX: duplicates the spec decl.
type FIX_TestEmbed_Slice []TestEmbed

// implements typeinfo.Inspector
func (*TestEmbed_Slice) Inspect() typeinfo.T {
	return &Z_TestEmbed_Info
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
	return &Z_TestFlow_Info
}

// return a valid markup map, creating it if necessary.
func (op *TestFlow) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// test_flow, a type of flow.
const Z_TestFlow_Type = "test_flow"

// ensure the command implements its specified slots:
var _ TestSlot = (*TestFlow)(nil)

var Z_TestFlow_Info = typeinfo.Flow{
	Name: Z_TestFlow_Type,
	Lede: "flow",
	Terms: []typeinfo.Term{{
		Name:     "slot",
		Label:    "slot",
		Optional: true,
		Type:     &Z_TestSlot_Info,
	}, {
		Name:     "txt",
		Label:    "txt",
		Optional: true,
		Type:     &Z_TestTxt_Info,
	}, {
		Name:     "num",
		Label:    "num",
		Optional: true,
		Type:     &Z_TestNum_Info,
	}, {
		Name:     "bool",
		Label:    "bool",
		Optional: true,
		Type:     &Z_TestBool_Info,
	}, {
		Name:     "slots",
		Label:    "slots",
		Optional: true,
		Repeats:  true,
		Type:     &Z_TestSlot_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_TestSlot_Info,
	},
}

// holds a slice of type test_flow
// FIX: duplicates the spec decl.
type FIX_TestFlow_Slice []TestFlow

// implements typeinfo.Inspector
func (*TestFlow_Slice) Inspect() typeinfo.T {
	return &Z_TestFlow_Info
}

// test_bool, a type of str enum.
const Z_TestBool_Type = "test_bool"

const (
	W_TestBool_True  = "$TRUE"
	W_TestBool_False = "$FALSE"
)

var Z_TestBool_Info = typeinfo.Str{
	Name: Z_TestBool_Type,
	Options: []string{
		W_TestBool_True,
		W_TestBool_False,
	},
}

// test_str, a type of str enum.
const Z_TestStr_Type = "test_str"

const (
	W_TestStr_One    = "$ONE"
	W_TestStr_Other  = "$OTHER"
	W_TestStr_Option = "$OPTION"
)

var Z_TestStr_Info = typeinfo.Str{
	Name: Z_TestStr_Type,
	Options: []string{
		W_TestStr_One,
		W_TestStr_Other,
		W_TestStr_Option,
	},
}

// test_txt, a type of str.
const Z_TestTxt_Type = "test_txt"

var Z_TestTxt_Info = typeinfo.Str{
	Name: Z_TestTxt_Type,
}

// test_num, a type of num.
const Z_TestNum_Type = "test_num"

var Z_TestNum_Info = typeinfo.Num{
	Name: Z_TestNum_Type,
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "testdl",
	Slot: z_slot_list,
	Flow: z_flow_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_TestSlot_Info,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_TestEmbed_Info,
	&Z_TestFlow_Info,
}
