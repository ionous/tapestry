// Code generated by Tapestry; edit at your own risk.
package grammar

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// grammar_maker, a type of slot.
const Z_GrammarMaker_Type = "grammar_maker"

var Z_GrammarMaker_Info = typeinfo.Slot{
	Name: Z_GrammarMaker_Type,
}

// holds a slice of slots
type GrammarMaker_Slots []GrammarMaker

// implements typeinfo.Inspector
func (*GrammarMaker_Slots) Inspect() typeinfo.T {
	return &Z_GrammarMaker_Info
}

// scanner_maker, a type of slot.
const Z_ScannerMaker_Type = "scanner_maker"

var Z_ScannerMaker_Info = typeinfo.Slot{
	Name: Z_ScannerMaker_Type,
}

// holds a slice of slots
type ScannerMaker_Slots []ScannerMaker

// implements typeinfo.Inspector
func (*ScannerMaker_Slots) Inspect() typeinfo.T {
	return &Z_ScannerMaker_Info
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
	return &Z_Action_Info
}

// return a valid markup map, creating it if necessary.
func (op *Action) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// action, a type of flow.
const Z_Action_Type = "action"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Action)(nil)

var Z_Action_Info = typeinfo.Flow{
	Name: Z_Action_Type,
	Lede: "action",
	Terms: []typeinfo.Term{{
		Name:  "action",
		Label: "_",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:     "arguments",
		Label:    "args",
		Optional: true,
		Repeats:  true,
		Type:     &assign.Z_Arg_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type action
// FIX: duplicates the spec decl.
type FIX_Action_Slice []Action

// implements typeinfo.Inspector
func (*Action_Slice) Inspect() typeinfo.T {
	return &Z_Action_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Sequence struct {
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Sequence) Inspect() typeinfo.T {
	return &Z_Sequence_Info
}

// return a valid markup map, creating it if necessary.
func (op *Sequence) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// sequence, a type of flow.
const Z_Sequence_Type = "sequence"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Sequence)(nil)

var Z_Sequence_Info = typeinfo.Flow{
	Name: Z_Sequence_Type,
	Lede: "sequence",
	Terms: []typeinfo.Term{{
		Name:    "series",
		Label:   "_",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type sequence
// FIX: duplicates the spec decl.
type FIX_Sequence_Slice []Sequence

// implements typeinfo.Inspector
func (*Sequence_Slice) Inspect() typeinfo.T {
	return &Z_Sequence_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_ChooseOne struct {
	Options ScannerMaker
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*ChooseOne) Inspect() typeinfo.T {
	return &Z_ChooseOne_Info
}

// return a valid markup map, creating it if necessary.
func (op *ChooseOne) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// choose_one, a type of flow.
const Z_ChooseOne_Type = "choose_one"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*ChooseOne)(nil)

var Z_ChooseOne_Info = typeinfo.Flow{
	Name: Z_ChooseOne_Type,
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:    "options",
		Label:   "of",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type choose_one
// FIX: duplicates the spec decl.
type FIX_ChooseOne_Slice []ChooseOne

// implements typeinfo.Inspector
func (*ChooseOne_Slice) Inspect() typeinfo.T {
	return &Z_ChooseOne_Info
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
	return &Z_Directive_Info
}

// return a valid markup map, creating it if necessary.
func (op *Directive) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// directive, a type of flow.
const Z_Directive_Type = "directive"

// ensure the command implements its specified slots:
var _ GrammarMaker = (*Directive)(nil)

var Z_Directive_Info = typeinfo.Flow{
	Name: Z_Directive_Type,
	Lede: "interpret",
	Terms: []typeinfo.Term{{
		Name:  "name",
		Label: "name",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:    "series",
		Label:   "with",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_GrammarMaker_Info,
	},
}

// holds a slice of type directive
// FIX: duplicates the spec decl.
type FIX_Directive_Slice []Directive

// implements typeinfo.Inspector
func (*Directive_Slice) Inspect() typeinfo.T {
	return &Z_Directive_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Noun struct {
	Kind   string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Noun) Inspect() typeinfo.T {
	return &Z_Noun_Info
}

// return a valid markup map, creating it if necessary.
func (op *Noun) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// noun, a type of flow.
const Z_Noun_Type = "noun"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Noun)(nil)

var Z_Noun_Info = typeinfo.Flow{
	Name: Z_Noun_Type,
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:  "kind",
		Label: "noun",
		Type:  &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type noun
// FIX: duplicates the spec decl.
type FIX_Noun_Slice []Noun

// implements typeinfo.Inspector
func (*Noun_Slice) Inspect() typeinfo.T {
	return &Z_Noun_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Refine struct {
	Series ScannerMaker
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Refine) Inspect() typeinfo.T {
	return &Z_Refine_Info
}

// return a valid markup map, creating it if necessary.
func (op *Refine) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// refine, a type of flow.
const Z_Refine_Type = "refine"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Refine)(nil)

var Z_Refine_Info = typeinfo.Flow{
	Name: Z_Refine_Type,
	Lede: "refine",
	Terms: []typeinfo.Term{{
		Name:    "series",
		Label:   "sequence",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type refine
// FIX: duplicates the spec decl.
type FIX_Refine_Slice []Refine

// implements typeinfo.Inspector
func (*Refine_Slice) Inspect() typeinfo.T {
	return &Z_Refine_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Reverse struct {
	Reverses ScannerMaker
	Markup   map[string]any
}

// implements typeinfo.Inspector
func (*Reverse) Inspect() typeinfo.T {
	return &Z_Reverse_Info
}

// return a valid markup map, creating it if necessary.
func (op *Reverse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// reverse, a type of flow.
const Z_Reverse_Type = "reverse"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Reverse)(nil)

var Z_Reverse_Info = typeinfo.Flow{
	Name: Z_Reverse_Type,
	Lede: "reverse",
	Terms: []typeinfo.Term{{
		Name:    "reverses",
		Label:   "_",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type reverse
// FIX: duplicates the spec decl.
type FIX_Reverse_Slice []Reverse

// implements typeinfo.Inspector
func (*Reverse_Slice) Inspect() typeinfo.T {
	return &Z_Reverse_Info
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
	return &Z_Focus_Info
}

// return a valid markup map, creating it if necessary.
func (op *Focus) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// focus, a type of flow.
const Z_Focus_Type = "focus"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Focus)(nil)

var Z_Focus_Info = typeinfo.Flow{
	Name: Z_Focus_Type,
	Lede: "focus",
	Terms: []typeinfo.Term{{
		Name:  "player",
		Label: "_",
		Type:  &prim.Z_Text_Info,
	}, {
		Name:    "series",
		Label:   "sequence",
		Repeats: true,
		Type:    &Z_ScannerMaker_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type focus
// FIX: duplicates the spec decl.
type FIX_Focus_Slice []Focus

// implements typeinfo.Inspector
func (*Focus_Slice) Inspect() typeinfo.T {
	return &Z_Focus_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Words struct {
	Words  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Words) Inspect() typeinfo.T {
	return &Z_Words_Info
}

// return a valid markup map, creating it if necessary.
func (op *Words) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// words, a type of flow.
const Z_Words_Type = "words"

// ensure the command implements its specified slots:
var _ ScannerMaker = (*Words)(nil)

var Z_Words_Info = typeinfo.Flow{
	Name: Z_Words_Type,
	Lede: "one",
	Terms: []typeinfo.Term{{
		Name:    "words",
		Label:   "word",
		Repeats: true,
		Type:    &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_ScannerMaker_Info,
	},
}

// holds a slice of type words
// FIX: duplicates the spec decl.
type FIX_Words_Slice []Words

// implements typeinfo.Inspector
func (*Words_Slice) Inspect() typeinfo.T {
	return &Z_Words_Info
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var Y_slot_List = []*typeinfo.Slot{
	&Z_GrammarMaker_Info,
	&Z_ScannerMaker_Info,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var Y_flow_List = []*typeinfo.Flow{
	&Z_Action_Info,
	&Z_Sequence_Info,
	&Z_ChooseOne_Info,
	&Z_Directive_Info,
	&Z_Noun_Info,
	&Z_Refine_Info,
	&Z_Reverse_Info,
	&Z_Focus_Info,
	&Z_Words_Info,
}

// a list of all command signatures
// ( for processing and verifying story files )
var Z_Signatures = map[uint64]interface{}{}
