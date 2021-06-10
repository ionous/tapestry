// Code generated by "makeops"; edit at your own risk.
package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

// AsNum Define the name of a number variable.
type AsNum struct {
	Var value.VariableName `if:"label=_"`
}

func (*AsNum) Compose() composer.Spec {
	return composer.Spec{
		Name: "as_num",
		Uses: "flow",
		Lede: "num",
	}
}

var _ ListIterator = (*AsNum)(nil)

// AsRec Define the name of a record variable.
type AsRec struct {
	Var value.VariableName `if:"label=_"`
}

func (*AsRec) Compose() composer.Spec {
	return composer.Spec{
		Name: "as_rec",
		Uses: "flow",
		Lede: "rec",
	}
}

var _ ListIterator = (*AsRec)(nil)

// AsTxt Define the name of a text variable.
type AsTxt struct {
	Var value.VariableName `if:"label=_"`
}

func (*AsTxt) Compose() composer.Spec {
	return composer.Spec{
		Name: "as_txt",
		Uses: "flow",
		Lede: "txt",
	}
}

var _ ListIterator = (*AsTxt)(nil)

// EraseEdge Erase at edge: Remove one or more values from a list
type EraseEdge struct {
	From   ListSource  `if:"label=_"`
	AtEdge rt.BoolEval `if:"label=at_front,optional"`
}

func (*EraseEdge) Compose() composer.Spec {
	return composer.Spec{
		Name: "erase_edge",
		Uses: "flow",
		Lede: "erase",
	}
}

var _ rt.Execute = (*EraseEdge)(nil)

// EraseIndex Erase at index: Remove one or more values from a list
type EraseIndex struct {
	Count   rt.NumberEval `if:"label=_"`
	From    ListSource    `if:"label=from"`
	AtIndex rt.NumberEval `if:"label=at_index"`
}

func (*EraseIndex) Compose() composer.Spec {
	return composer.Spec{
		Name: "erase_index",
		Uses: "flow",
		Lede: "erase",
	}
}

var _ rt.Execute = (*EraseIndex)(nil)

// Erasing Erase elements from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.
type Erasing struct {
	Count   rt.NumberEval `if:"label=_"`
	From    ListSource    `if:"label=from"`
	AtIndex rt.NumberEval `if:"label=at_index"`
	As      string        `if:"label=as,type=text"`
	Do      core.Activity `if:"label=do"`
}

func (*Erasing) Compose() composer.Spec {
	return composer.Spec{
		Name: "erasing",
		Uses: "flow",
	}
}

var _ rt.Execute = (*Erasing)(nil)

// ErasingEdge Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.
type ErasingEdge struct {
	From   ListSource    `if:"label=_"`
	AtEdge rt.BoolEval   `if:"label=at_front,optional"`
	As     string        `if:"label=as,type=text"`
	Do     core.Activity `if:"label=do"`
	Else   core.Brancher `if:"label=else,optional"`
}

func (*ErasingEdge) Compose() composer.Spec {
	return composer.Spec{
		Name: "erasing_edge",
		Uses: "flow",
		Lede: "erasing",
	}
}

var _ rt.Execute = (*ErasingEdge)(nil)

// FromNumList Uses a list of numbers
type FromNumList struct {
	Var value.VariableName `if:"label=_"`
}

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_num_list",
		Uses: "flow",
		Lede: "nums",
	}
}

var _ ListSource = (*FromNumList)(nil)

// FromRecList Uses a list of records
type FromRecList struct {
	Var value.VariableName `if:"label=_"`
}

func (*FromRecList) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_rec_list",
		Uses: "flow",
		Lede: "recs",
	}
}

var _ ListSource = (*FromRecList)(nil)

// FromTxtList Uses a list of text
type FromTxtList struct {
	Var value.VariableName `if:"label=_"`
}

func (*FromTxtList) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_txt_list",
		Uses: "flow",
		Lede: "txts",
	}
}

var _ ListSource = (*FromTxtList)(nil)

// IntoNumList Targets a list of numbers
type IntoNumList struct {
	Var value.VariableName `if:"label=_"`
}

func (*IntoNumList) Compose() composer.Spec {
	return composer.Spec{
		Name: "into_num_list",
		Uses: "flow",
		Lede: "nums",
	}
}

var _ ListTarget = (*IntoNumList)(nil)

// IntoRecList Targets a list of records
type IntoRecList struct {
	Var value.VariableName `if:"label=_"`
}

func (*IntoRecList) Compose() composer.Spec {
	return composer.Spec{
		Name: "into_rec_list",
		Uses: "flow",
		Lede: "recs",
	}
}

var _ ListTarget = (*IntoRecList)(nil)

// IntoTxtList Targets a list of text
type IntoTxtList struct {
	Var value.VariableName `if:"label=_"`
}

func (*IntoTxtList) Compose() composer.Spec {
	return composer.Spec{
		Name: "into_txt_list",
		Uses: "flow",
		Lede: "txts",
	}
}

