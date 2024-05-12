// list
package list

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// Erase at edge: Remove one or more values from a list.
type EraseEdge struct {
	Target assign.Address
	AtEdge rtti.BoolEval
	Markup map[string]any
}

// erase_edge, a type of flow.
var Zt_EraseEdge typeinfo.Flow

// Implements [typeinfo.Instance]
func (*EraseEdge) TypeInfo() typeinfo.T {
	return &Zt_EraseEdge
}

// Implements [typeinfo.Markup]
func (op *EraseEdge) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*EraseEdge)(nil)

// Holds a slice of type EraseEdge.
type EraseEdge_Slice []EraseEdge

// Implements [typeinfo.Instance] for a slice of EraseEdge.
func (*EraseEdge_Slice) TypeInfo() typeinfo.T {
	return &Zt_EraseEdge
}

// Implements [typeinfo.Repeats] for a slice of EraseEdge.
func (op *EraseEdge_Slice) Repeats() bool {
	return len(*op) > 0
}

// Erase at index: Remove one or more values from a list.
type EraseIndex struct {
	Count   rtti.NumberEval
	Target  assign.Address
	AtIndex rtti.NumberEval
	Markup  map[string]any
}

// erase_index, a type of flow.
var Zt_EraseIndex typeinfo.Flow

// Implements [typeinfo.Instance]
func (*EraseIndex) TypeInfo() typeinfo.T {
	return &Zt_EraseIndex
}

// Implements [typeinfo.Markup]
func (op *EraseIndex) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*EraseIndex)(nil)

// Holds a slice of type EraseIndex.
type EraseIndex_Slice []EraseIndex

// Implements [typeinfo.Instance] for a slice of EraseIndex.
func (*EraseIndex_Slice) TypeInfo() typeinfo.T {
	return &Zt_EraseIndex
}

// Implements [typeinfo.Repeats] for a slice of EraseIndex.
func (op *EraseIndex_Slice) Repeats() bool {
	return len(*op) > 0
}

// Erase elements from the front or back of a list.
// Runs a pattern with a list containing the erased values.
// If nothing was erased, the pattern will be called with an empty list.
type Erasing struct {
	Count   rtti.NumberEval
	Target  assign.Address
	AtIndex rtti.NumberEval
	As      string
	Exe     []rtti.Execute
	Markup  map[string]any
}

// erasing, a type of flow.
var Zt_Erasing typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Erasing) TypeInfo() typeinfo.T {
	return &Zt_Erasing
}

// Implements [typeinfo.Markup]
func (op *Erasing) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*Erasing)(nil)

// Holds a slice of type Erasing.
type Erasing_Slice []Erasing

// Implements [typeinfo.Instance] for a slice of Erasing.
func (*Erasing_Slice) TypeInfo() typeinfo.T {
	return &Zt_Erasing
}

// Implements [typeinfo.Repeats] for a slice of Erasing.
func (op *Erasing_Slice) Repeats() bool {
	return len(*op) > 0
}

// Erase one element from the front or back of a list.
// Runs an activity with a list containing the erased values;
// the list can be empty if nothing was erased.
type ErasingEdge struct {
	Target assign.Address
	AtEdge rtti.BoolEval
	As     string
	Exe    []rtti.Execute
	Else   core.Brancher
	Markup map[string]any
}

// erasing_edge, a type of flow.
var Zt_ErasingEdge typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ErasingEdge) TypeInfo() typeinfo.T {
	return &Zt_ErasingEdge
}

// Implements [typeinfo.Markup]
func (op *ErasingEdge) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ErasingEdge)(nil)

// Holds a slice of type ErasingEdge.
type ErasingEdge_Slice []ErasingEdge

// Implements [typeinfo.Instance] for a slice of ErasingEdge.
func (*ErasingEdge_Slice) TypeInfo() typeinfo.T {
	return &Zt_ErasingEdge
}

// Implements [typeinfo.Repeats] for a slice of ErasingEdge.
func (op *ErasingEdge_Slice) Repeats() bool {
	return len(*op) > 0
}

