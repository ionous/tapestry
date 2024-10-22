// Boolean operations and loop handling.
package logic

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// brancher, a type of slot.
var Zt_Brancher = typeinfo.Slot{
	Name: "brancher",
	Markup: map[string]any{
		"--": "Helper for the else statements of [ChooseBranch].",
	},
}

// Holds a single slot.
type Brancher_Slot struct{ Value Brancher }

// Implements [typeinfo.Instance] for a single slot.
func (*Brancher_Slot) TypeInfo() typeinfo.T {
	return &Zt_Brancher
}

// Holds a slice of slots.
type Brancher_Slots []Brancher

// Implements [typeinfo.Instance] for a slice of slots.
func (*Brancher_Slots) TypeInfo() typeinfo.T {
	return &Zt_Brancher
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Brancher_Slots) Repeats() bool {
	return len(*op) > 0
}

// This always returns true.
type Always struct {
	Markup map[string]any `json:",omitempty"`
}

// always, a type of flow.
var Zt_Always typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Always) TypeInfo() typeinfo.T {
	return &Zt_Always
}

// Implements [typeinfo.Markup]
func (op *Always) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*Always)(nil)

// Holds a slice of type Always.
type Always_Slice []Always

// Implements [typeinfo.Instance] for a slice of Always.
func (*Always_Slice) TypeInfo() typeinfo.T {
	return &Zt_Always
}

// Implements [typeinfo.Repeats] for a slice of Always.
func (op *Always_Slice) Repeats() bool {
	return len(*op) > 0
}

// This always returns false.
type Never struct {
	Markup map[string]any `json:",omitempty"`
}

// never, a type of flow.
var Zt_Never typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Never) TypeInfo() typeinfo.T {
	return &Zt_Never
}

// Implements [typeinfo.Markup]
func (op *Never) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*Never)(nil)

// Holds a slice of type Never.
type Never_Slice []Never

// Implements [typeinfo.Instance] for a slice of Never.
func (*Never_Slice) TypeInfo() typeinfo.T {
	return &Zt_Never
}

// Implements [typeinfo.Repeats] for a slice of Never.
func (op *Never_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine the "truthiness" of a value.
// Bool values simply return their value.
// Num values: are true when not exactly zero.
// Text values: are true whenever they contain content.
// List values: are true whenever the list is non-empty.
// ( note this is similar to python, and different than javascript. )
// Record values: are true whenever they have been initialized.
// ( only sub-records start uninitialized; record variables are always true. )
type IsValue struct {
	Value  rtti.Assignment
	Markup map[string]any `json:",omitempty"`
}

// is_value, a type of flow.
var Zt_IsValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsValue) TypeInfo() typeinfo.T {
	return &Zt_IsValue
}

// Implements [typeinfo.Markup]
func (op *IsValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsValue)(nil)

// Holds a slice of type IsValue.
type IsValue_Slice []IsValue

// Implements [typeinfo.Instance] for a slice of IsValue.
func (*IsValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsValue
}

// Implements [typeinfo.Repeats] for a slice of IsValue.
func (op *IsValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Check that every condition in a set of conditions returns true.
// Stops after finding a failed condition.
// An empty list returns false.
type IsAll struct {
	Test   []rtti.BoolEval
	Markup map[string]any `json:",omitempty"`
}

// is_all, a type of flow.
var Zt_IsAll typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsAll) TypeInfo() typeinfo.T {
	return &Zt_IsAll
}

// Implements [typeinfo.Markup]
func (op *IsAll) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsAll)(nil)

// Holds a slice of type IsAll.
type IsAll_Slice []IsAll

// Implements [typeinfo.Instance] for a slice of IsAll.
func (*IsAll_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsAll
}

// Implements [typeinfo.Repeats] for a slice of IsAll.
func (op *IsAll_Slice) Repeats() bool {
	return len(*op) > 0
}

// Check whether any condition in a set of conditions returns true.
// Stops after finding the first successful condition.
// An empty list returns false.
type IsAny struct {
	Test   []rtti.BoolEval
	Markup map[string]any `json:",omitempty"`
}

// is_any, a type of flow.
var Zt_IsAny typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsAny) TypeInfo() typeinfo.T {
	return &Zt_IsAny
}

// Implements [typeinfo.Markup]
func (op *IsAny) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsAny)(nil)

// Holds a slice of type IsAny.
type IsAny_Slice []IsAny

