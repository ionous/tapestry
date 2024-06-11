// Player input parsing.
package grammar

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// scanner_maker, a type of slot.
var Zt_ScannerMaker = typeinfo.Slot{
	Name: "scanner_maker",
	Markup: map[string]any{
		"comment": "Commands which interpret player input.",
	},
}

// Holds a single slot.
type ScannerMaker_Slot struct{ Value ScannerMaker }

// Implements [typeinfo.Instance] for a single slot.
func (*ScannerMaker_Slot) TypeInfo() typeinfo.T {
	return &Zt_ScannerMaker
}

// Holds a slice of slots.
type ScannerMaker_Slots []ScannerMaker

// Implements [typeinfo.Instance] for a slice of slots.
func (*ScannerMaker_Slots) TypeInfo() typeinfo.T {
	return &Zt_ScannerMaker
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *ScannerMaker_Slots) Repeats() bool {
	return len(*op) > 0
}

// Starts a parser scanner.
// This is generated by story statement "DefineNamedGrammar"
// so that grammar doesn't have to import weave.
type Directive struct {
	Name   string
	Series []ScannerMaker
	Markup map[string]any
}

// directive, a type of flow.
var Zt_Directive typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Directive) TypeInfo() typeinfo.T {
	return &Zt_Directive
}

// Implements [typeinfo.Markup]
func (op *Directive) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Holds a slice of type Directive.
type Directive_Slice []Directive

// Implements [typeinfo.Instance] for a slice of Directive.
func (*Directive_Slice) TypeInfo() typeinfo.T {
	return &Zt_Directive
}

// Implements [typeinfo.Repeats] for a slice of Directive.
func (op *Directive_Slice) Repeats() bool {
	return len(*op) > 0
}

// Run a script defined action.
// The action will be parameterized with the actor performing the action,
// and up to two other nouns. The actor is always the player's,
// the nouns are whichever were matched by [Noun].
// This command usually appears last in a series of sub-scanners,
// after all of the earlier scanners have successfully matched.
type Action struct {
	Action    string
	Arguments []call.Arg
	Markup    map[string]any
}

// action, a type of flow.
var Zt_Action typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Action) TypeInfo() typeinfo.T {
	return &Zt_Action
}

// Implements [typeinfo.Markup]
func (op *Action) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Action)(nil)

// Holds a slice of type Action.
type Action_Slice []Action

// Implements [typeinfo.Instance] for a slice of Action.
func (*Action_Slice) TypeInfo() typeinfo.T {
	return &Zt_Action
}

// Implements [typeinfo.Repeats] for a slice of Action.
func (op *Action_Slice) Repeats() bool {
	return len(*op) > 0
}

// Require that all of its sub-scanners match.
type Sequence struct {
	Series []ScannerMaker
	Markup map[string]any
}

// sequence, a type of flow.
var Zt_Sequence typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Sequence) TypeInfo() typeinfo.T {
	return &Zt_Sequence
}

// Implements [typeinfo.Markup]
func (op *Sequence) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Sequence)(nil)

// Holds a slice of type Sequence.
type Sequence_Slice []Sequence

// Implements [typeinfo.Instance] for a slice of Sequence.
func (*Sequence_Slice) TypeInfo() typeinfo.T {
	return &Zt_Sequence
}

// Implements [typeinfo.Repeats] for a slice of Sequence.
func (op *Sequence_Slice) Repeats() bool {
	return len(*op) > 0
}

// Attempts to match exactly one of its sub-scanners.
// Its sub-scanners are evaluated in the order listed;
// stopping after the first successful match has occurred.
type ChooseOne struct {
	Options []ScannerMaker
	Markup  map[string]any
}

// choose_one, a type of flow.
var Zt_ChooseOne typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ChooseOne) TypeInfo() typeinfo.T {
	return &Zt_ChooseOne
}

// Implements [typeinfo.Markup]
func (op *ChooseOne) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*ChooseOne)(nil)

// Holds a slice of type ChooseOne.
type ChooseOne_Slice []ChooseOne

// Implements [typeinfo.Instance] for a slice of ChooseOne.
func (*ChooseOne_Slice) TypeInfo() typeinfo.T {
	return &Zt_ChooseOne
}

// Implements [typeinfo.Repeats] for a slice of ChooseOne.
func (op *ChooseOne_Slice) Repeats() bool {
	return len(*op) > 0
}

// Attempt to find a noun matching one or more words
// entered by the player. Out of all the nouns in a game world
// the specific nouns matched will, by default, depend on what
// objects are visible to the player.
type Noun struct {
	Kind   string
	Markup map[string]any
}

// noun, a type of flow.
var Zt_Noun typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Noun) TypeInfo() typeinfo.T {
	return &Zt_Noun
}

// Implements [typeinfo.Markup]
func (op *Noun) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Noun)(nil)

// Holds a slice of type Noun.
type Noun_Slice []Noun

// Implements [typeinfo.Instance] for a slice of Noun.
func (*Noun_Slice) TypeInfo() typeinfo.T {
	return &Zt_Noun
}