// Loops over the elements in the passed list, or runs the 'else' activity if empty.
type ListEach struct {
	List   rtti.Assignment
	As     string
	Exe    []rtti.Execute
	Else   core.Brancher
	Markup map[string]any
}

// list_each, a type of flow.
var Zt_ListEach typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListEach) TypeInfo() typeinfo.T {
	return &Zt_ListEach
}

// Implements [typeinfo.Markup]
func (op *ListEach) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListEach)(nil)

// Holds a slice of type ListEach.
type ListEach_Slice []ListEach

// Implements [typeinfo.Instance] for a slice of ListEach.
func (*ListEach_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListEach
}

// Implements [typeinfo.Repeats] for a slice of ListEach.
func (op *ListEach_Slice) Repeats() bool {
	return len(*op) > 0
}

// Search a list for a specific value.
type ListFind struct {
	Value  rtti.Assignment
	List   rtti.Assignment
	Markup map[string]any
}

// list_find, a type of flow.
var Zt_ListFind typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListFind) TypeInfo() typeinfo.T {
	return &Zt_ListFind
}

// Implements [typeinfo.Markup]
func (op *ListFind) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*ListFind)(nil)
var _ rtti.NumberEval = (*ListFind)(nil)

// Holds a slice of type ListFind.
type ListFind_Slice []ListFind

// Implements [typeinfo.Instance] for a slice of ListFind.
func (*ListFind_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListFind
}

// Implements [typeinfo.Repeats] for a slice of ListFind.
func (op *ListFind_Slice) Repeats() bool {
	return len(*op) > 0
}

// Transform the values from a list.
// The named pattern gets with with two parameters for each value in the list:
// 'in' as each value from the list, and 'out' as the var passed to the gather.
type ListGather struct {
	Target assign.Address
	From   rtti.Assignment
	Using  string
	Markup map[string]any
}

// list_gather, a type of flow.
var Zt_ListGather typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListGather) TypeInfo() typeinfo.T {
	return &Zt_ListGather
}

// Implements [typeinfo.Markup]
func (op *ListGather) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Holds a slice of type ListGather.
type ListGather_Slice []ListGather

// Implements [typeinfo.Instance] for a slice of ListGather.
func (*ListGather_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListGather
}

// Implements [typeinfo.Repeats] for a slice of ListGather.
func (op *ListGather_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determines the number of values in a list.
type ListLen struct {
	List   rtti.Assignment
	Markup map[string]any
}

// list_len, a type of flow.
var Zt_ListLen typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListLen) TypeInfo() typeinfo.T {
	return &Zt_ListLen
}

// Implements [typeinfo.Markup]
func (op *ListLen) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumberEval = (*ListLen)(nil)

// Holds a slice of type ListLen.
type ListLen_Slice []ListLen

// Implements [typeinfo.Instance] for a slice of ListLen.
func (*ListLen_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListLen
}

// Implements [typeinfo.Repeats] for a slice of ListLen.
func (op *ListLen_Slice) Repeats() bool {
	return len(*op) > 0
}

type MakeTextList struct {
	Values []rtti.TextEval
	Markup map[string]any
}

// make_text_list, a type of flow.
var Zt_MakeTextList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeTextList) TypeInfo() typeinfo.T {
	return &Zt_MakeTextList
}

// Implements [typeinfo.Markup]
func (op *MakeTextList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextListEval = (*MakeTextList)(nil)

// Holds a slice of type MakeTextList.
type MakeTextList_Slice []MakeTextList

// Implements [typeinfo.Instance] for a slice of MakeTextList.
func (*MakeTextList_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeTextList
}

// Implements [typeinfo.Repeats] for a slice of MakeTextList.
func (op *MakeTextList_Slice) Repeats() bool {
	return len(*op) > 0
}

type MakeNumList struct {
	Values []rtti.NumberEval
	Markup map[string]any
}

// make_num_list, a type of flow.
var Zt_MakeNumList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeNumList) TypeInfo() typeinfo.T {
	return &Zt_MakeNumList
}