// Implements [typeinfo.Instance] for a slice of IsAny.
func (*IsAny_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsAny
}

// Implements [typeinfo.Repeats] for a slice of IsAny.
func (op *IsAny_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine the opposite of a condition.
type Not struct {
	Test   rtti.BoolEval
	Markup map[string]any `json:",omitempty"`
}

// not, a type of flow.
var Zt_Not typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Not) TypeInfo() typeinfo.T {
	return &Zt_Not
}

// Implements [typeinfo.Markup]
func (op *Not) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*Not)(nil)

// Holds a slice of type Not.
type Not_Slice []Not

// Implements [typeinfo.Instance] for a slice of Not.
func (*Not_Slice) TypeInfo() typeinfo.T {
	return &Zt_Not
}

// Implements [typeinfo.Repeats] for a slice of Not.
func (op *Not_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine the "falsiness" of a value.
// This is the opposite of [TrueValue].
type NotValue struct {
	Value  rtti.Assignment
	Markup map[string]any `json:",omitempty"`
}

// not_value, a type of flow.
var Zt_NotValue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NotValue) TypeInfo() typeinfo.T {
	return &Zt_NotValue
}

// Implements [typeinfo.Markup]
func (op *NotValue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*NotValue)(nil)

// Holds a slice of type NotValue.
type NotValue_Slice []NotValue

// Implements [typeinfo.Instance] for a slice of NotValue.
func (*NotValue_Slice) TypeInfo() typeinfo.T {
	return &Zt_NotValue
}

// Implements [typeinfo.Repeats] for a slice of NotValue.
func (op *NotValue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Check that every condition in a set of conditions returns false.
// Stops after finding any successful condition.
// An empty list returns false.
type NotAll struct {
	Test   []rtti.BoolEval
	Markup map[string]any `json:",omitempty"`
}

// not_all, a type of flow.
var Zt_NotAll typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NotAll) TypeInfo() typeinfo.T {
	return &Zt_NotAll
}

// Implements [typeinfo.Markup]
func (op *NotAll) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*NotAll)(nil)

// Holds a slice of type NotAll.
type NotAll_Slice []NotAll

// Implements [typeinfo.Instance] for a slice of NotAll.
func (*NotAll_Slice) TypeInfo() typeinfo.T {
	return &Zt_NotAll
}

// Implements [typeinfo.Repeats] for a slice of NotAll.
func (op *NotAll_Slice) Repeats() bool {
	return len(*op) > 0
}

// Check whether any condition in a set of conditions returns false.
// Stops after finding any failed condition.
// An empty list returns false.
type NotAny struct {
	Test   []rtti.BoolEval
	Markup map[string]any `json:",omitempty"`
}

// not_any, a type of flow.
var Zt_NotAny typeinfo.Flow

// Implements [typeinfo.Instance]
func (*NotAny) TypeInfo() typeinfo.T {
	return &Zt_NotAny
}

// Implements [typeinfo.Markup]
func (op *NotAny) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*NotAny)(nil)

// Holds a slice of type NotAny.
type NotAny_Slice []NotAny

// Implements [typeinfo.Instance] for a slice of NotAny.
func (*NotAny_Slice) TypeInfo() typeinfo.T {
	return &Zt_NotAny
}

// Implements [typeinfo.Repeats] for a slice of NotAny.
func (op *NotAny_Slice) Repeats() bool {
	return len(*op) > 0
}

// Select a block of statements to run based on a true/false check.
type ChooseBranch struct {
	Condition rtti.BoolEval
	Args      []call.Arg
	Exe       []rtti.Execute
	Else      Brancher
	Markup    map[string]any `json:",omitempty"`
}

// choose_branch, a type of flow.
var Zt_ChooseBranch typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ChooseBranch) TypeInfo() typeinfo.T {
	return &Zt_ChooseBranch
}

// Implements [typeinfo.Markup]
func (op *ChooseBranch) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ChooseBranch)(nil)
var _ Brancher = (*ChooseBranch)(nil)

// Holds a slice of type ChooseBranch.
type ChooseBranch_Slice []ChooseBranch

// Implements [typeinfo.Instance] for a slice of ChooseBranch.
func (*ChooseBranch_Slice) TypeInfo() typeinfo.T {
	return &Zt_ChooseBranch
}

