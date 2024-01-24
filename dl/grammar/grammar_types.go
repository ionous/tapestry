// Code generated by Tapestry; edit at your own risk.
package grammar

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// grammar_maker, a type of slot.
var Zt_GrammarMaker = typeinfo.Slot{
	Name: "grammar_maker",
	Markup: map[string]any{
		"comment": "Helper for defining parser grammars.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_GrammarMaker_Slot struct{ Value GrammarMaker }

// implements typeinfo.Inspector for a single slot.
func (*FIX_GrammarMaker_Slot) Inspect() typeinfo.T {
	return &Zt_GrammarMaker
}

// holds a slice of slots
type GrammarMaker_Slots []GrammarMaker

// implements typeinfo.Inspector for a series of slots.
func (*GrammarMaker_Slots) Inspect() typeinfo.T {
	return &Zt_GrammarMaker
}

// scanner_maker, a type of slot.
var Zt_ScannerMaker = typeinfo.Slot{
	Name: "scanner_maker",
	Markup: map[string]any{
		"comment": "Helper for defining input scanners.",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_ScannerMaker_Slot struct{ Value ScannerMaker }

// implements typeinfo.Inspector for a single slot.
func (*FIX_ScannerMaker_Slot) Inspect() typeinfo.T {
	return &Zt_ScannerMaker
}

// holds a slice of slots
type ScannerMaker_Slots []ScannerMaker

// implements typeinfo.Inspector for a series of slots.
func (*ScannerMaker_Slots) Inspect() typeinfo.T {
	return &Zt_ScannerMaker
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Action struct {
	Action    string
	Arguments assign.Arg
	Markup    map[string]any
}

// implements typeinfo.Inspector
func (*Action) Inspect() typeinfo.T {
	return &Zt_Action
}

// return a valid markup map, creating it if necessary.
func (op *Action) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Action)(nil)

// action, a type of flow.
var Zt_Action = typeinfo.Flow{
	Name: "action",
	Lede: "action",
	Terms: []typeinfo.Term{{
		Name:  "action",
		Label: "_",
		Type:  &prim.Zt_Text,
	}, {
		Name:     "arguments",
		Label:    "args",
		Optional: true,
		Repeats:  true,
		Type:     &assign.Zt_Arg,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "makes a parser scanner producing a script defined action.",
	},
}

// holds a slice of type action
// FIX: duplicates the spec decl.
type FIX_Action_Slice []Action

// implements typeinfo.Inspector
func (*Action_Slice) Inspect() typeinfo.T {
	return &Zt_Action
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Sequence struct {
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Sequence) Inspect() typeinfo.T {
	return &Zt_Sequence
}

// return a valid markup map, creating it if necessary.
func (op *Sequence) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Sequence)(nil)

// sequence, a type of flow.
var Zt_Sequence = typeinfo.Flow{
	Name: "sequence",
	Lede: "sequence",
	Terms: []typeinfo.Term{{
		Name:    "series",
		Label:   "_",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "makes a parser scanner.",
	},
}

// holds a slice of type sequence
// FIX: duplicates the spec decl.
type FIX_Sequence_Slice []Sequence

// implements typeinfo.Inspector
func (*Sequence_Slice) Inspect() typeinfo.T {
	return &Zt_Sequence
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ChooseOne struct {
	Options ScannerMaker
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*ChooseOne) Inspect() typeinfo.T {
	return &Zt_ChooseOne
}

// return a valid markup map, creating it if necessary.
func (op *ChooseOne) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*ChooseOne)(nil)

// choose_one, a type of flow.
var Zt_ChooseOne = typeinfo.Flow{
	Name: "choose_one",
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:    "options",
		Label:   "of",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "makes a parser scanner.",
	},
}

// holds a slice of type choose_one
// FIX: duplicates the spec decl.
type FIX_ChooseOne_Slice []ChooseOne

// implements typeinfo.Inspector
func (*ChooseOne_Slice) Inspect() typeinfo.T {
	return &Zt_ChooseOne
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Directive struct {
	Name   string
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Directive) Inspect() typeinfo.T {
	return &Zt_Directive
}

// return a valid markup map, creating it if necessary.
func (op *Directive) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ GrammarMaker = (*Directive)(nil)

// directive, a type of flow.
var Zt_Directive = typeinfo.Flow{
	Name: "directive",
	Lede: "interpret",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "name",
		Type:  &prim.Zt_Text,
	}, {
		Name:    "series",
		Label:   "with",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_GrammarMaker,
	},
	Markup: map[string]any{
		"comment": "starts a parser scanner.",
	},
}

// holds a slice of type directive
// FIX: duplicates the spec decl.
type FIX_Directive_Slice []Directive

// implements typeinfo.Inspector
func (*Directive_Slice) Inspect() typeinfo.T {
	return &Zt_Directive
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Noun struct {
	Kind   string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Noun) Inspect() typeinfo.T {
	return &Zt_Noun
}

// return a valid markup map, creating it if necessary.
func (op *Noun) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Noun)(nil)

// noun, a type of flow.
var Zt_Noun = typeinfo.Flow{
	Name: "noun",
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:  "kind",
		Label: "noun",
		Type:  &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "makes a parser scanner.",
	},
}

// holds a slice of type noun
// FIX: duplicates the spec decl.
type FIX_Noun_Slice []Noun

// implements typeinfo.Inspector
func (*Noun_Slice) Inspect() typeinfo.T {
	return &Zt_Noun
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Refine struct {
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Refine) Inspect() typeinfo.T {
	return &Zt_Refine
}

// return a valid markup map, creating it if necessary.
func (op *Refine) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Refine)(nil)

// refine, a type of flow.
var Zt_Refine = typeinfo.Flow{
	Name: "refine",
	Lede: "refine",
	Terms: []typeinfo.Term{{
		Name:    "series",
		Label:   "sequence",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "Change to the bounds of the most recent result.",
	},
}

// holds a slice of type refine
// FIX: duplicates the spec decl.
type FIX_Refine_Slice []Refine

// implements typeinfo.Inspector
func (*Refine_Slice) Inspect() typeinfo.T {
	return &Zt_Refine
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Reverse struct {
	Reverses ScannerMaker
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*Reverse) Inspect() typeinfo.T {
	return &Zt_Reverse
}

// return a valid markup map, creating it if necessary.
func (op *Reverse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Reverse)(nil)

// reverse, a type of flow.
var Zt_Reverse = typeinfo.Flow{
	Name: "reverse",
	Lede: "reverse",
	Terms: []typeinfo.Term{{
		Name:    "reverses",
		Label:   "_",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "Swap the first and last matching results.",
	},
}

// holds a slice of type reverse
// FIX: duplicates the spec decl.
type FIX_Reverse_Slice []Reverse

// implements typeinfo.Inspector
func (*Reverse_Slice) Inspect() typeinfo.T {
	return &Zt_Reverse
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Focus struct {
	Player string
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Focus) Inspect() typeinfo.T {
	return &Zt_Focus
}

// return a valid markup map, creating it if necessary.
func (op *Focus) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Focus)(nil)

// focus, a type of flow.
var Zt_Focus = typeinfo.Flow{
	Name: "focus",
	Lede: "focus",
	Terms: []typeinfo.Term{{
		Name:  "player",
		Label: "_",
		Type:  &prim.Zt_Text,
	}, {
		Name:    "series",
		Label:   "sequence",
		Repeats: true,
		Type:    &Zt_ScannerMaker,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "Select a specific set of bounds for the scanner.",
	},
}

// holds a slice of type focus
// FIX: duplicates the spec decl.
type FIX_Focus_Slice []Focus

// implements typeinfo.Inspector
func (*Focus_Slice) Inspect() typeinfo.T {
	return &Zt_Focus
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Words struct {
	Words  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Words) Inspect() typeinfo.T {
	return &Zt_Words
}

// return a valid markup map, creating it if necessary.
func (op *Words) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Words)(nil)

// words, a type of flow.
var Zt_Words = typeinfo.Flow{
	Name: "words",
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:    "words",
		Label:   "word",
		Repeats: true,
		Type:    &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_ScannerMaker,
	},
	Markup: map[string]any{
		"comment": "makes a parser scanner.",
	},
}

// holds a slice of type words
// FIX: duplicates the spec decl.
type FIX_Words_Slice []Words

// implements typeinfo.Inspector
func (*Words_Slice) Inspect() typeinfo.T {
	return &Zt_Words
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "grammar",
	Slot: z_slot_list,
	Flow: z_flow_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_GrammarMaker,
	&Zt_ScannerMaker,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Action,
	&Zt_Sequence,
	&Zt_ChooseOne,
	&Zt_Directive,
	&Zt_Noun,
	&Zt_Refine,
	&Zt_Reverse,
	&Zt_Focus,
	&Zt_Words,
}
