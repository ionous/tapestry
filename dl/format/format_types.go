// Text printing and output control.
// The default Tapestry runtime will process printed text according to its [markup rules](https://pkg.go.dev/git.sr.ht/~ionous/tapestry/web/markup).
package format

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// counter, a type of slot.
var Zt_Counter = typeinfo.Slot{
	Name: "counter",
	Markup: map[string]any{
		"comment": "A slot used internally for generating unique names during weave.",
	},
}

// Holds a single slot.
type Counter_Slot struct{ Value Counter }

// Implements [typeinfo.Instance] for a single slot.
func (*Counter_Slot) TypeInfo() typeinfo.T {
	return &Zt_Counter
}

// Holds a slice of slots.
type Counter_Slots []Counter

// Implements [typeinfo.Instance] for a slice of slots.
func (*Counter_Slots) TypeInfo() typeinfo.T {
	return &Zt_Counter
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Counter_Slots) Repeats() bool {
	return len(*op) > 0
}

// Add a single blank line ( unless a blank line was just written ).
// See also the <p> markup.
type ParagraphBreak struct {
	Markup map[string]any
}

// paragraph_break, a type of flow.
var Zt_ParagraphBreak typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ParagraphBreak) TypeInfo() typeinfo.T {
	return &Zt_ParagraphBreak
}

// Implements [typeinfo.Markup]
func (op *ParagraphBreak) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*ParagraphBreak)(nil)

// Holds a slice of type ParagraphBreak.
type ParagraphBreak_Slice []ParagraphBreak

// Implements [typeinfo.Instance] for a slice of ParagraphBreak.
func (*ParagraphBreak_Slice) TypeInfo() typeinfo.T {
	return &Zt_ParagraphBreak
}

// Implements [typeinfo.Repeats] for a slice of ParagraphBreak.
func (op *ParagraphBreak_Slice) Repeats() bool {
	return len(*op) > 0
}

// Start a new line ( if not already at a new line ).
// See also the <wbr> markup.
type SoftBreak struct {
	Markup map[string]any
}

// soft_break, a type of flow.
var Zt_SoftBreak typeinfo.Flow

// Implements [typeinfo.Instance]
func (*SoftBreak) TypeInfo() typeinfo.T {
	return &Zt_SoftBreak
}

// Implements [typeinfo.Markup]
func (op *SoftBreak) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*SoftBreak)(nil)

// Holds a slice of type SoftBreak.
type SoftBreak_Slice []SoftBreak

// Implements [typeinfo.Instance] for a slice of SoftBreak.
func (*SoftBreak_Slice) TypeInfo() typeinfo.T {
	return &Zt_SoftBreak
}

// Implements [typeinfo.Repeats] for a slice of SoftBreak.
func (op *SoftBreak_Slice) Repeats() bool {
	return len(*op) > 0
}

// Start a new line.
// See also the <br> markup.
type LineBreak struct {
	Markup map[string]any
}

// line_break, a type of flow.
var Zt_LineBreak typeinfo.Flow

// Implements [typeinfo.Instance]
func (*LineBreak) TypeInfo() typeinfo.T {
	return &Zt_LineBreak
}

// Implements [typeinfo.Markup]
func (op *LineBreak) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*LineBreak)(nil)

// Holds a slice of type LineBreak.
type LineBreak_Slice []LineBreak

// Implements [typeinfo.Instance] for a slice of LineBreak.
func (*LineBreak_Slice) TypeInfo() typeinfo.T {
	return &Zt_LineBreak
}

// Implements [typeinfo.Repeats] for a slice of LineBreak.
func (op *LineBreak_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns some text selected from a set of predefined values. When called multiple times, this returns each one of the values in their specified order, then it loops back to the first value again.
type CycleText struct {
	Name   string
	Parts  []rtti.TextEval
	Markup map[string]any
}

// cycle_text, a type of flow.
var Zt_CycleText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*CycleText) TypeInfo() typeinfo.T {
	return &Zt_CycleText
}