// Implements [typeinfo.Repeats] for a slice of ChooseBranch.
func (op *ChooseBranch_Slice) Repeats() bool {
	return len(*op) > 0
}

// Run a set of statements after a condition has failed.
type ChooseNothingElse struct {
	Exe    []rtti.Execute
	Markup map[string]any `json:",omitempty"`
}

// choose_nothing_else, a type of flow.
var Zt_ChooseNothingElse typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ChooseNothingElse) TypeInfo() typeinfo.T {
	return &Zt_ChooseNothingElse
}

// Implements [typeinfo.Markup]
func (op *ChooseNothingElse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Brancher = (*ChooseNothingElse)(nil)

// Holds a slice of type ChooseNothingElse.
type ChooseNothingElse_Slice []ChooseNothingElse

// Implements [typeinfo.Instance] for a slice of ChooseNothingElse.
func (*ChooseNothingElse_Slice) TypeInfo() typeinfo.T {
	return &Zt_ChooseNothingElse
}

// Implements [typeinfo.Repeats] for a slice of ChooseNothingElse.
func (op *ChooseNothingElse_Slice) Repeats() bool {
	return len(*op) > 0
}

// Pick one of two possible text values based on a condition.
// ( This acts similar to a ternary. )
type ChooseNum struct {
	If     rtti.BoolEval
	Args   []call.Arg
	True   rtti.NumEval
	False  rtti.NumEval
	Markup map[string]any `json:",omitempty"`
}

// choose_num, a type of flow.
var Zt_ChooseNum typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ChooseNum) TypeInfo() typeinfo.T {
	return &Zt_ChooseNum
}

// Implements [typeinfo.Markup]
func (op *ChooseNum) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumEval = (*ChooseNum)(nil)

// Holds a slice of type ChooseNum.
type ChooseNum_Slice []ChooseNum

// Implements [typeinfo.Instance] for a slice of ChooseNum.
func (*ChooseNum_Slice) TypeInfo() typeinfo.T {
	return &Zt_ChooseNum
}

// Implements [typeinfo.Repeats] for a slice of ChooseNum.
func (op *ChooseNum_Slice) Repeats() bool {
	return len(*op) > 0
}

// Pick one of two possible text values based on a condition.
// ( This acts similar to a ternary. )
type ChooseText struct {
	If     rtti.BoolEval
	Args   []call.Arg
	True   rtti.TextEval
	False  rtti.TextEval
	Markup map[string]any `json:",omitempty"`
}

// choose_text, a type of flow.
var Zt_ChooseText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ChooseText) TypeInfo() typeinfo.T {
	return &Zt_ChooseText
}

// Implements [typeinfo.Markup]
func (op *ChooseText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*ChooseText)(nil)

// Holds a slice of type ChooseText.
type ChooseText_Slice []ChooseText

// Implements [typeinfo.Instance] for a slice of ChooseText.
func (*ChooseText_Slice) TypeInfo() typeinfo.T {
	return &Zt_ChooseText
}

// Implements [typeinfo.Repeats] for a slice of ChooseText.
func (op *ChooseText_Slice) Repeats() bool {
	return len(*op) > 0
}

// In a repeating loop, exit the loop;
// or, in a rule, stop processing rules.
type Break struct {
	Markup map[string]any `json:",omitempty"`
}

// break, a type of flow.
var Zt_Break typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Break) TypeInfo() typeinfo.T {
	return &Zt_Break
}

// Implements [typeinfo.Markup]
func (op *Break) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*Break)(nil)

// Holds a slice of type Break.
type Break_Slice []Break

// Implements [typeinfo.Instance] for a slice of Break.
func (*Break_Slice) TypeInfo() typeinfo.T {
	return &Zt_Break
}

// Implements [typeinfo.Repeats] for a slice of Break.
func (op *Break_Slice) Repeats() bool {
	return len(*op) > 0
}

// In a repeating loop, try the next iteration of the loop;
// or, in a rule, continue to the next rule.
type Continue struct {
	Markup map[string]any `json:",omitempty"`
}

// continue, a type of flow.
var Zt_Continue typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Continue) TypeInfo() typeinfo.T {
	return &Zt_Continue
}

// Implements [typeinfo.Markup]
func (op *Continue) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*Continue)(nil)

// Holds a slice of type Continue.
type Continue_Slice []Continue

// Implements [typeinfo.Instance] for a slice of Continue.
func (*Continue_Slice) TypeInfo() typeinfo.T {
	return &Zt_Continue
}