// Implements [typeinfo.Markup]
func (op *MakeNumList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumListEval = (*MakeNumList)(nil)

// Holds a slice of type MakeNumList.
type MakeNumList_Slice []MakeNumList

// Implements [typeinfo.Instance] for a slice of MakeNumList.
func (*MakeNumList_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeNumList
}

// Implements [typeinfo.Repeats] for a slice of MakeNumList.
func (op *MakeNumList_Slice) Repeats() bool {
	return len(*op) > 0
}

type MakeRecordList struct {
	Values []rtti.RecordEval
	Kind   rtti.TextEval
	Markup map[string]any
}

// make_record_list, a type of flow.
var Zt_MakeRecordList typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeRecordList) TypeInfo() typeinfo.T {
	return &Zt_MakeRecordList
}

// Implements [typeinfo.Markup]
func (op *MakeRecordList) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.RecordListEval = (*MakeRecordList)(nil)

// Holds a slice of type MakeRecordList.
type MakeRecordList_Slice []MakeRecordList

// Implements [typeinfo.Instance] for a slice of MakeRecordList.
func (*MakeRecordList_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeRecordList
}

// Implements [typeinfo.Repeats] for a slice of MakeRecordList.
func (op *MakeRecordList_Slice) Repeats() bool {
	return len(*op) > 0
}

// Transform the values from one list and place the results in another list.
// The designated pattern is called with each value from the 'from list', one value at a time.
type ListMap struct {
	Target      assign.Address
	List        rtti.Assignment
	PatternName string
	Markup      map[string]any
}

// list_map, a type of flow.
var Zt_ListMap typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListMap) TypeInfo() typeinfo.T {
	return &Zt_ListMap
}

// Implements [typeinfo.Markup]
func (op *ListMap) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListMap)(nil)

// Holds a slice of type ListMap.
type ListMap_Slice []ListMap

// Implements [typeinfo.Instance] for a slice of ListMap.
func (*ListMap_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListMap
}

// Implements [typeinfo.Repeats] for a slice of ListMap.
func (op *ListMap_Slice) Repeats() bool {
	return len(*op) > 0
}

// Combine all of the values in a list into a single value.
// The designated pattern is called with two parameters:
//  1. each element of the list; and,
//  2. the value being combined.
//
// And, that pattern is expected to return the newly updated value.
type ListReduce struct {
	Target      assign.Address
	List        rtti.Assignment
	PatternName string
	Markup      map[string]any
}

// list_reduce, a type of flow.
var Zt_ListReduce typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListReduce) TypeInfo() typeinfo.T {
	return &Zt_ListReduce
}

// Implements [typeinfo.Markup]
func (op *ListReduce) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListReduce)(nil)

// Holds a slice of type ListReduce.
type ListReduce_Slice []ListReduce

// Implements [typeinfo.Instance] for a slice of ListReduce.
func (*ListReduce_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListReduce
}

// Implements [typeinfo.Repeats] for a slice of ListReduce.
func (op *ListReduce_Slice) Repeats() bool {
	return len(*op) > 0
}

// Reverse a list.
type ListReverse struct {
	Target assign.Address
	Markup map[string]any
}

// list_reverse, a type of flow.
var Zt_ListReverse typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListReverse) TypeInfo() typeinfo.T {
	return &Zt_ListReverse
}

// Implements [typeinfo.Markup]
func (op *ListReverse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListReverse)(nil)

// Holds a slice of type ListReverse.
type ListReverse_Slice []ListReverse

// Implements [typeinfo.Instance] for a slice of ListReverse.
func (*ListReverse_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListReverse
}

// Implements [typeinfo.Repeats] for a slice of ListReverse.
func (op *ListReverse_Slice) Repeats() bool {
	return len(*op) > 0
}

// Create a new list from a section of another list.
// Start is optional, if omitted slice starts at the first element.
// If start is greater the length, an empty array is returned.
// Slice doesnt include the ending index.
// Negatives indices indicates an offset from the end.
// When end is omitted, copy up to and including the last element;
// and do the same if the end is greater than the length
type ListSlice struct {
	List   rtti.Assignment
	Start  rtti.NumberEval
	End    rtti.NumberEval
	Markup map[string]any
}

// list_slice, a type of flow.
var Zt_ListSlice typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListSlice) TypeInfo() typeinfo.T {
	return &Zt_ListSlice
}