var _ ListTarget = (*IntoTxtList)(nil)

// ListAt Get a value from a list. The first element is is index 1.
type ListAt struct {
	List  rt.Assignment `if:"label=_"`
	Index rt.NumberEval `if:"label=index"`
}

func (*ListAt) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_at",
		Uses: "flow",
		Lede: "get",
	}
}

var _ rt.NumberEval = (*ListAt)(nil)
var _ rt.TextEval = (*ListAt)(nil)
var _ rt.RecordEval = (*ListAt)(nil)

// ListEach Loops over the elements in the passed list, or runs the &#x27;else&#x27; activity if empty.
type ListEach struct {
	List rt.Assignment `if:"label=across"`
	As   ListIterator  `if:"label=as"`
	Do   core.Activity `if:"label=do"`
	Else core.Brancher `if:"label=else,optional"`
}

func (*ListEach) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_each",
		Uses: "flow",
		Lede: "repeating",
	}
}

var _ rt.Execute = (*ListEach)(nil)

// ListFind Search a list for a specific value.
type ListFind struct {
	Value rt.Assignment `if:"label=_"`
	List  rt.Assignment `if:"label=list"`
}

func (*ListFind) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_find",
		Uses: "flow",
		Lede: "find",
	}
}

var _ rt.BoolEval = (*ListFind)(nil)
var _ rt.NumberEval = (*ListFind)(nil)

// ListGather Transform the values from a list. The named pattern gets called once for each value in the list. It get called with two parameters: &#x27;in&#x27; as each value from the list, and &#x27;out&#x27; as the var passed to the gather.
type ListGather struct {
	Var   value.VariableName `if:"label=_"`
	From  ListSource         `if:"label=from"`
	Using string             `if:"label=_,type=text"`
}

func (*ListGather) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_gather",
		Uses: "flow",
		Lede: "gather",
	}
}

// ListLen Determines the number of values in a list.
type ListLen struct {
	List rt.Assignment `if:"label=_"`
}

func (*ListLen) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_len",
		Uses: "flow",
		Lede: "len",
	}
}

var _ rt.NumberEval = (*ListLen)(nil)

// ListMap Transform the values from one list and place the results in another list. The designated pattern is called with each value from the &#x27;from list&#x27;, one value at a time.
type ListMap struct {
	ToList       string        `if:"label=_,type=text"`
	FromList     rt.Assignment `if:"label=from_list"`
	UsingPattern string        `if:"label=using,type=text"`
}

func (*ListMap) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_map",
		Uses: "flow",
		Lede: "map",
	}
}

var _ rt.Execute = (*ListMap)(nil)

// ListReduce Transform the values from one list by combining them into a single value. The named pattern is called with two parameters: &#x27;in&#x27; ( each element of the list ) and &#x27;out&#x27; ( ex. a record ).
type ListReduce struct {
	IntoValue    string        `if:"label=into,type=text"`
	FromList     rt.Assignment `if:"label=from_list"`
	UsingPattern string        `if:"label=using,type=text"`
}

func (*ListReduce) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_reduce",
		Uses: "flow",
		Lede: "reduce",
	}
}

var _ rt.Execute = (*ListReduce)(nil)

// ListReverse Reverse a list.
type ListReverse struct {
	List ListSource `if:"label=_"`
}

func (*ListReverse) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_reverse",
		Uses: "flow",
		Lede: "reverse",
	}
}

var _ rt.Execute = (*ListReverse)(nil)

// ListSet Overwrite an existing value in a list.
type ListSet struct {
	List  string        `if:"label=_,type=text"`
	Index rt.NumberEval `if:"label=index"`
	From  rt.Assignment `if:"label=from"`
}

func (*ListSet) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_set",
		Uses: "flow",
		Lede: "set",
	}
}

var _ rt.Execute = (*ListSet)(nil)

// ListSlice Create a new list from a section of another list.,Start is optional, if omitted slice starts at the first element.,If start is greater the length, an empty array is returned.,Slice doesnt include the ending index.,Negatives indices indicates an offset from the end.,When end is omitted, copy up to and including the last element;,and do the same if the end is greater than the length
type ListSlice struct {
	List  rt.Assignment `if:"label=_"`
	Start rt.NumberEval `if:"label=start,optional"`
	End   rt.NumberEval `if:"label=end,optional"`
}

func (*ListSlice) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_slice",
		Uses: "flow",
		Lede: "slice",
	}
}

var _ rt.NumListEval = (*ListSlice)(nil)
var _ rt.TextListEval = (*ListSlice)(nil)
var _ rt.RecordListEval = (*ListSlice)(nil)

// ListSortNumbers
type ListSortNumbers struct {
	Var        value.VariableName `if:"label=_"`
	ByField    string             `if:"label=by_field,type=text"`
	Descending rt.BoolEval        `if:"label=descending,optional"`
}

func (*ListSortNumbers) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_sort_numbers",
		Uses: "flow",
		Lede: "sort",
	}
}