// Implements [typeinfo.Repeats] for a slice of Continue.
func (op *Continue_Slice) Repeats() bool {
	return len(*op) > 0
}

// Keep running a series of actions while a condition succeeds.
type Repeat struct {
	Condition rtti.BoolEval
	Initial   []call.Arg
	Args      []call.Arg
	Exe       []rtti.Execute
	Markup    map[string]any `json:",omitempty"`
}

// repeat, a type of flow.
var Zt_Repeat typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Repeat) TypeInfo() typeinfo.T {
	return &Zt_Repeat
}

// Implements [typeinfo.Markup]
func (op *Repeat) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*Repeat)(nil)

// Holds a slice of type Repeat.
type Repeat_Slice []Repeat

// Implements [typeinfo.Instance] for a slice of Repeat.
func (*Repeat_Slice) TypeInfo() typeinfo.T {
	return &Zt_Repeat
}

// Implements [typeinfo.Repeats] for a slice of Repeat.
func (op *Repeat_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_Always = typeinfo.Flow{
		Name:  "always",
		Lede:  "always",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": "This always returns true.",
		},
	}
	Zt_Never = typeinfo.Flow{
		Name:  "never",
		Lede:  "never",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": "This always returns false.",
		},
	}
	Zt_IsValue = typeinfo.Flow{
		Name: "is_value",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"--": "The value to test.",
			},
			Type: &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Determine the \"truthiness\" of a value.", "Bool values simply return their value.", "Num values: are true when not exactly zero.", "Text values: are true whenever they contain content.", "List values: are true whenever the list is non-empty.", "( note this is similar to python, and different than javascript. )", "Record values: are true whenever they have been initialized.", "( only sub-records start uninitialized; record variables are always true. )"},
		},
	}
	Zt_IsAll = typeinfo.Flow{
		Name: "is_all",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:    "test",
			Label:   "all",
			Repeats: true,
			Markup: map[string]any{
				"--": "One or more conditions to check for success.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Check that every condition in a set of conditions returns true.", "Stops after finding a failed condition.", "An empty list returns false."},
		},
	}
	Zt_IsAny = typeinfo.Flow{
		Name: "is_any",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:    "test",
			Label:   "any",
			Repeats: true,
			Markup: map[string]any{
				"--": "One or more conditions to check for success.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Check whether any condition in a set of conditions returns true.", "Stops after finding the first successful condition.", "An empty list returns false."},
		},
	}
	Zt_Not = typeinfo.Flow{
		Name: "not",
		Lede: "not",
		Terms: []typeinfo.Term{{
			Name: "test",
			Markup: map[string]any{
				"--": "The condition to check.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": "Determine the opposite of a condition.",
		},
	}
	Zt_NotValue = typeinfo.Flow{
		Name: "not_value",
		Lede: "not",
		Terms: []typeinfo.Term{{
			Name:  "value",
			Label: "value",
			Markup: map[string]any{
				"--": "The value to test.",
			},
			Type: &rtti.Zt_Assignment,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Determine the \"falsiness\" of a value.", "This is the opposite of [TrueValue]."},
		},
	}
	Zt_NotAll = typeinfo.Flow{
		Name: "not_all",
		Lede: "not",
		Terms: []typeinfo.Term{{
			Name:    "test",
			Label:   "all",
			Repeats: true,
			Markup: map[string]any{
				"--": "One or more conditions to check for failure.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Check that every condition in a set of conditions returns false.", "Stops after finding any successful condition.", "An empty list returns false."},
		},
	}
	Zt_NotAny = typeinfo.Flow{
		Name: "not_any",
		Lede: "not",
		Terms: []typeinfo.Term{{
			Name:    "test",
			Label:   "any",
			Repeats: true,
			Markup: map[string]any{
				"--": "One or more conditions to check for failure.",
			},
			Type: &rtti.Zt_BoolEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"--": []string{"Check whether any condition in a set of conditions returns false.", "Stops after finding any failed condition.", "An empty list returns false."},
		},
	}
	Zt_ChooseBranch = typeinfo.Flow{
		Name: "choose_branch",
		Lede: "if",
		Terms: []typeinfo.Term{{
			Name: "condition",
			Markup: map[string]any{
				"--": "The condition to test.",
			},
			Type: &rtti.Zt_BoolEval,
		}, {
			Name:     "args",
			Label:    "assuming",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"--": "A set of local variables available while testing the condition and while running the do/else statements. These are initialized before testing the condition.",
			},
			Type: &call.Zt_Arg,
		}, {
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"--": "Statements which run when the condition succeeded.",
			},
			Type: &rtti.Zt_Execute,
		}, {
			Name:     "else",
			Label:    "else",
			Optional: true,
			Markup: map[string]any{
				"--": "An optional set of statements to evaluate when the condition failed.",
			},
			Type: &Zt_Brancher,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
			&Zt_Brancher,
		},
		Markup: map[string]any{
			"--": "Select a block of statements to run based on a true/false check.",
		},
	}
	Zt_ChooseNothingElse = typeinfo.Flow{
		Name: "choose_nothing_else",
		Lede: "finally",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"--": "One or more statements to run.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Brancher,
		},
		Markup: map[string]any{
			"--": "Run a set of statements after a condition has failed.",
		},
	}
	Zt_ChooseNum = typeinfo.Flow{
		Name: "choose_num",
		Lede: "num",
		Terms: []typeinfo.Term{{
			Name:  "if",
			Label: "if",
			Markup: map[string]any{
				"--": "The condition to test.",
			},
			Type: &rtti.Zt_BoolEval,
		}, {
			Name:     "args",
			Label:    "assuming",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"--": "A set of local variables available while testing the condition and while running the do/else statements. These are initialized before testing the condition.",
			},
			Type: &call.Zt_Arg,
		}, {
			Name:  "true",
			Label: "then",
			Markup: map[string]any{
				"--": []string{"The number to use if the condition succeeds.", "( The eval is only processed if the condition succeeded. )"},
			},
			Type: &rtti.Zt_NumEval,
		}, {
			Name:     "false",
			Label:    "else",
			Optional: true,
			Markup: map[string]any{
				"--": []string{"The number to use if the condition fails.", "( The eval is only processed if the condition failed. )"},
			},
			Type: &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumEval,
		},
		Markup: map[string]any{
			"--": []string{"Pick one of two possible text values based on a condition.", "( This acts similar to a ternary. )"},
		},
	}
	Zt_ChooseText = typeinfo.Flow{
		Name: "choose_text",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:  "if",
			Label: "if",
			Markup: map[string]any{
				"--": "The condition to test.",
			},
			Type: &rtti.Zt_BoolEval,
		}, {
			Name:     "args",
			Label:    "assuming",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"--": "A set of local variables available while testing the condition and while running the do/else statements. These are initialized before testing the condition.",
			},
			Type: &call.Zt_Arg,
		}, {
			Name:  "true",
			Label: "then",
			Markup: map[string]any{
				"--": []string{"The text to use if the condition succeeds.", "( The eval is only processed if the condition succeeded. )"},
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:     "false",
			Label:    "else",
			Optional: true,
			Markup: map[string]any{
				"--": []string{"The text to use if the condition fails.", "( The eval is only processed if the condition failed. )"},
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"--": []string{"Pick one of two possible text values based on a condition.", "( This acts similar to a ternary. )"},
		},
	}
	Zt_Break = typeinfo.Flow{
		Name:  "break",
		Lede:  "break",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"--": []string{"In a repeating loop, exit the loop;", "or, in a rule, stop processing rules."},
		},
	}
	Zt_Continue = typeinfo.Flow{
		Name:  "continue",
		Lede:  "continue",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"--": []string{"In a repeating loop, try the next iteration of the loop;", "or, in a rule, continue to the next rule."},
		},
	}
	Zt_Repeat = typeinfo.Flow{
		Name: "repeat",
		Lede: "repeat",
		Terms: []typeinfo.Term{{
			Name:  "condition",
			Label: "if",
			Markup: map[string]any{
				"--": []string{"The condition to check before running the loop.", "If it succeeds, execute the loop;", "if it fails, stop executing the loop."},
			},
			Type: &rtti.Zt_BoolEval,
		}, {
			Name:     "initial",
			Label:    "initially",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"--": []string{"A set of variables available to the loop;", "evaluated just once, before the loop's first run."},
			},
			Type: &call.Zt_Arg,
		}, {
			Name:     "args",
			Label:    "assuming",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"--": []string{"A set of variables available to the loop;", "evaluated before each iteration of the loop.", "( These take precedence over the initial variables.", "If the same names appear in both sets of variables, the ones here win. )"},
			},
			Type: &call.Zt_Arg,
		}, {
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"--": "The statements to execute every loop.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"--": "Keep running a series of actions while a condition succeeds.",
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "logic",
	Comment: []string{
		"Boolean operations and loop handling.",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Brancher,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Always,
	&Zt_Never,
	&Zt_IsValue,
	&Zt_IsAll,
	&Zt_IsAny,
	&Zt_Not,
	&Zt_NotValue,
	&Zt_NotAll,
	&Zt_NotAny,
	&Zt_ChooseBranch,
	&Zt_ChooseNothingElse,
	&Zt_ChooseNum,
	&Zt_ChooseText,
	&Zt_Break,
	&Zt_Continue,
	&Zt_Repeat,
}

// gob like registration
func Register(reg func(any)) {
	reg((*Always)(nil))
	reg((*Never)(nil))
	reg((*IsValue)(nil))
	reg((*IsAll)(nil))
	reg((*IsAny)(nil))
	reg((*Not)(nil))
	reg((*NotValue)(nil))
	reg((*NotAll)(nil))
	reg((*NotAny)(nil))
	reg((*ChooseBranch)(nil))
	reg((*ChooseNothingElse)(nil))
	reg((*ChooseNum)(nil))
	reg((*ChooseText)(nil))
	reg((*Break)(nil))
	reg((*Continue)(nil))
	reg((*Repeat)(nil))
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	1979437068831463006:  (*Always)(nil),            /* bool_eval=Always */
	9570569845423374482:  (*Break)(nil),             /* execute=Break */
	3156233792812716886:  (*Continue)(nil),          /* execute=Continue */
	13697022905922221509: (*ChooseNothingElse)(nil), /* brancher=Finally do: */
	6524366950360243674:  (*ChooseBranch)(nil),      /* brancher=If:assuming:do: */
	12195526980856142720: (*ChooseBranch)(nil),      /* execute=If:assuming:do: */
	16752471159562852415: (*ChooseBranch)(nil),      /* brancher=If:assuming:do:else: */
	2092791308408463217:  (*ChooseBranch)(nil),      /* execute=If:assuming:do:else: */
	11676187955438326921: (*ChooseBranch)(nil),      /* brancher=If:do: */
	16551038912311542599: (*ChooseBranch)(nil),      /* execute=If:do: */
	11846460753008131314: (*ChooseBranch)(nil),      /* brancher=If:do:else: */
	9882017885672780228:  (*ChooseBranch)(nil),      /* execute=If:do:else: */
	9557764360653562078:  (*IsAll)(nil),             /* bool_eval=Is all: */
	10841523351090882945: (*IsAny)(nil),             /* bool_eval=Is any: */
	14617237694045471748: (*IsValue)(nil),           /* bool_eval=Is value: */
	1310533520550597035:  (*Never)(nil),             /* bool_eval=Never */
	1619507005253470411:  (*NotAll)(nil),            /* bool_eval=Not all: */
	588371806463953504:   (*NotAny)(nil),            /* bool_eval=Not any: */
	14241778230257487825: (*NotValue)(nil),          /* bool_eval=Not value: */
	3572677870333466638:  (*Not)(nil),               /* bool_eval=Not: */
	12220459187031741460: (*ChooseNum)(nil),         /* num_eval=Num if:assuming:then: */
	2863639051637372837:  (*ChooseNum)(nil),         /* num_eval=Num if:assuming:then:else: */
	9841785069654362751:  (*ChooseNum)(nil),         /* num_eval=Num if:then: */
	2293377426593441548:  (*ChooseNum)(nil),         /* num_eval=Num if:then:else: */
	765031688964517969:   (*Repeat)(nil),            /* execute=Repeat if:assuming:do: */
	8691975868338143968:  (*Repeat)(nil),            /* execute=Repeat if:do: */
	3143868389463768902:  (*Repeat)(nil),            /* execute=Repeat if:initially:assuming:do: */
	6032823766141856509:  (*Repeat)(nil),            /* execute=Repeat if:initially:do: */
	4784360512497235820:  (*ChooseText)(nil),        /* text_eval=Text if:assuming:then: */
	13980719859951632205: (*ChooseText)(nil),        /* text_eval=Text if:assuming:then:else: */
	4706788097495762503:  (*ChooseText)(nil),        /* text_eval=Text if:then: */
	12221021609112050372: (*ChooseText)(nil),        /* text_eval=Text if:then:else: */
}