// Implements [typeinfo.Markup]
func (op *ListSlice) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumListEval = (*ListSlice)(nil)
var _ rtti.TextListEval = (*ListSlice)(nil)
var _ rtti.RecordListEval = (*ListSlice)(nil)

// Holds a slice of type ListSlice.
type ListSlice_Slice []ListSlice

// Implements [typeinfo.Instance] for a slice of ListSlice.
func (*ListSlice_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListSlice
}

// Implements [typeinfo.Repeats] for a slice of ListSlice.
func (op *ListSlice_Slice) Repeats() bool {
	return len(*op) > 0
}

type ListSortNumbers struct {
	Target     assign.Address
	ByField    string
	Descending rtti.BoolEval
	Markup     map[string]any
}

// list_sort_numbers, a type of flow.
var Zt_ListSortNumbers typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListSortNumbers) TypeInfo() typeinfo.T {
	return &Zt_ListSortNumbers
}

// Implements [typeinfo.Markup]
func (op *ListSortNumbers) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListSortNumbers)(nil)

// Holds a slice of type ListSortNumbers.
type ListSortNumbers_Slice []ListSortNumbers

// Implements [typeinfo.Instance] for a slice of ListSortNumbers.
func (*ListSortNumbers_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListSortNumbers
}

// Implements [typeinfo.Repeats] for a slice of ListSortNumbers.
func (op *ListSortNumbers_Slice) Repeats() bool {
	return len(*op) > 0
}

// Rearrange the elements in the named list by using the designated pattern to test pairs of elements.
type ListSortText struct {
	Target     assign.Address
	ByField    string
	Descending rtti.BoolEval
	UsingCase  rtti.BoolEval
	Markup     map[string]any
}

// list_sort_text, a type of flow.
var Zt_ListSortText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListSortText) TypeInfo() typeinfo.T {
	return &Zt_ListSortText
}

// Implements [typeinfo.Markup]
func (op *ListSortText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListSortText)(nil)

// Holds a slice of type ListSortText.
type ListSortText_Slice []ListSortText

// Implements [typeinfo.Instance] for a slice of ListSortText.
func (*ListSortText_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListSortText
}

// Implements [typeinfo.Repeats] for a slice of ListSortText.
func (op *ListSortText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Modify a list by adding and removing elements.
// The type of the elements being added must match the type of the list.
// Text cant be added to a list of numbers, numbers cant be added to a list of text.
// If the starting index is negative, this begins that many elements from the end of the array;
// if list's length plus the start is less than zero, this begins from index zero.
// If the remove count is missing, this removes all elements from the start to the end;
// if the remove count is zero or negative, no elements are removed.
type ListSplice struct {
	Target assign.Address
	Start  rtti.NumberEval
	Remove rtti.NumberEval
	Insert rtti.Assignment
	Markup map[string]any
}

// list_splice, a type of flow.
var Zt_ListSplice typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListSplice) TypeInfo() typeinfo.T {
	return &Zt_ListSplice
}

// Implements [typeinfo.Markup]
func (op *ListSplice) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListSplice)(nil)
var _ rtti.NumListEval = (*ListSplice)(nil)
var _ rtti.TextListEval = (*ListSplice)(nil)
var _ rtti.RecordListEval = (*ListSplice)(nil)

// Holds a slice of type ListSplice.
type ListSplice_Slice []ListSplice

// Implements [typeinfo.Instance] for a slice of ListSplice.
func (*ListSplice_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListSplice
}

// Implements [typeinfo.Repeats] for a slice of ListSplice.
func (op *ListSplice_Slice) Repeats() bool {
	return len(*op) > 0
}

// Add a value to a list.
type ListPush struct {
	Value  rtti.Assignment
	Target assign.Address
	AtEdge rtti.BoolEval
	Markup map[string]any
}

// list_push, a type of flow.
var Zt_ListPush typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ListPush) TypeInfo() typeinfo.T {
	return &Zt_ListPush
}

// Implements [typeinfo.Markup]
func (op *ListPush) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ListPush)(nil)