// Implements [typeinfo.Markup]
func (op *CycleText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Counter = (*CycleText)(nil)
var _ rtti.TextEval = (*CycleText)(nil)

// Holds a slice of type CycleText.
type CycleText_Slice []CycleText

// Implements [typeinfo.Instance] for a slice of CycleText.
func (*CycleText_Slice) TypeInfo() typeinfo.T {
	return &Zt_CycleText
}

// Implements [typeinfo.Repeats] for a slice of CycleText.
func (op *CycleText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns some text selected from a set of predefined values. When called multiple times, this returns each one of the values in a randomized order. After returning all of the available options, it begins again with a new ordering.
type ShuffleText struct {
	Name    string
	Parts   []rtti.TextEval
	Indices Shuffler
	Markup  map[string]any
}

// shuffle_text, a type of flow.
var Zt_ShuffleText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*ShuffleText) TypeInfo() typeinfo.T {
	return &Zt_ShuffleText
}

// Implements [typeinfo.Markup]
func (op *ShuffleText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Counter = (*ShuffleText)(nil)
var _ rtti.TextEval = (*ShuffleText)(nil)

// Holds a slice of type ShuffleText.
type ShuffleText_Slice []ShuffleText

// Implements [typeinfo.Instance] for a slice of ShuffleText.
func (*ShuffleText_Slice) TypeInfo() typeinfo.T {
	return &Zt_ShuffleText
}

// Implements [typeinfo.Repeats] for a slice of ShuffleText.
func (op *ShuffleText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Returns some text selected from a set of predefined values. When called multiple times, this returns each of its inputs in turn. After returning all of the available options, it sticks to using the last option.
//
// As a special case, if there was only ever one option: it returns that option followed by nothing ( the empty string ) forever after.
type StoppingText struct {
	Name   string
	Parts  []rtti.TextEval
	Markup map[string]any
}

// stopping_text, a type of flow.
var Zt_StoppingText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*StoppingText) TypeInfo() typeinfo.T {
	return &Zt_StoppingText
}

// Implements [typeinfo.Markup]
func (op *StoppingText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Counter = (*StoppingText)(nil)
var _ rtti.TextEval = (*StoppingText)(nil)

// Holds a slice of type StoppingText.
type StoppingText_Slice []StoppingText

// Implements [typeinfo.Instance] for a slice of StoppingText.
func (*StoppingText_Slice) TypeInfo() typeinfo.T {
	return &Zt_StoppingText
}

// Implements [typeinfo.Repeats] for a slice of StoppingText.
func (op *StoppingText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Capture any and all text printed by the game, and return it as a single string of continuous text. New lines are stored as line feeds ('\n').
type BufferText struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// buffer_text, a type of flow.
var Zt_BufferText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*BufferText) TypeInfo() typeinfo.T {
	return &Zt_BufferText
}

// Implements [typeinfo.Markup]
func (op *BufferText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*BufferText)(nil)

// Holds a slice of type BufferText.
type BufferText_Slice []BufferText

// Implements [typeinfo.Instance] for a slice of BufferText.
func (*BufferText_Slice) TypeInfo() typeinfo.T {
	return &Zt_BufferText
}

// Implements [typeinfo.Repeats] for a slice of BufferText.
func (op *BufferText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Display some text to the player.
// The default runtime will format the text according to the rules specified by the Tapestry markup package:
// https://pkg.go.dev/git.sr.ht/~ionous/tapestry/web/markup
type PrintText struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// print_text, a type of flow.
var Zt_PrintText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintText) TypeInfo() typeinfo.T {
	return &Zt_PrintText
}

// Implements [typeinfo.Markup]
func (op *PrintText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*PrintText)(nil)

// Holds a slice of type PrintText.
type PrintText_Slice []PrintText

// Implements [typeinfo.Instance] for a slice of PrintText.
func (*PrintText_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintText
}

// Implements [typeinfo.Repeats] for a slice of PrintText.
func (op *PrintText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Collect printed text and separate that text by single spaces.
type PrintWords struct {
	Separator rtti.TextEval
	Exe       []rtti.Execute
	Markup    map[string]any
}

// print_words, a type of flow.
var Zt_PrintWords typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintWords) TypeInfo() typeinfo.T {
	return &Zt_PrintWords
}

// Implements [typeinfo.Markup]
func (op *PrintWords) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.Execute = (*PrintWords)(nil)
var _ rtti.TextEval = (*PrintWords)(nil)

// Holds a slice of type PrintWords.
type PrintWords_Slice []PrintWords

// Implements [typeinfo.Instance] for a slice of PrintWords.
func (*PrintWords_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintWords
}

// Implements [typeinfo.Repeats] for a slice of PrintWords.
func (op *PrintWords_Slice) Repeats() bool {
	return len(*op) > 0
}

// Collect printed text and surround the output with parenthesis '()'.
// If no text is printed, no parentheses are printed.
type PrintParens struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// print_parens, a type of flow.
var Zt_PrintParens typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintParens) TypeInfo() typeinfo.T {
	return &Zt_PrintParens
}

// Implements [typeinfo.Markup]
func (op *PrintParens) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*PrintParens)(nil)
var _ rtti.Execute = (*PrintParens)(nil)

// Holds a slice of type PrintParens.
type PrintParens_Slice []PrintParens

// Implements [typeinfo.Instance] for a slice of PrintParens.
func (*PrintParens_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintParens
}

// Implements [typeinfo.Repeats] for a slice of PrintParens.
func (op *PrintParens_Slice) Repeats() bool {
	return len(*op) > 0
}

// Separates words with commas, and 'and'.
type PrintCommas struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// print_commas, a type of flow.
var Zt_PrintCommas typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintCommas) TypeInfo() typeinfo.T {
	return &Zt_PrintCommas
}