var _ rt.Execute = (*ListSortNumbers)(nil)

// ListSortText Rearrange the elements in the named list by using the designated pattern to test pairs of elements.
type ListSortText struct {
	Var        value.VariableName `if:"label=_"`
	ByField    string             `if:"label=by_field,type=text"`
	Descending rt.BoolEval        `if:"label=descending,optional"`
	UsingCase  rt.BoolEval        `if:"label=using_case,optional"`
}

func (*ListSortText) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_sort_text",
		Uses: "flow",
		Lede: "sort",
	}
}

var _ rt.Execute = (*ListSortText)(nil)

// ListSortUsing Rearrange the elements in the named list by using the designated pattern to test pairs of elements.
type ListSortUsing struct {
	Var   value.VariableName `if:"label=_"`
	Using string             `if:"label=using,type=text"`
}

func (*ListSortUsing) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_sort_using",
		Uses: "flow",
		Lede: "sort",
	}
}

var _ rt.Execute = (*ListSortUsing)(nil)

// ListSplice Modify a list by adding and removing elements. Note: the type of the elements being added must match the type of the list. Text cant be added to a list of numbers, numbers cant be added to a list of text. If the starting index is negative, it will begin that many elements from the end of the array. If list&#x27;s length + the start is less than 0, it will begin from index 0. If the remove count is missing, it removes all elements from the start to the end; if it is 0 or negative, no elements are removed.
type ListSplice struct {
	List   string        `if:"label=_,type=text"`
	Start  rt.NumberEval `if:"label=start"`
	Remove rt.NumberEval `if:"label=remove"`
	Insert rt.Assignment `if:"label=insert"`
}

func (*ListSplice) Compose() composer.Spec {
	return composer.Spec{
		Name: "list_splice",
		Uses: "flow",
		Lede: "splice",
	}
}

var _ rt.Execute = (*ListSplice)(nil)
var _ rt.NumListEval = (*ListSplice)(nil)
var _ rt.TextListEval = (*ListSplice)(nil)
var _ rt.RecordListEval = (*ListSplice)(nil)

// PutEdge Add a value to a list
type PutEdge struct {
	From   rt.Assignment `if:"label=_"`
	Into   ListTarget    `if:"label=into"`
	AtEdge rt.BoolEval   `if:"label=at_front,optional"`
}

func (*PutEdge) Compose() composer.Spec {
	return composer.Spec{
		Name: "put_edge",
		Uses: "flow",
		Lede: "put",
	}
}

var _ rt.Execute = (*PutEdge)(nil)

// PutIndex Replace one value in a list with another
type PutIndex struct {
	From    rt.Assignment `if:"label=_"`
	Into    ListTarget    `if:"label=into"`
	AtIndex rt.NumberEval `if:"label=at_index"`
}

func (*PutIndex) Compose() composer.Spec {
	return composer.Spec{
		Name: "put_index",
		Uses: "flow",
		Lede: "put",
	}
}

var _ rt.Execute = (*PutIndex)(nil)

// Range Generates a series of numbers r[i] &#x3D; (start + step*i) where i&gt;&#x3D;0.,Start and step default to 1, stop defaults to start;,the inputs are truncated to produce whole numbers;,a zero step returns an error.,A positive step ends the series when the returned value would exceed stop,while a negative step ends before generating a value less than stop.
type Range struct {
	To     rt.NumberEval `if:"label=_"`
	From   rt.NumberEval `if:"label=from,optional"`
	ByStep rt.NumberEval `if:"label=by_step,optional"`
}

func (*Range) Compose() composer.Spec {
	return composer.Spec{
		Name: "range",
		Uses: "flow",
	}
}

var _ rt.NumListEval = (*Range)(nil)
var Slots = []interface{}{
	(*ListIterator)(nil),
	(*ListSource)(nil),
	(*ListTarget)(nil),
}
var Slats = []composer.Composer{
	(*AsNum)(nil),
	(*AsRec)(nil),
	(*AsTxt)(nil),
	(*EraseEdge)(nil),
	(*EraseIndex)(nil),
	(*Erasing)(nil),
	(*ErasingEdge)(nil),
	(*FromNumList)(nil),
	(*FromRecList)(nil),
	(*FromTxtList)(nil),
	(*IntoNumList)(nil),
	(*IntoRecList)(nil),
	(*IntoTxtList)(nil),
	(*ListAt)(nil),
	(*ListEach)(nil),
	(*ListFind)(nil),
	(*ListGather)(nil),
	(*ListLen)(nil),
	(*ListMap)(nil),
	(*ListReduce)(nil),
	(*ListReverse)(nil),
	(*ListSet)(nil),
	(*ListSlice)(nil),
	(*ListSortNumbers)(nil),
	(*ListSortText)(nil),
	(*ListSortUsing)(nil),
	(*ListSplice)(nil),
	(*PutEdge)(nil),
	(*PutIndex)(nil),
	(*Range)(nil),
}