// Holds a slice of type ListPush.
type ListPush_Slice []ListPush

// Implements [typeinfo.Instance] for a slice of ListPush.
func (*ListPush_Slice) TypeInfo() typeinfo.T {
	return &Zt_ListPush
}

// Implements [typeinfo.Repeats] for a slice of ListPush.
func (op *ListPush_Slice) Repeats() bool {
	return len(*op) > 0
}

// Generates a series of numbers r[i] = (start + step*i) where i>=0.
// Start and step default to 1, stop defaults to start;
// the inputs are truncated to produce whole numbers;
// a zero step returns an error.
// A positive step ends the series when the returned value would exceed stop
// while a negative step ends before generating a value less than stop.
type Range struct {
	To     rtti.NumberEval
	From   rtti.NumberEval
	ByStep rtti.NumberEval
	Markup map[string]any
}

// range, a type of flow.
var Zt_Range typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Range) TypeInfo() typeinfo.T {
	return &Zt_Range
}

// Implements [typeinfo.Markup]
func (op *Range) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumListEval = (*Range)(nil)

// Holds a slice of type Range.
type Range_Slice []Range

// Implements [typeinfo.Instance] for a slice of Range.
func (*Range_Slice) TypeInfo() typeinfo.T {
	return &Zt_Range
}