// Implements [typeinfo.Repeats] for a slice of Noun.
func (op *Noun_Slice) Repeats() bool {
	return len(*op) > 0
}

// Uses the last noun in a series of sub-scanners as a source
// for finding for the first noun in the series.
// For example, in the phrase "take book from table"
// if there was one book on a table, and many books on a bookshelf,
// refine would ensure that the book on the table was selected.
type Refine struct {
	Series []ScannerMaker
	Markup map[string]any
}

// refine, a type of flow.
var Zt_Refine typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Refine) TypeInfo() typeinfo.T {
	return &Zt_Refine
}

// Implements [typeinfo.Markup]
func (op *Refine) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Refine)(nil)

// Holds a slice of type Refine.
type Refine_Slice []Refine

// Implements [typeinfo.Instance] for a slice of Refine.
func (*Refine_Slice) TypeInfo() typeinfo.T {
	return &Zt_Refine
}

// Implements [typeinfo.Repeats] for a slice of Refine.
func (op *Refine_Slice) Repeats() bool {
	return len(*op) > 0
}

// Assuming that the first and last sub-scanners in a set
// match nouns, reverse the order of those nouns when sending
// them to an action.
type Reverse struct {
	Reverses []ScannerMaker
	Markup   map[string]any
}

// reverse, a type of flow.
var Zt_Reverse typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Reverse) TypeInfo() typeinfo.T {
	return &Zt_Reverse
}

// Implements [typeinfo.Markup]
func (op *Reverse) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Reverse)(nil)

// Holds a slice of type Reverse.
type Reverse_Slice []Reverse

// Implements [typeinfo.Instance] for a slice of Reverse.
func (*Reverse_Slice) TypeInfo() typeinfo.T {
	return &Zt_Reverse
}

// Implements [typeinfo.Repeats] for a slice of Reverse.
func (op *Reverse_Slice) Repeats() bool {
	return len(*op) > 0
}

// Select a specific set of bounds for the scanner to use when matching nouns.
// Currently, the set of scanners is defined in golang,
// and cannot be extended by script.
// They are:
//   - "" - the empty string for the current locale ( the default bounds )
//   - "compass" - for nouns of kind "directions" ( north, east, etc. )
//   - "player" - for the player's inventory
//   - "debugging" - for all objects of kind "objects"
//
// see: https://pkg.go.dev/git.sr.ht/~ionous/tapestry/support/play#Survey
type Focus struct {
	Player string
	Series []ScannerMaker
	Markup map[string]any
}

// focus, a type of flow.
var Zt_Focus typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Focus) TypeInfo() typeinfo.T {
	return &Zt_Focus
}

// Implements [typeinfo.Markup]
func (op *Focus) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Focus)(nil)

// Holds a slice of type Focus.
type Focus_Slice []Focus

// Implements [typeinfo.Instance] for a slice of Focus.
func (*Focus_Slice) TypeInfo() typeinfo.T {
	return &Zt_Focus
}

// Implements [typeinfo.Repeats] for a slice of Focus.
func (op *Focus_Slice) Repeats() bool {
	return len(*op) > 0
}

// Match text exactly as typed by the player.
// To match two words that must appear together
// use a [Sequence] containing two separate word commands.
// For example, to match the phrase "take off", use:
//
//	Sequence:
//	  - One word: "take"
//	  - One word: "off"
type Words struct {
	Words  []string
	Markup map[string]any
}

// words, a type of flow.
var Zt_Words typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Words) TypeInfo() typeinfo.T {
	return &Zt_Words
}

// Implements [typeinfo.Markup]
func (op *Words) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ ScannerMaker = (*Words)(nil)

// Holds a slice of type Words.
type Words_Slice []Words

// Implements [typeinfo.Instance] for a slice of Words.
func (*Words_Slice) TypeInfo() typeinfo.T {
	return &Zt_Words
}

