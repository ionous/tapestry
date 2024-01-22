// Code generated by Tapestry; edit at your own risk.
package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_EraseEdge struct {
	Target assign.Address
	AtEdge rtti.BoolEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*EraseEdge) Inspect() typeinfo.T {
	return &Z_EraseEdge_Info
}

// return a valid markup map, creating it if necessary.
func (op *EraseEdge) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// erase_edge, a type of flow.
const Z_EraseEdge_Type = "erase_edge"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*EraseEdge)(nil)

var Z_EraseEdge_Info = typeinfo.Flow{
	Name: Z_EraseEdge_Type,
	Lede: "erase",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:     "at_edge",
		Label:    "at_front",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Erase at edge: Remove one or more values from a list.",
	},
}

// holds a slice of type erase_edge
// FIX: duplicates the spec decl.
type FIX_EraseEdge_Slice []EraseEdge

// implements typeinfo.Inspector
func (*EraseEdge_Slice) Inspect() typeinfo.T {
	return &Z_EraseEdge_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_EraseIndex struct {
	Count   rtti.NumberEval
	Target  assign.Address
	AtIndex rtti.NumberEval
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*EraseIndex) Inspect() typeinfo.T {
	return &Z_EraseIndex_Info
}

// return a valid markup map, creating it if necessary.
func (op *EraseIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// erase_index, a type of flow.
const Z_EraseIndex_Type = "erase_index"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*EraseIndex)(nil)

var Z_EraseIndex_Info = typeinfo.Flow{
	Name: Z_EraseIndex_Type,
	Lede: "erase",
	Terms: []typeinfo.Term{{
		Name:  "count",
		Label: "_",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:  "target",
		Label: "from",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "at_index",
		Label: "at_index",
		Type:  &rtti.Z_NumberEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Erase at index: Remove one or more values from a list.",
	},
}

// holds a slice of type erase_index
// FIX: duplicates the spec decl.
type FIX_EraseIndex_Slice []EraseIndex

// implements typeinfo.Inspector
func (*EraseIndex_Slice) Inspect() typeinfo.T {
	return &Z_EraseIndex_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Erasing struct {
	Count   rtti.NumberEval
	Target  assign.Address
	AtIndex rtti.NumberEval
	As      string
	Exe     rtti.Execute
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*Erasing) Inspect() typeinfo.T {
	return &Z_Erasing_Info
}

// return a valid markup map, creating it if necessary.
func (op *Erasing) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// erasing, a type of flow.
const Z_Erasing_Type = "erasing"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*Erasing)(nil)

var Z_Erasing_Info = typeinfo.Flow{
	Name: Z_Erasing_Type,
	Lede: "erasing",
	Terms: []typeinfo.Term{{
		Name:  "count",
		Label: "_",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:  "target",
		Label: "from",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "at_index",
		Label: "at_index",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:  "as",
		Label: "as",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:    "exe",
		Label:   "do",
		Repeats: true,
		Type:    &rtti.Z_Execute_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Erase elements from the front or back of a list.", "Runs a pattern with a list containing the erased values.", "If nothing was erased, the pattern will be called with an empty list."},
	},
}

// holds a slice of type erasing
// FIX: duplicates the spec decl.
type FIX_Erasing_Slice []Erasing

// implements typeinfo.Inspector
func (*Erasing_Slice) Inspect() typeinfo.T {
	return &Z_Erasing_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ErasingEdge struct {
	Target assign.Address
	AtEdge rtti.BoolEval
	As     string
	Exe    rtti.Execute
	Else   core.Brancher
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ErasingEdge) Inspect() typeinfo.T {
	return &Z_ErasingEdge_Info
}

// return a valid markup map, creating it if necessary.
func (op *ErasingEdge) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// erasing_edge, a type of flow.
const Z_ErasingEdge_Type = "erasing_edge"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ErasingEdge)(nil)

var Z_ErasingEdge_Info = typeinfo.Flow{
	Name: Z_ErasingEdge_Type,
	Lede: "erasing",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:     "at_edge",
		Label:    "at_front",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}, {
		Name:  "as",
		Label: "as",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:    "exe",
		Label:   "do",
		Repeats: true,
		Type:    &rtti.Z_Execute_Info,
	}, {
		Name:     "else",
		Label:    "else",
		Optional: true,
		Type:     &core.Z_Brancher_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.",
	},
}

// holds a slice of type erasing_edge
// FIX: duplicates the spec decl.
type FIX_ErasingEdge_Slice []ErasingEdge

// implements typeinfo.Inspector
func (*ErasingEdge_Slice) Inspect() typeinfo.T {
	return &Z_ErasingEdge_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListEach struct {
	List   rtti.Assignment
	As     string
	Exe    rtti.Execute
	Else   core.Brancher
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListEach) Inspect() typeinfo.T {
	return &Z_ListEach_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListEach) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_each, a type of flow.
const Z_ListEach_Type = "list_each"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListEach)(nil)

var Z_ListEach_Info = typeinfo.Flow{
	Name: Z_ListEach_Type,
	Lede: "repeating",
	Terms: []typeinfo.Term{{
		Name:  "list",
		Label: "across",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "as",
		Label: "as",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:    "exe",
		Label:   "do",
		Repeats: true,
		Type:    &rtti.Z_Execute_Info,
	}, {
		Name:     "else",
		Label:    "else",
		Optional: true,
		Type:     &core.Z_Brancher_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Loops over the elements in the passed list, or runs the 'else' activity if empty.",
	},
}

// holds a slice of type list_each
// FIX: duplicates the spec decl.
type FIX_ListEach_Slice []ListEach

// implements typeinfo.Inspector
func (*ListEach_Slice) Inspect() typeinfo.T {
	return &Z_ListEach_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListFind struct {
	Value  rtti.Assignment
	List   rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListFind) Inspect() typeinfo.T {
	return &Z_ListFind_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListFind) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_find, a type of flow.
const Z_ListFind_Type = "list_find"

// ensure the command implements its specified slots:
var _ rtti.BoolEval = (*ListFind)(nil)
var _ rtti.NumberEval = (*ListFind)(nil)

var Z_ListFind_Info = typeinfo.Flow{
	Name: Z_ListFind_Type,
	Lede: "find",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "list",
		Label: "in_list",
		Type:  &rtti.Z_Assignment_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_BoolEval_Info,
		&rtti.Z_NumberEval_Info,
	},
	Markup: map[string]any{
		"comment": "Search a list for a specific value.",
	},
}

// holds a slice of type list_find
// FIX: duplicates the spec decl.
type FIX_ListFind_Slice []ListFind

// implements typeinfo.Inspector
func (*ListFind_Slice) Inspect() typeinfo.T {
	return &Z_ListFind_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListGather struct {
	Target assign.Address
	From   rtti.Assignment
	Using  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListGather) Inspect() typeinfo.T {
	return &Z_ListGather_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListGather) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_gather, a type of flow.
const Z_ListGather_Type = "list_gather"

var Z_ListGather_Info = typeinfo.Flow{
	Name: Z_ListGather_Type,
	Lede: "gather",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "from",
		Label: "from",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "using",
		Label: "using",
		Type:  &prim.Z_Text_Info,
	}},
	Markup: map[string]any{
		"comment": []interface{}{"Transform the values from a list.", "The named pattern gets with with two parameters for each value in the list:", "'in' as each value from the list, and 'out' as the var passed to the gather."},
	},
}

// holds a slice of type list_gather
// FIX: duplicates the spec decl.
type FIX_ListGather_Slice []ListGather

// implements typeinfo.Inspector
func (*ListGather_Slice) Inspect() typeinfo.T {
	return &Z_ListGather_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListLen struct {
	List   rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListLen) Inspect() typeinfo.T {
	return &Z_ListLen_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListLen) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_len, a type of flow.
const Z_ListLen_Type = "list_len"

// ensure the command implements its specified slots:
var _ rtti.NumberEval = (*ListLen)(nil)

var Z_ListLen_Info = typeinfo.Flow{
	Name: Z_ListLen_Type,
	Lede: "len",
	Terms: []typeinfo.Term{{
		Name:  "list",
		Label: "_",
		Type:  &rtti.Z_Assignment_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_NumberEval_Info,
	},
	Markup: map[string]any{
		"comment": "Determines the number of values in a list.",
	},
}

// holds a slice of type list_len
// FIX: duplicates the spec decl.
type FIX_ListLen_Slice []ListLen

// implements typeinfo.Inspector
func (*ListLen_Slice) Inspect() typeinfo.T {
	return &Z_ListLen_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_MakeTextList struct {
	Values rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*MakeTextList) Inspect() typeinfo.T {
	return &Z_MakeTextList_Info
}

// return a valid markup map, creating it if necessary.
func (op *MakeTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// make_text_list, a type of flow.
const Z_MakeTextList_Type = "make_text_list"

// ensure the command implements its specified slots:
var _ rtti.TextListEval = (*MakeTextList)(nil)

var Z_MakeTextList_Info = typeinfo.Flow{
	Name: Z_MakeTextList_Type,
	Lede: "list",
	Terms: []typeinfo.Term{{
		Name:    "values",
		Label:   "of_text",
		Repeats: true,
		Type:    &rtti.Z_TextEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_TextListEval_Info,
	},
}

// holds a slice of type make_text_list
// FIX: duplicates the spec decl.
type FIX_MakeTextList_Slice []MakeTextList

// implements typeinfo.Inspector
func (*MakeTextList_Slice) Inspect() typeinfo.T {
	return &Z_MakeTextList_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_MakeNumList struct {
	Values rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*MakeNumList) Inspect() typeinfo.T {
	return &Z_MakeNumList_Info
}

// return a valid markup map, creating it if necessary.
func (op *MakeNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// make_num_list, a type of flow.
const Z_MakeNumList_Type = "make_num_list"

// ensure the command implements its specified slots:
var _ rtti.NumListEval = (*MakeNumList)(nil)

var Z_MakeNumList_Info = typeinfo.Flow{
	Name: Z_MakeNumList_Type,
	Lede: "list",
	Terms: []typeinfo.Term{{
		Name:    "values",
		Label:   "of_numbers",
		Repeats: true,
		Type:    &rtti.Z_NumberEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_NumListEval_Info,
	},
}

// holds a slice of type make_num_list
// FIX: duplicates the spec decl.
type FIX_MakeNumList_Slice []MakeNumList

// implements typeinfo.Inspector
func (*MakeNumList_Slice) Inspect() typeinfo.T {
	return &Z_MakeNumList_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_MakeRecordList struct {
	Values rtti.RecordEval
	Kind   rtti.TextEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*MakeRecordList) Inspect() typeinfo.T {
	return &Z_MakeRecordList_Info
}

// return a valid markup map, creating it if necessary.
func (op *MakeRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// make_record_list, a type of flow.
const Z_MakeRecordList_Type = "make_record_list"

// ensure the command implements its specified slots:
var _ rtti.RecordListEval = (*MakeRecordList)(nil)

var Z_MakeRecordList_Info = typeinfo.Flow{
	Name: Z_MakeRecordList_Type,
	Lede: "list",
	Terms: []typeinfo.Term{{
		Name:    "values",
		Label:   "of_records",
		Repeats: true,
		Type:    &rtti.Z_RecordEval_Info,
	}, {
		Name:  "kind",
		Label: "of_type",
		Type:  &rtti.Z_TextEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_RecordListEval_Info,
	},
}

// holds a slice of type make_record_list
// FIX: duplicates the spec decl.
type FIX_MakeRecordList_Slice []MakeRecordList

// implements typeinfo.Inspector
func (*MakeRecordList_Slice) Inspect() typeinfo.T {
	return &Z_MakeRecordList_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListMap struct {
	Target      assign.Address
	List        rtti.Assignment
	PatternName string
	Markup      map[string]any
}

// implements typeinfo.Inspector
func (*ListMap) Inspect() typeinfo.T {
	return &Z_ListMap_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListMap) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_map, a type of flow.
const Z_ListMap_Type = "list_map"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListMap)(nil)

var Z_ListMap_Info = typeinfo.Flow{
	Name: Z_ListMap_Type,
	Lede: "map",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "list",
		Label: "from_list",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "pattern_name",
		Label: "using",
		Type:  &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Transform the values from one list and place the results in another list.", "The designated pattern is called with each value from the 'from list', one value at a time."},
	},
}

// holds a slice of type list_map
// FIX: duplicates the spec decl.
type FIX_ListMap_Slice []ListMap

// implements typeinfo.Inspector
func (*ListMap_Slice) Inspect() typeinfo.T {
	return &Z_ListMap_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListReduce struct {
	Target      assign.Address
	List        rtti.Assignment
	PatternName string
	Markup      map[string]any
}

// implements typeinfo.Inspector
func (*ListReduce) Inspect() typeinfo.T {
	return &Z_ListReduce_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListReduce) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_reduce, a type of flow.
const Z_ListReduce_Type = "list_reduce"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListReduce)(nil)

var Z_ListReduce_Info = typeinfo.Flow{
	Name: Z_ListReduce_Type,
	Lede: "reduce",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "into",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "list",
		Label: "from_list",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "pattern_name",
		Label: "using",
		Type:  &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Combine all of the values in a list into a single value.", "The designated pattern is called with two parameters:", "  1. each element of the list; and,", "  2. the value being combined.", "And, that pattern is expected to return the newly updated value."},
	},
}

// holds a slice of type list_reduce
// FIX: duplicates the spec decl.
type FIX_ListReduce_Slice []ListReduce

// implements typeinfo.Inspector
func (*ListReduce_Slice) Inspect() typeinfo.T {
	return &Z_ListReduce_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListReverse struct {
	Target assign.Address
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListReverse) Inspect() typeinfo.T {
	return &Z_ListReverse_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListReverse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_reverse, a type of flow.
const Z_ListReverse_Type = "list_reverse"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListReverse)(nil)

var Z_ListReverse_Info = typeinfo.Flow{
	Name: Z_ListReverse_Type,
	Lede: "reverse",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "list",
		Type:  &assign.Z_Address_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Reverse a list.",
	},
}

// holds a slice of type list_reverse
// FIX: duplicates the spec decl.
type FIX_ListReverse_Slice []ListReverse

// implements typeinfo.Inspector
func (*ListReverse_Slice) Inspect() typeinfo.T {
	return &Z_ListReverse_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListSlice struct {
	List   rtti.Assignment
	Start  rtti.NumberEval
	End    rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListSlice) Inspect() typeinfo.T {
	return &Z_ListSlice_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListSlice) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_slice, a type of flow.
const Z_ListSlice_Type = "list_slice"

// ensure the command implements its specified slots:
var _ rtti.NumListEval = (*ListSlice)(nil)
var _ rtti.TextListEval = (*ListSlice)(nil)
var _ rtti.RecordListEval = (*ListSlice)(nil)

var Z_ListSlice_Info = typeinfo.Flow{
	Name: Z_ListSlice_Type,
	Lede: "slice",
	Terms: []typeinfo.Term{{
		Name:  "list",
		Label: "_",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:     "start",
		Label:    "start",
		Optional: true,
		Type:     &rtti.Z_NumberEval_Info,
	}, {
		Name:     "end",
		Label:    "end",
		Optional: true,
		Type:     &rtti.Z_NumberEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_NumListEval_Info,
		&rtti.Z_TextListEval_Info,
		&rtti.Z_RecordListEval_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Create a new list from a section of another list.", "Start is optional, if omitted slice starts at the first element.", "If start is greater the length, an empty array is returned.", "Slice doesnt include the ending index.", "Negatives indices indicates an offset from the end.", "When end is omitted, copy up to and including the last element;", "and do the same if the end is greater than the length"},
	},
}

// holds a slice of type list_slice
// FIX: duplicates the spec decl.
type FIX_ListSlice_Slice []ListSlice

// implements typeinfo.Inspector
func (*ListSlice_Slice) Inspect() typeinfo.T {
	return &Z_ListSlice_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListSortNumbers struct {
	Target     assign.Address
	ByField    string
	Descending rtti.BoolEval
	Markup     map[string]any
}

// implements typeinfo.Inspector
func (*ListSortNumbers) Inspect() typeinfo.T {
	return &Z_ListSortNumbers_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListSortNumbers) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_sort_numbers, a type of flow.
const Z_ListSortNumbers_Type = "list_sort_numbers"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListSortNumbers)(nil)

var Z_ListSortNumbers_Info = typeinfo.Flow{
	Name: Z_ListSortNumbers_Type,
	Lede: "sort_numbers",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "by_field",
		Label: "by_field",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:     "descending",
		Label:    "descending",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
}

// holds a slice of type list_sort_numbers
// FIX: duplicates the spec decl.
type FIX_ListSortNumbers_Slice []ListSortNumbers

// implements typeinfo.Inspector
func (*ListSortNumbers_Slice) Inspect() typeinfo.T {
	return &Z_ListSortNumbers_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListSortText struct {
	Target     assign.Address
	ByField    string
	Descending rtti.BoolEval
	UsingCase  rtti.BoolEval
	Markup     map[string]any
}

// implements typeinfo.Inspector
func (*ListSortText) Inspect() typeinfo.T {
	return &Z_ListSortText_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListSortText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_sort_text, a type of flow.
const Z_ListSortText_Type = "list_sort_text"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListSortText)(nil)

var Z_ListSortText_Info = typeinfo.Flow{
	Name: Z_ListSortText_Type,
	Lede: "sort_texts",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "by_field",
		Label: "by_field",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:     "descending",
		Label:    "descending",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}, {
		Name:     "using_case",
		Label:    "using_case",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
	},
}

// holds a slice of type list_sort_text
// FIX: duplicates the spec decl.
type FIX_ListSortText_Slice []ListSortText

// implements typeinfo.Inspector
func (*ListSortText_Slice) Inspect() typeinfo.T {
	return &Z_ListSortText_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListSplice struct {
	Target assign.Address
	Start  rtti.NumberEval
	Remove rtti.NumberEval
	Insert rtti.Assignment
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListSplice) Inspect() typeinfo.T {
	return &Z_ListSplice_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListSplice) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_splice, a type of flow.
const Z_ListSplice_Type = "list_splice"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListSplice)(nil)
var _ rtti.NumListEval = (*ListSplice)(nil)
var _ rtti.TextListEval = (*ListSplice)(nil)
var _ rtti.RecordListEval = (*ListSplice)(nil)

var Z_ListSplice_Info = typeinfo.Flow{
	Name: Z_ListSplice_Type,
	Lede: "splice",
	Terms: []typeinfo.Term{{
		Name:  "target",
		Label: "_",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:  "start",
		Label: "start",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:  "remove",
		Label: "remove",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:  "insert",
		Label: "insert",
		Type:  &rtti.Z_Assignment_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
		&rtti.Z_NumListEval_Info,
		&rtti.Z_TextListEval_Info,
		&rtti.Z_RecordListEval_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Modify a list by adding and removing elements.", "The type of the elements being added must match the type of the list.", "Text cant be added to a list of numbers, numbers cant be added to a list of text.", "If the starting index is negative, this begins that many elements from the end of the array;", "if list's length plus the start is less than zero, this begins from index zero.", "If the remove count is missing, this removes all elements from the start to the end;", "if the remove count is zero or negative, no elements are removed."},
	},
}

// holds a slice of type list_splice
// FIX: duplicates the spec decl.
type FIX_ListSplice_Slice []ListSplice

// implements typeinfo.Inspector
func (*ListSplice_Slice) Inspect() typeinfo.T {
	return &Z_ListSplice_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ListPush struct {
	Value  rtti.Assignment
	Target assign.Address
	AtEdge rtti.BoolEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*ListPush) Inspect() typeinfo.T {
	return &Z_ListPush_Info
}

// return a valid markup map, creating it if necessary.
func (op *ListPush) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// list_push, a type of flow.
const Z_ListPush_Type = "list_push"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*ListPush)(nil)

var Z_ListPush_Info = typeinfo.Flow{
	Name: Z_ListPush_Type,
	Lede: "push",
	Terms: []typeinfo.Term{{
		Name:  "value",
		Label: "_",
		Type:  &rtti.Z_Assignment_Info,
	}, {
		Name:  "target",
		Label: "into",
		Type:  &assign.Z_Address_Info,
	}, {
		Name:     "at_edge",
		Label:    "at_front",
		Optional: true,
		Type:     &rtti.Z_BoolEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_Info,
	},
	Markup: map[string]any{
		"comment": "Add a value to a list.",
	},
}

// holds a slice of type list_push
// FIX: duplicates the spec decl.
type FIX_ListPush_Slice []ListPush

// implements typeinfo.Inspector
func (*ListPush_Slice) Inspect() typeinfo.T {
	return &Z_ListPush_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Range struct {
	To     rtti.NumberEval
	From   rtti.NumberEval
	ByStep rtti.NumberEval
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Range) Inspect() typeinfo.T {
	return &Z_Range_Info
}

// return a valid markup map, creating it if necessary.
func (op *Range) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// range, a type of flow.
const Z_Range_Type = "range"

// ensure the command implements its specified slots:
var _ rtti.NumListEval = (*Range)(nil)

var Z_Range_Info = typeinfo.Flow{
	Name: Z_Range_Type,
	Lede: "range",
	Terms: []typeinfo.Term{{
		Name:  "to",
		Label: "_",
		Type:  &rtti.Z_NumberEval_Info,
	}, {
		Name:     "from",
		Label:    "from",
		Optional: true,
		Type:     &rtti.Z_NumberEval_Info,
	}, {
		Name:     "by_step",
		Label:    "by_step",
		Optional: true,
		Type:     &rtti.Z_NumberEval_Info,
	}},
	Slots: []*typeinfo.Slot{
		&rtti.Z_NumListEval_Info,
	},
	Markup: map[string]any{
		"comment": []interface{}{"Generates a series of numbers r[i] = (start + step*i) where i>=0.", "Start and step default to 1, stop defaults to start;", "the inputs are truncated to produce whole numbers;", "a zero step returns an error.", "A positive step ends the series when the returned value would exceed stop", "while a negative step ends before generating a value less than stop."},
	},
}

// holds a slice of type range
// FIX: duplicates the spec decl.
type FIX_Range_Slice []Range

// implements typeinfo.Inspector
func (*Range_Slice) Inspect() typeinfo.T {
	return &Z_Range_Info
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var Y_flow_List = []*typeinfo.Flow{
	&Z_EraseEdge_Info,
	&Z_EraseIndex_Info,
	&Z_Erasing_Info,
	&Z_ErasingEdge_Info,
	&Z_ListEach_Info,
	&Z_ListFind_Info,
	&Z_ListGather_Info,
	&Z_ListLen_Info,
	&Z_MakeTextList_Info,
	&Z_MakeNumList_Info,
	&Z_MakeRecordList_Info,
	&Z_ListMap_Info,
	&Z_ListReduce_Info,
	&Z_ListReverse_Info,
	&Z_ListSlice_Info,
	&Z_ListSortNumbers_Info,
	&Z_ListSortText_Info,
	&Z_ListSplice_Info,
	&Z_ListPush_Info,
	&Z_Range_Info,
}

// a list of all command signatures
// ( for processing and verifying story files )
var Z_Signatures = map[uint64]interface{}{}