// Implements [typeinfo.Repeats] for a slice of Range.
func (op *Range_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_EraseEdge = typeinfo.Flow{
		Name: "erase_edge",
		Lede: "erase",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:     "at_edge",
			Label:    "at_front",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Erase at edge: Remove one or more values from a list.",
		},
	}
	Zt_EraseIndex = typeinfo.Flow{
		Name: "erase_index",
		Lede: "erase",
		Terms: []typeinfo.Term{{
			Name: "count",
			Type: &rtti.Zt_NumberEval,
		}, {
			Name:  "target",
			Label: "from",
			Type:  &assign.Zt_Address,
		}, {
			Name:  "at_index",
			Label: "at_index",
			Type:  &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Erase at index: Remove one or more values from a list.",
		},
	}
	Zt_Erasing = typeinfo.Flow{
		Name: "erasing",
		Lede: "erasing",
		Terms: []typeinfo.Term{{
			Name: "count",
			Type: &rtti.Zt_NumberEval,
		}, {
			Name:  "target",
			Label: "from",
			Type:  &assign.Zt_Address,
		}, {
			Name:  "at_index",
			Label: "at_index",
			Type:  &rtti.Zt_NumberEval,
		}, {
			Name:  "as",
			Label: "as",
			Type:  &prim.Zt_Text,
		}, {
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Type:    &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Erase elements from the front or back of a list.", "Runs a pattern with a list containing the erased values.", "If nothing was erased, the pattern will be called with an empty list."},
		},
	}
	Zt_ErasingEdge = typeinfo.Flow{
		Name: "erasing_edge",
		Lede: "erasing",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:     "at_edge",
			Label:    "at_front",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}, {
			Name:  "as",
			Label: "as",
			Type:  &prim.Zt_Text,
		}, {
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Type:    &rtti.Zt_Execute,
		}, {
			Name:     "else",
			Label:    "else",
			Optional: true,
			Type:     &core.Zt_Brancher,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Erase one element from the front or back of a list.", "Runs an activity with a list containing the erased values;", "the list can be empty if nothing was erased."},
		},
	}
	Zt_ListEach = typeinfo.Flow{
		Name: "list_each",
		Lede: "repeating",
		Terms: []typeinfo.Term{{
			Name:  "list",
			Label: "across",
			Type:  &rtti.Zt_Assignment,
		}, {
			Name:  "as",
			Label: "as",
			Type:  &prim.Zt_Text,
		}, {
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Type:    &rtti.Zt_Execute,
		}, {
			Name:     "else",
			Label:    "else",
			Optional: true,
			Type:     &core.Zt_Brancher,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Loops over the elements in the passed list, or runs the 'else' activity if empty.",
		},
	}
	Zt_ListFind = typeinfo.Flow{
		Name: "list_find",
		Lede: "find",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_Assignment,
		}, {
			Name:  "list",
			Label: "in_list",
			Type:  &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumberEval,
		},
		Markup: map[string]any{
			"comment": "Search a list for a specific value.",
		},
	}
	Zt_ListGather = typeinfo.Flow{
		Name: "list_gather",
		Lede: "gather",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:  "from",
			Label: "from",
			Type:  &rtti.Zt_Assignment,
		}, {
			Name:  "using",
			Label: "using",
			Type:  &prim.Zt_Text,
		}},
		Markup: map[string]any{
			"comment": []interface{}{"Transform the values from a list.", "The named pattern gets with with two parameters for each value in the list:", "'in' as each value from the list, and 'out' as the var passed to the gather."},
		},
	}
	Zt_ListLen = typeinfo.Flow{
		Name: "list_len",
		Lede: "len",
		Terms: []typeinfo.Term{{
			Name: "list",
			Type: &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumberEval,
		},
		Markup: map[string]any{
			"comment": "Determines the number of values in a list.",
		},
	}
	Zt_MakeTextList = typeinfo.Flow{
		Name: "make_text_list",
		Lede: "list",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "of_text",
			Repeats: true,
			Type:    &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextListEval,
		},
	}
	Zt_MakeNumList = typeinfo.Flow{
		Name: "make_num_list",
		Lede: "list",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "of_numbers",
			Repeats: true,
			Type:    &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumListEval,
		},
	}
	Zt_MakeRecordList = typeinfo.Flow{
		Name: "make_record_list",
		Lede: "list",
		Terms: []typeinfo.Term{{
			Name:    "values",
			Label:   "of_records",
			Repeats: true,
			Type:    &rtti.Zt_RecordEval,
		}, {
			Name:  "kind",
			Label: "of_type",
			Type:  &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_RecordListEval,
		},
	}
	Zt_ListMap = typeinfo.Flow{
		Name: "list_map",
		Lede: "map",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:  "list",
			Label: "from_list",
			Type:  &rtti.Zt_Assignment,
		}, {
			Name:  "pattern_name",
			Label: "using",
			Type:  &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Transform the values from one list and place the results in another list.", "The designated pattern is called with each value from the 'from list', one value at a time."},
		},
	}
	Zt_ListReduce = typeinfo.Flow{
		Name: "list_reduce",
		Lede: "reduce",
		Terms: []typeinfo.Term{{
			Name:  "target",
			Label: "into",
			Type:  &assign.Zt_Address,
		}, {
			Name:  "list",
			Label: "from_list",
			Type:  &rtti.Zt_Assignment,
		}, {
			Name:  "pattern_name",
			Label: "using",
			Type:  &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Combine all of the values in a list into a single value.", "The designated pattern is called with two parameters:", "  1. each element of the list; and,", "  2. the value being combined.", "And, that pattern is expected to return the newly updated value."},
		},
	}
	Zt_ListReverse = typeinfo.Flow{
		Name: "list_reverse",
		Lede: "reverse",
		Terms: []typeinfo.Term{{
			Name:  "target",
			Label: "list",
			Type:  &assign.Zt_Address,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Reverse a list.",
		},
	}
	Zt_ListSlice = typeinfo.Flow{
		Name: "list_slice",
		Lede: "slice",
		Terms: []typeinfo.Term{{
			Name: "list",
			Type: &rtti.Zt_Assignment,
		}, {
			Name:     "start",
			Label:    "start",
			Optional: true,
			Type:     &rtti.Zt_NumberEval,
		}, {
			Name:     "end",
			Label:    "end",
			Optional: true,
			Type:     &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumListEval,
			&rtti.Zt_TextListEval,
			&rtti.Zt_RecordListEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Create a new list from a section of another list.", "Start is optional, if omitted slice starts at the first element.", "If start is greater the length, an empty array is returned.", "Slice doesnt include the ending index.", "Negatives indices indicates an offset from the end.", "When end is omitted, copy up to and including the last element;", "and do the same if the end is greater than the length"},
		},
	}
	Zt_ListSortNumbers = typeinfo.Flow{
		Name: "list_sort_numbers",
		Lede: "sort_numbers",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:  "by_field",
			Label: "by_field",
			Type:  &prim.Zt_Text,
		}, {
			Name:     "descending",
			Label:    "descending",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
	Zt_ListSortText = typeinfo.Flow{
		Name: "list_sort_text",
		Lede: "sort_texts",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:  "by_field",
			Label: "by_field",
			Type:  &prim.Zt_Text,
		}, {
			Name:     "descending",
			Label:    "descending",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}, {
			Name:     "using_case",
			Label:    "using_case",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
		},
	}
	Zt_ListSplice = typeinfo.Flow{
		Name: "list_splice",
		Lede: "splice",
		Terms: []typeinfo.Term{{
			Name: "target",
			Type: &assign.Zt_Address,
		}, {
			Name:  "start",
			Label: "start",
			Type:  &rtti.Zt_NumberEval,
		}, {
			Name:  "remove",
			Label: "remove",
			Type:  &rtti.Zt_NumberEval,
		}, {
			Name:  "insert",
			Label: "insert",
			Type:  &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
			&rtti.Zt_NumListEval,
			&rtti.Zt_TextListEval,
			&rtti.Zt_RecordListEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Modify a list by adding and removing elements.", "The type of the elements being added must match the type of the list.", "Text cant be added to a list of numbers, numbers cant be added to a list of text.", "If the starting index is negative, this begins that many elements from the end of the array;", "if list's length plus the start is less than zero, this begins from index zero.", "If the remove count is missing, this removes all elements from the start to the end;", "if the remove count is zero or negative, no elements are removed."},
		},
	}
	Zt_ListPush = typeinfo.Flow{
		Name: "list_push",
		Lede: "push",
		Terms: []typeinfo.Term{{
			Name: "value",
			Type: &rtti.Zt_Assignment,
		}, {
			Name:  "target",
			Label: "into",
			Type:  &assign.Zt_Address,
		}, {
			Name:     "at_edge",
			Label:    "at_front",
			Optional: true,
			Type:     &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Add a value to a list.",
		},
	}
	Zt_Range = typeinfo.Flow{
		Name: "range",
		Lede: "range",
		Terms: []typeinfo.Term{{
			Name: "to",
			Type: &rtti.Zt_NumberEval,
		}, {
			Name:     "from",
			Label:    "from",
			Optional: true,
			Type:     &rtti.Zt_NumberEval,
		}, {
			Name:     "by_step",
			Label:    "by_step",
			Optional: true,
			Type:     &rtti.Zt_NumberEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumListEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Generates a series of numbers r[i] = (start + step*i) where i>=0.", "Start and step default to 1, stop defaults to start;", "the inputs are truncated to produce whole numbers;", "a zero step returns an error.", "A positive step ends the series when the returned value would exceed stop", "while a negative step ends before generating a value less than stop."},
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "list",
	Comment: []string{
		"list",
	},

	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_EraseEdge,
	&Zt_EraseIndex,
	&Zt_Erasing,
	&Zt_ErasingEdge,
	&Zt_ListEach,
	&Zt_ListFind,
	&Zt_ListGather,
	&Zt_ListLen,
	&Zt_MakeTextList,
	&Zt_MakeNumList,
	&Zt_MakeRecordList,
	&Zt_ListMap,
	&Zt_ListReduce,
	&Zt_ListReverse,
	&Zt_ListSlice,
	&Zt_ListSortNumbers,
	&Zt_ListSortText,
	&Zt_ListSplice,
	&Zt_ListPush,
	&Zt_Range,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	6334415563934548256:  (*ListGather)(nil),      /* Gather:from:using: */
	17857642077015906043: (*EraseEdge)(nil),       /* execute=Erase: */
	4911242881414594201:  (*EraseEdge)(nil),       /* execute=Erase:atFront: */
	13326390992756169124: (*EraseIndex)(nil),      /* execute=Erase:from:atIndex: */
	15309883842271607141: (*ErasingEdge)(nil),     /* execute=Erasing:as:do: */
	2341467540630172606:  (*ErasingEdge)(nil),     /* execute=Erasing:as:do:else: */
	7006351070379896671:  (*ErasingEdge)(nil),     /* execute=Erasing:atFront:as:do: */
	12034742036302137452: (*ErasingEdge)(nil),     /* execute=Erasing:atFront:as:do:else: */
	1044384912965145788:  (*Erasing)(nil),         /* execute=Erasing:from:atIndex:as:do: */
	8547752949201735569:  (*ListFind)(nil),        /* bool_eval=Find:inList: */
	16815906459082105780: (*ListFind)(nil),        /* number_eval=Find:inList: */
	3478260273963207965:  (*ListLen)(nil),         /* number_eval=Len: */
	11141869806069158915: (*MakeNumList)(nil),     /* num_list_eval=List ofNumbers: */
	10609280349940760977: (*MakeRecordList)(nil),  /* record_list_eval=List ofRecords:ofType: */
	15650595833095485421: (*MakeTextList)(nil),    /* text_list_eval=List ofText: */
	8449127989109999373:  (*ListMap)(nil),         /* execute=Map:fromList:using: */
	14590825769568398889: (*ListPush)(nil),        /* execute=Push:into: */
	17497959320325918107: (*ListPush)(nil),        /* execute=Push:into:atFront: */
	120416590109430143:   (*Range)(nil),           /* num_list_eval=Range: */
	15503705420922978310: (*Range)(nil),           /* num_list_eval=Range:byStep: */
	16618866959380663563: (*Range)(nil),           /* num_list_eval=Range:from: */
	14227857065891717050: (*Range)(nil),           /* num_list_eval=Range:from:byStep: */
	18245549119758376391: (*ListReduce)(nil),      /* execute=Reduce into:fromList:using: */
	7084717997213120806:  (*ListEach)(nil),        /* execute=Repeating across:as:do: */
	12445157229684471803: (*ListEach)(nil),        /* execute=Repeating across:as:do:else: */
	177314099445105829:   (*ListReverse)(nil),     /* execute=Reverse list: */
	4235921801420235638:  (*ListSlice)(nil),       /* num_list_eval=Slice: */
	13273073049578089927: (*ListSlice)(nil),       /* record_list_eval=Slice: */
	18323981472330239313: (*ListSlice)(nil),       /* text_list_eval=Slice: */
	3713929053224137387:  (*ListSlice)(nil),       /* num_list_eval=Slice:end: */
	326673439235441194:   (*ListSlice)(nil),       /* record_list_eval=Slice:end: */
	8469880138850798532:  (*ListSlice)(nil),       /* text_list_eval=Slice:end: */
	6763121597476813124:  (*ListSlice)(nil),       /* num_list_eval=Slice:start: */
	10126987075066562677: (*ListSlice)(nil),       /* record_list_eval=Slice:start: */
	2045310658543284955:  (*ListSlice)(nil),       /* text_list_eval=Slice:start: */
	14495675636779114361: (*ListSlice)(nil),       /* num_list_eval=Slice:start:end: */
	3241896595896148736:  (*ListSlice)(nil),       /* record_list_eval=Slice:start:end: */
	8901512565003460886:  (*ListSlice)(nil),       /* text_list_eval=Slice:start:end: */
	2873147130324862012:  (*ListSortNumbers)(nil), /* execute=SortNumbers:byField: */
	16697045456605499852: (*ListSortNumbers)(nil), /* execute=SortNumbers:byField:descending: */
	16004888373963195994: (*ListSortText)(nil),    /* execute=SortTexts:byField: */
	10015011362106184366: (*ListSortText)(nil),    /* execute=SortTexts:byField:descending: */
	10595238214248400404: (*ListSortText)(nil),    /* execute=SortTexts:byField:descending:usingCase: */
	10680774202307610784: (*ListSortText)(nil),    /* execute=SortTexts:byField:usingCase: */
	13203130291219794646: (*ListSplice)(nil),      /* execute=Splice:start:remove:insert: */
	6201472222981604265:  (*ListSplice)(nil),      /* num_list_eval=Splice:start:remove:insert: */
	15778591428898251294: (*ListSplice)(nil),      /* record_list_eval=Splice:start:remove:insert: */
	11160578659475180120: (*ListSplice)(nil),      /* text_list_eval=Splice:start:remove:insert: */
}