// Implements [typeinfo.Markup]
func (op *PrintCommas) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*PrintCommas)(nil)
var _ rtti.Execute = (*PrintCommas)(nil)

// Holds a slice of type PrintCommas.
type PrintCommas_Slice []PrintCommas

// Implements [typeinfo.Instance] for a slice of PrintCommas.
func (*PrintCommas_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintCommas
}

// Implements [typeinfo.Repeats] for a slice of PrintCommas.
func (op *PrintCommas_Slice) Repeats() bool {
	return len(*op) > 0
}

// Group text into an unordered list <ul>.
type Rows struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// rows, a type of flow.
var Zt_Rows typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Rows) TypeInfo() typeinfo.T {
	return &Zt_Rows
}

// Implements [typeinfo.Markup]
func (op *Rows) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Rows)(nil)

// Holds a slice of type Rows.
type Rows_Slice []Rows

// Implements [typeinfo.Instance] for a slice of Rows.
func (*Rows_Slice) TypeInfo() typeinfo.T {
	return &Zt_Rows
}

// Implements [typeinfo.Repeats] for a slice of Rows.
func (op *Rows_Slice) Repeats() bool {
	return len(*op) > 0
}

// Group text into a single line <li> as part of a list of lines.
// See also: 'rows'.
type Row struct {
	Exe    []rtti.Execute
	Markup map[string]any
}

// row, a type of flow.
var Zt_Row typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Row) TypeInfo() typeinfo.T {
	return &Zt_Row
}

// Implements [typeinfo.Markup]
func (op *Row) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Row)(nil)

// Holds a slice of type Row.
type Row_Slice []Row

// Implements [typeinfo.Instance] for a slice of Row.
func (*Row_Slice) TypeInfo() typeinfo.T {
	return &Zt_Row
}