// Implements [typeinfo.Repeats] for a slice of Words.
func (op *Words_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_Directive = typeinfo.Flow{
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
		Markup: map[string]any{
			"comment":  []interface{}{"Starts a parser scanner.", "This is generated by story statement \"DefineNamedGrammar\"", "so that grammar doesn't have to import weave."},
			"internal": true,
		},
	}
	Zt_Action = typeinfo.Flow{
		Name: "action",
		Lede: "action",
		Terms: []typeinfo.Term{{
			Name: "action",
			Markup: map[string]any{
				"comment": []interface{}{"Name of the action.", "See also: [DefineAction]."},
			},
			Type: &prim.Zt_Text,
		}, {
			Name:     "arguments",
			Label:    "args",
			Optional: true,
			Repeats:  true,
			Markup: map[string]any{
				"comment": "Optional additional arguments to pass to the action.",
			},
			Type: &call.Zt_Arg,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Run a script defined action.", "The action will be parameterized with the actor performing the action,", "and up to two other nouns. The actor is always the player's,", "the nouns are whichever were matched by [Noun].", "This command usually appears last in a series of sub-scanners,", "after all of the earlier scanners have successfully matched."},
		},
	}
	Zt_Sequence = typeinfo.Flow{
		Name: "sequence",
		Lede: "sequence",
		Terms: []typeinfo.Term{{
			Name:    "series",
			Repeats: true,
			Markup: map[string]any{
				"comment": "A series of other scanners to match.",
			},
			Type: &Zt_ScannerMaker,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": "Require that all of its sub-scanners match.",
		},
	}
	Zt_ChooseOne = typeinfo.Flow{
		Name: "choose_one",
		Lede: "one",
		Terms: []typeinfo.Term{{
			Name:    "options",
			Label:   "of",
			Repeats: true,
			Markup: map[string]any{
				"comment": "A series of other scanners to match.",
			},
			Type: &Zt_ScannerMaker,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Attempts to match exactly one of its sub-scanners.", "Its sub-scanners are evaluated in the order listed;", "stopping after the first successful match has occurred."},
		},
	}
	Zt_Noun = typeinfo.Flow{
		Name: "noun",
		Lede: "one",
		Terms: []typeinfo.Term{{
			Name:  "kind",
			Label: "noun",
			Markup: map[string]any{
				"comment": []interface{}{"The kind of the noun that should be allowed to match.", "Most often this is the most generic kind: \"objects\".", "That allows the action to respond to the player", "when an unexpected or undesired noun was entered."},
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Attempt to find a noun matching one or more words", "entered by the player. Out of all the nouns in a game world", "the specific nouns matched will, by default, depend on what", "objects are visible to the player."},
		},
	}
	Zt_Refine = typeinfo.Flow{
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
			"comment": []interface{}{"Uses the last noun in a series of sub-scanners as a source", "for finding for the first noun in the series.", "For example, in the phrase \"take book from table\"", "if there was one book on a table, and many books on a bookshelf,", "refine would ensure that the book on the table was selected."},
		},
	}
	Zt_Reverse = typeinfo.Flow{
		Name: "reverse",
		Lede: "reverse",
		Terms: []typeinfo.Term{{
			Name:    "reverses",
			Repeats: true,
			Type:    &Zt_ScannerMaker,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Assuming that the first and last sub-scanners in a set", "match nouns, reverse the order of those nouns when sending", "them to an action."},
		},
	}
	Zt_Focus = typeinfo.Flow{
		Name: "focus",
		Lede: "focus",
		Terms: []typeinfo.Term{{
			Name: "player",
			Type: &prim.Zt_Text,
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
			"comment": []interface{}{"Select a specific set of bounds for the scanner to use when matching nouns.", "Currently, the set of scanners is defined in golang,", "and cannot be extended by script.", "They are:", "  * \"\" - the empty string for the current locale ( the default bounds )", "  * \"compass\" - for nouns of kind \"directions\" ( north, east, etc. )", "  * \"player\" - for the player's inventory", "  * \"debugging\" - for all objects of kind \"objects\"", "see: https://pkg.go.dev/git.sr.ht/~ionous/tapestry/support/play#Survey"},
		},
	}
	Zt_Words = typeinfo.Flow{
		Name: "words",
		Lede: "one",
		Terms: []typeinfo.Term{{
			Name:    "words",
			Label:   "word",
			Repeats: true,
			Markup: map[string]any{
				"comment": []interface{}{"One or more synonyms to match.", "Only the first word matched will be used.", "For example: \"pluck\", \"pick\", \"strum\"", "would match either \"pluck the guitar\",", "or \"strum the guitar\"."},
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_ScannerMaker,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Match text exactly as typed by the player.", "To match two words that must appear together", "use a [Sequence] containing two separate word commands.", "For example, to match the phrase \"take off\", use:", "  Sequence:", "    - One word: \"take\"", "    - One word: \"off\""},
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "grammar",
	Comment: []string{
		"Player input parsing.",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_ScannerMaker,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Directive,
	&Zt_Action,
	&Zt_Sequence,
	&Zt_ChooseOne,
	&Zt_Noun,
	&Zt_Refine,
	&Zt_Reverse,
	&Zt_Focus,
	&Zt_Words,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	17030018957559107353: (*Directive)(nil), /* Interpret name:with: */
	12048905879374467271: (*Action)(nil),    /* scanner_maker=Action: */
	967998274944030280:   (*Action)(nil),    /* scanner_maker=Action:args: */
	1756442538083378424:  (*Focus)(nil),     /* scanner_maker=Focus:sequence: */
	10964817074887037945: (*Noun)(nil),      /* scanner_maker=One noun: */
	16418039705711067622: (*ChooseOne)(nil), /* scanner_maker=One of: */
	16180319172078511701: (*Words)(nil),     /* scanner_maker=One word: */
	11402479949132197621: (*Refine)(nil),    /* scanner_maker=Refine sequence: */
	15857934419606450901: (*Reverse)(nil),   /* scanner_maker=Reverse: */
	10728359537834940094: (*Sequence)(nil),  /* scanner_maker=Sequence: */
}