// Implements [typeinfo.Repeats] for a slice of Row.
func (op *Row_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_ParagraphBreak = typeinfo.Flow{
		Name:  "paragraph_break",
		Lede:  "paragraph_break",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Add a single blank line ( unless a blank line was just written ).", "See also the <p> markup."},
		},
	}
	Zt_SoftBreak = typeinfo.Flow{
		Name:  "soft_break",
		Lede:  "soft_break",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Start a new line ( if not already at a new line ).", "See also the <wbr> markup."},
		},
	}
	Zt_LineBreak = typeinfo.Flow{
		Name:  "line_break",
		Lede:  "line_break",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Start a new line.", "See also the <br> markup."},
		},
	}
	Zt_CycleText = typeinfo.Flow{
		Name: "cycle_text",
		Lede: "cycle",
		Terms: []typeinfo.Term{{
			Name:     "name",
			Label:    "name",
			Optional: true,
			Markup: map[string]any{
				"comment": "An optional name used for controlling internal state.  When omitted, weave automatically generates a globally unique name. Commands with the same name will share internal state.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "parts",
			Label:   "text",
			Repeats: true,
			Markup: map[string]any{
				"comment": "One or more pieces of text to cycle through.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Counter,
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Returns some text selected from a set of predefined values. When called multiple times, this returns each one of the values in their specified order, then it loops back to the first value again.",
		},
	}
	Zt_ShuffleText = typeinfo.Flow{
		Name: "shuffle_text",
		Lede: "shuffle",
		Terms: []typeinfo.Term{{
			Name:     "name",
			Label:    "name",
			Optional: true,
			Markup: map[string]any{
				"comment": "An optional name used for controlling internal state.  When omitted, weave automatically generates a globally unique name. Commands with the same name will share internal state.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "parts",
			Label:   "text",
			Repeats: true,
			Markup: map[string]any{
				"comment": "One or more pieces of text to shuffle through.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:    "indices",
			Label:   "indices",
			Private: true,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Counter,
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Returns some text selected from a set of predefined values. When called multiple times, this returns each one of the values in a randomized order. After returning all of the available options, it begins again with a new ordering.",
		},
	}
	Zt_StoppingText = typeinfo.Flow{
		Name: "stopping_text",
		Lede: "stopping",
		Terms: []typeinfo.Term{{
			Name:     "name",
			Label:    "name",
			Optional: true,
			Markup: map[string]any{
				"comment": "An optional name used for controlling internal state. When omitted, weave automatically generates a globally unique name. Commands with the same name will share internal state.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "parts",
			Label:   "text",
			Repeats: true,
			Markup: map[string]any{
				"comment": "One or more pieces of text to shift through.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Counter,
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Returns some text selected from a set of predefined values. When called multiple times, this returns each of its inputs in turn. After returning all of the available options, it sticks to using the last option.", "", "As a special case, if there was only ever one option: it returns that option followed by nothing ( the empty string ) forever after."},
		},
	}
	Zt_BufferText = typeinfo.Flow{
		Name: "buffer_text",
		Lede: "buffer",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"comment": "The statements to capture text output from.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Capture any and all text printed by the game, and return it as a single string of continuous text. New lines are stored as line feeds ('\\n').",
		},
	}
	Zt_PrintText = typeinfo.Flow{
		Name: "print_text",
		Lede: "print",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to print.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Display some text to the player.", "The default runtime will format the text according to the rules specified by the Tapestry markup package:", "https://pkg.go.dev/git.sr.ht/~ionous/tapestry/web/markup"},
		},
	}
	Zt_PrintWords = typeinfo.Flow{
		Name: "print_words",
		Lede: "print",
		Terms: []typeinfo.Term{{
			Name:     "separator",
			Label:    "separator",
			Optional: true,
			Markup: map[string]any{
				"comment": "Optional text to place between adjoining words.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:    "exe",
			Label:   "words",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Runs one or more statements, and collects any text printed by them.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Collect printed text and separate that text by single spaces.",
		},
	}
	Zt_PrintParens = typeinfo.Flow{
		Name: "print_parens",
		Lede: "print",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "parentheses",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Runs one or more statements, and collects any text printed by them.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Collect printed text and surround the output with parenthesis '()'.", "If no text is printed, no parentheses are printed."},
		},
	}
	Zt_PrintCommas = typeinfo.Flow{
		Name: "print_commas",
		Lede: "print",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "commas",
			Repeats: true,
			Type:    &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": "Separates words with commas, and 'and'.",
		},
	}
	Zt_Rows = typeinfo.Flow{
		Name: "rows",
		Lede: "rows",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Runs one or more statements, and collects any text printed by them.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Group text into an unordered list <ul>.",
		},
	}
	Zt_Row = typeinfo.Flow{
		Name: "row",
		Lede: "row",
		Terms: []typeinfo.Term{{
			Name:    "exe",
			Label:   "do",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Runs one or more statements, and collects any text printed by them.",
			},
			Type: &rtti.Zt_Execute,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Group text into a single line <li> as part of a list of lines.", "See also: 'rows'."},
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "format",
	Comment: []string{
		"Text printing and output control.",
		"The default Tapestry runtime will process printed text according to its [markup rules](https://pkg.go.dev/git.sr.ht/~ionous/tapestry/web/markup).",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Counter,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_ParagraphBreak,
	&Zt_SoftBreak,
	&Zt_LineBreak,
	&Zt_CycleText,
	&Zt_ShuffleText,
	&Zt_StoppingText,
	&Zt_BufferText,
	&Zt_PrintText,
	&Zt_PrintWords,
	&Zt_PrintParens,
	&Zt_PrintCommas,
	&Zt_Rows,
	&Zt_Row,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	9767668117811810575:  (*BufferText)(nil),     /* text_eval=Buffer do: */
	16098131496381194958: (*CycleText)(nil),      /* counter=Cycle name:text: */
	5355971188045229340:  (*CycleText)(nil),      /* text_eval=Cycle name:text: */
	17596073119249480739: (*CycleText)(nil),      /* counter=Cycle text: */
	1237120803959249173:  (*CycleText)(nil),      /* text_eval=Cycle text: */
	10898429598193857104: (*LineBreak)(nil),      /* execute=LineBreak */
	1194153657675604478:  (*ParagraphBreak)(nil), /* execute=ParagraphBreak */
	16169738297367022876: (*PrintCommas)(nil),    /* execute=Print commas: */
	6231219704730380469:  (*PrintCommas)(nil),    /* text_eval=Print commas: */
	4206645811149297220:  (*PrintParens)(nil),    /* execute=Print parentheses: */
	8219093008786386363:  (*PrintParens)(nil),    /* text_eval=Print parentheses: */
	4149419216708670664:  (*PrintWords)(nil),     /* execute=Print separator:words: */
	4219359027975954467:  (*PrintWords)(nil),     /* text_eval=Print separator:words: */
	1331651249232124175:  (*PrintWords)(nil),     /* execute=Print words: */
	17978150574109115948: (*PrintWords)(nil),     /* text_eval=Print words: */
	4512128922644282356:  (*PrintText)(nil),      /* execute=Print: */
	14820902432990466009: (*Row)(nil),            /* text_eval=Row do: */
	42419598592471524:    (*Rows)(nil),           /* text_eval=Rows do: */
	12460624099586212271: (*ShuffleText)(nil),    /* counter=Shuffle name:text: */
	8909818107999898193:  (*ShuffleText)(nil),    /* text_eval=Shuffle name:text: */
	3444877746271964624:  (*ShuffleText)(nil),    /* counter=Shuffle text: */
	7835310741853066190:  (*ShuffleText)(nil),    /* text_eval=Shuffle text: */
	17335248920749226950: (*SoftBreak)(nil),      /* execute=SoftBreak */
	13115056552370612412: (*StoppingText)(nil),   /* counter=Stopping name:text: */
	11830555676954637550: (*StoppingText)(nil),   /* text_eval=Stopping name:text: */
	13363393271236249653: (*StoppingText)(nil),   /* counter=Stopping text: */
	9145628730349656131:  (*StoppingText)(nil),   /* text_eval=Stopping text: */
}
