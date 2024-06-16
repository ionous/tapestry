// Text manipulation and transformation.
package text

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// Count the number of characters in some text.
type TextLen struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// text_len, a type of flow.
var Zt_TextLen typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextLen) TypeInfo() typeinfo.T {
	return &Zt_TextLen
}

// Implements [typeinfo.Markup]
func (op *TextLen) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.NumEval = (*TextLen)(nil)

// Holds a slice of type TextLen.
type TextLen_Slice []TextLen

// Implements [typeinfo.Instance] for a slice of TextLen.
func (*TextLen_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextLen
}

// Implements [typeinfo.Repeats] for a slice of TextLen.
func (op *TextLen_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether one piece of text contains a second piece of text.
//
// The [rt.NumEval] version returns the first index at which the text appears,
// or zero if not found.
type FindText struct {
	Text    rtti.TextEval
	Subtext rtti.TextEval
	Markup  map[string]any
}

// find_text, a type of flow.
var Zt_FindText typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FindText) TypeInfo() typeinfo.T {
	return &Zt_FindText
}

// Implements [typeinfo.Markup]
func (op *FindText) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*FindText)(nil)
var _ rtti.NumEval = (*FindText)(nil)

// Holds a slice of type FindText.
type FindText_Slice []FindText

// Implements [typeinfo.Instance] for a slice of FindText.
func (*FindText_Slice) TypeInfo() typeinfo.T {
	return &Zt_FindText
}

// Implements [typeinfo.Repeats] for a slice of FindText.
func (op *FindText_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether text starts in a particular way.
type TextStartsWith struct {
	Text    rtti.TextEval
	Subtext rtti.TextEval
	Markup  map[string]any
}

// text_starts_with, a type of flow.
var Zt_TextStartsWith typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextStartsWith) TypeInfo() typeinfo.T {
	return &Zt_TextStartsWith
}

// Implements [typeinfo.Markup]
func (op *TextStartsWith) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*TextStartsWith)(nil)

// Holds a slice of type TextStartsWith.
type TextStartsWith_Slice []TextStartsWith

// Implements [typeinfo.Instance] for a slice of TextStartsWith.
func (*TextStartsWith_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextStartsWith
}

// Implements [typeinfo.Repeats] for a slice of TextStartsWith.
func (op *TextStartsWith_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether text ends in a particular way.
type TextEndsWith struct {
	Text    rtti.TextEval
	Subtext rtti.TextEval
	Markup  map[string]any
}

// text_ends_with, a type of flow.
var Zt_TextEndsWith typeinfo.Flow

// Implements [typeinfo.Instance]
func (*TextEndsWith) TypeInfo() typeinfo.T {
	return &Zt_TextEndsWith
}

// Implements [typeinfo.Markup]
func (op *TextEndsWith) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*TextEndsWith)(nil)

// Holds a slice of type TextEndsWith.
type TextEndsWith_Slice []TextEndsWith

// Implements [typeinfo.Instance] for a slice of TextEndsWith.
func (*TextEndsWith_Slice) TypeInfo() typeinfo.T {
	return &Zt_TextEndsWith
}

// Implements [typeinfo.Repeats] for a slice of TextEndsWith.
func (op *TextEndsWith_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether text is completely without content.
// Even spaces are considered content. The text "" is considered empty,
// the text " " is considered *not* empty.
type IsEmpty struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// is_empty, a type of flow.
var Zt_IsEmpty typeinfo.Flow

// Implements [typeinfo.Instance]
func (*IsEmpty) TypeInfo() typeinfo.T {
	return &Zt_IsEmpty
}

// Implements [typeinfo.Markup]
func (op *IsEmpty) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*IsEmpty)(nil)

// Holds a slice of type IsEmpty.
type IsEmpty_Slice []IsEmpty

// Implements [typeinfo.Instance] for a slice of IsEmpty.
func (*IsEmpty_Slice) TypeInfo() typeinfo.T {
	return &Zt_IsEmpty
}

// Implements [typeinfo.Repeats] for a slice of IsEmpty.
func (op *IsEmpty_Slice) Repeats() bool {
	return len(*op) > 0
}

// Determine whether text matches a regular expression.
// The expressions used are defined by go.
// https://pkg.go.dev/regexp/syntax
// https://github.com/google/re2/wiki/Syntax
type Matches struct {
	Text   rtti.TextEval
	Match  string
	Cache  MatchCache
	Markup map[string]any
}

// matches, a type of flow.
var Zt_Matches typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Matches) TypeInfo() typeinfo.T {
	return &Zt_Matches
}

// Implements [typeinfo.Markup]
func (op *Matches) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.BoolEval = (*Matches)(nil)

// Holds a slice of type Matches.
type Matches_Slice []Matches

// Implements [typeinfo.Instance] for a slice of Matches.
func (*Matches_Slice) TypeInfo() typeinfo.T {
	return &Zt_Matches
}

// Implements [typeinfo.Repeats] for a slice of Matches.
func (op *Matches_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy some text, changing its first letter to uppercase.
type Capitalize struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// capitalize, a type of flow.
var Zt_Capitalize typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Capitalize) TypeInfo() typeinfo.T {
	return &Zt_Capitalize
}

// Implements [typeinfo.Markup]
func (op *Capitalize) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Capitalize)(nil)

// Holds a slice of type Capitalize.
type Capitalize_Slice []Capitalize

// Implements [typeinfo.Instance] for a slice of Capitalize.
func (*Capitalize_Slice) TypeInfo() typeinfo.T {
	return &Zt_Capitalize
}

// Implements [typeinfo.Repeats] for a slice of Capitalize.
func (op *Capitalize_Slice) Repeats() bool {
	return len(*op) > 0
}

// Combine text to produce new text.
type Join struct {
	Sep    rtti.TextEval
	Parts  []rtti.TextEval
	Markup map[string]any
}

// join, a type of flow.
var Zt_Join typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Join) TypeInfo() typeinfo.T {
	return &Zt_Join
}

// Implements [typeinfo.Markup]
func (op *Join) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Join)(nil)

// Holds a slice of type Join.
type Join_Slice []Join

// Implements [typeinfo.Instance] for a slice of Join.
func (*Join_Slice) TypeInfo() typeinfo.T {
	return &Zt_Join
}

// Implements [typeinfo.Repeats] for a slice of Join.
func (op *Join_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy some text, changing every letter into lowercase.
// For example, turns "QUIET" into "quiet.
type MakeLowercase struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// make_lowercase, a type of flow.
var Zt_MakeLowercase typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeLowercase) TypeInfo() typeinfo.T {
	return &Zt_MakeLowercase
}

// Implements [typeinfo.Markup]
func (op *MakeLowercase) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*MakeLowercase)(nil)

// Holds a slice of type MakeLowercase.
type MakeLowercase_Slice []MakeLowercase

// Implements [typeinfo.Instance] for a slice of MakeLowercase.
func (*MakeLowercase_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeLowercase
}

// Implements [typeinfo.Repeats] for a slice of MakeLowercase.
func (op *MakeLowercase_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy some text with its contents flipped back to front.
// For example, turns "Tapestry" into 'yrtsepaT'.
type MakeReversed struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// make_reversed, a type of flow.
var Zt_MakeReversed typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeReversed) TypeInfo() typeinfo.T {
	return &Zt_MakeReversed
}

// Implements [typeinfo.Markup]
func (op *MakeReversed) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*MakeReversed)(nil)

// Holds a slice of type MakeReversed.
type MakeReversed_Slice []MakeReversed

// Implements [typeinfo.Instance] for a slice of MakeReversed.
func (*MakeReversed_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeReversed
}

// Implements [typeinfo.Repeats] for a slice of MakeReversed.
func (op *MakeReversed_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy text, changing the start of each sentence so that it starts with a capital letter. ( Currently, "sentences" are considered to be a series of characters ending with a full-stop followed by a space. )
// For example, "see the doctor run. run doctor. run." into "See the doctor run. Run doctor. Run."
type MakeSentenceCase struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// make_sentence_case, a type of flow.
var Zt_MakeSentenceCase typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeSentenceCase) TypeInfo() typeinfo.T {
	return &Zt_MakeSentenceCase
}

// Implements [typeinfo.Markup]
func (op *MakeSentenceCase) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*MakeSentenceCase)(nil)

// Holds a slice of type MakeSentenceCase.
type MakeSentenceCase_Slice []MakeSentenceCase

// Implements [typeinfo.Instance] for a slice of MakeSentenceCase.
func (*MakeSentenceCase_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeSentenceCase
}

// Implements [typeinfo.Repeats] for a slice of MakeSentenceCase.
func (op *MakeSentenceCase_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy some text, making every word start with a capital letter.
// For example, turns "empire apple" into "Empire Apple".
type MakeTitleCase struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// make_title_case, a type of flow.
var Zt_MakeTitleCase typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeTitleCase) TypeInfo() typeinfo.T {
	return &Zt_MakeTitleCase
}

// Implements [typeinfo.Markup]
func (op *MakeTitleCase) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*MakeTitleCase)(nil)

// Holds a slice of type MakeTitleCase.
type MakeTitleCase_Slice []MakeTitleCase

// Implements [typeinfo.Instance] for a slice of MakeTitleCase.
func (*MakeTitleCase_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeTitleCase
}

// Implements [typeinfo.Repeats] for a slice of MakeTitleCase.
func (op *MakeTitleCase_Slice) Repeats() bool {
	return len(*op) > 0
}

// Copy some text, changing every letter into uppercase.
// For example, transforms "loud" into "LOUD".
type MakeUppercase struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// make_uppercase, a type of flow.
var Zt_MakeUppercase typeinfo.Flow

// Implements [typeinfo.Instance]
func (*MakeUppercase) TypeInfo() typeinfo.T {
	return &Zt_MakeUppercase
}

// Implements [typeinfo.Markup]
func (op *MakeUppercase) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*MakeUppercase)(nil)

// Holds a slice of type MakeUppercase.
type MakeUppercase_Slice []MakeUppercase

// Implements [typeinfo.Instance] for a slice of MakeUppercase.
func (*MakeUppercase_Slice) TypeInfo() typeinfo.T {
	return &Zt_MakeUppercase
}

// Implements [typeinfo.Repeats] for a slice of MakeUppercase.
func (op *MakeUppercase_Slice) Repeats() bool {
	return len(*op) > 0
}

// Pluralize a word.
// The singular form of a word can have more than one plural form.
// For example: "person" can be "people" or "persons".
// If more than one exists, this chooses arbitrarily.
//
// Note, The transformation uses predefined rules and some explicit mappings.
// The story command [DefinePlural] can add new mappings.
type Pluralize struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// pluralize, a type of flow.
var Zt_Pluralize typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Pluralize) TypeInfo() typeinfo.T {
	return &Zt_Pluralize
}

// Implements [typeinfo.Markup]
func (op *Pluralize) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Pluralize)(nil)

// Holds a slice of type Pluralize.
type Pluralize_Slice []Pluralize

// Implements [typeinfo.Instance] for a slice of Pluralize.
func (*Pluralize_Slice) TypeInfo() typeinfo.T {
	return &Zt_Pluralize
}

// Implements [typeinfo.Repeats] for a slice of Pluralize.
func (op *Pluralize_Slice) Repeats() bool {
	return len(*op) > 0
}

// Change a plural word into its singular form.
// A plural word only has one singular form.
// For example, given the word "people", return "person".
// See [pluralize] for more information.
type Singularize struct {
	Text   rtti.TextEval
	Markup map[string]any
}

// singularize, a type of flow.
var Zt_Singularize typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Singularize) TypeInfo() typeinfo.T {
	return &Zt_Singularize
}

// Implements [typeinfo.Markup]
func (op *Singularize) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*Singularize)(nil)

// Holds a slice of type Singularize.
type Singularize_Slice []Singularize

// Implements [typeinfo.Instance] for a slice of Singularize.
func (*Singularize_Slice) TypeInfo() typeinfo.T {
	return &Zt_Singularize
}

// Implements [typeinfo.Repeats] for a slice of Singularize.
func (op *Singularize_Slice) Repeats() bool {
	return len(*op) > 0
}

// Express a number using numerals.
// For example, given the number `1` return the text "1".
//
// The [story.Execute] version prints the text for the player.
type PrintNumDigits struct {
	Num    rtti.NumEval
	Markup map[string]any
}

// print_num_digits, a type of flow.
var Zt_PrintNumDigits typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintNumDigits) TypeInfo() typeinfo.T {
	return &Zt_PrintNumDigits
}

// Implements [typeinfo.Markup]
func (op *PrintNumDigits) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*PrintNumDigits)(nil)
var _ rtti.Execute = (*PrintNumDigits)(nil)

// Holds a slice of type PrintNumDigits.
type PrintNumDigits_Slice []PrintNumDigits

// Implements [typeinfo.Instance] for a slice of PrintNumDigits.
func (*PrintNumDigits_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintNumDigits
}

// Implements [typeinfo.Repeats] for a slice of PrintNumDigits.
func (op *PrintNumDigits_Slice) Repeats() bool {
	return len(*op) > 0
}

// Express an integer in plain english.
// For example, given the number `1` return the text "one".
// It converts floating point numbers to integer by truncating:
// given `1.6`, it returns "one".
//
// The [story.Execute] version prints the text for the player.
type PrintNumWords struct {
	Num    rtti.NumEval
	Markup map[string]any
}

// print_num_words, a type of flow.
var Zt_PrintNumWords typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PrintNumWords) TypeInfo() typeinfo.T {
	return &Zt_PrintNumWords
}

// Implements [typeinfo.Markup]
func (op *PrintNumWords) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ rtti.TextEval = (*PrintNumWords)(nil)
var _ rtti.Execute = (*PrintNumWords)(nil)

// Holds a slice of type PrintNumWords.
type PrintNumWords_Slice []PrintNumWords

// Implements [typeinfo.Instance] for a slice of PrintNumWords.
func (*PrintNumWords_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintNumWords
}

// Implements [typeinfo.Repeats] for a slice of PrintNumWords.
func (op *PrintNumWords_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_TextLen = typeinfo.Flow{
		Name: "text_len",
		Lede: "text",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "length",
			Markup: map[string]any{
				"comment": "The text to measure.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_NumEval,
		},
		Markup: map[string]any{
			"comment": "Count the number of characters in some text.",
		},
	}
	Zt_FindText = typeinfo.Flow{
		Name: "find_text",
		Lede: "find",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to search within.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "subtext",
			Label: "text",
			Markup: map[string]any{
				"comment": "The text to find.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
			&rtti.Zt_NumEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Determine whether one piece of text contains a second piece of text.", "", "The [rt.NumEval] version returns the first index at which the text appears,", "or zero if not found."},
		},
	}
	Zt_TextStartsWith = typeinfo.Flow{
		Name: "text_starts_with",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "text",
			Markup: map[string]any{
				"comment": "The text to search within.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "subtext",
			Label: "prefix",
			Markup: map[string]any{
				"comment": "The text to find.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": "Determine whether text starts in a particular way.",
		},
	}
	Zt_TextEndsWith = typeinfo.Flow{
		Name: "text_ends_with",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "text",
			Markup: map[string]any{
				"comment": "The text to search within.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "subtext",
			Label: "suffix",
			Markup: map[string]any{
				"comment": "The text to find.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": "Determine whether text ends in a particular way.",
		},
	}
	Zt_IsEmpty = typeinfo.Flow{
		Name: "is_empty",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "empty",
			Markup: map[string]any{
				"comment": "The text to check for content.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Determine whether text is completely without content.", "Even spaces are considered content. The text \"\" is considered empty,", "the text \" \" is considered *not* empty."},
		},
	}
	Zt_Matches = typeinfo.Flow{
		Name: "matches",
		Lede: "is",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "text",
			Markup: map[string]any{
				"comment": "The text to match the expression against.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:  "match",
			Label: "expression",
			Markup: map[string]any{
				"comment": "The expression to match against the text.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "cache",
			Label:   "cache",
			Private: true,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_BoolEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Determine whether text matches a regular expression.", "The expressions used are defined by go.", "https://pkg.go.dev/regexp/syntax", "https://github.com/google/re2/wiki/Syntax"},
		},
	}
	Zt_Capitalize = typeinfo.Flow{
		Name: "capitalize",
		Lede: "capitalize",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to copy, then capitalize.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Copy some text, changing its first letter to uppercase.",
		},
	}
	Zt_Join = typeinfo.Flow{
		Name: "join",
		Lede: "join",
		Terms: []typeinfo.Term{{
			Name:     "sep",
			Optional: true,
			Markup: map[string]any{
				"comment": "Optionally, a separator to put between each piece of text.",
			},
			Type: &rtti.Zt_TextEval,
		}, {
			Name:    "parts",
			Label:   "parts",
			Repeats: true,
			Markup: map[string]any{
				"comment": "The pieces of text to combine.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": "Combine text to produce new text.",
		},
	}
	Zt_MakeLowercase = typeinfo.Flow{
		Name: "make_lowercase",
		Lede: "lower",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to copy, then lowercase.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Copy some text, changing every letter into lowercase.", "For example, turns \"QUIET\" into \"quiet."},
		},
	}
	Zt_MakeReversed = typeinfo.Flow{
		Name: "make_reversed",
		Lede: "reverse",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "text",
			Markup: map[string]any{
				"comment": "The text to copy and then reverse.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Copy some text with its contents flipped back to front.", "For example, turns \"Tapestry\" into 'yrtsepaT'."},
		},
	}
	Zt_MakeSentenceCase = typeinfo.Flow{
		Name: "make_sentence_case",
		Lede: "sentence",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to copy and then transform.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Copy text, changing the start of each sentence so that it starts with a capital letter. ( Currently, \"sentences\" are considered to be a series of characters ending with a full-stop followed by a space. )", "For example, \"see the doctor run. run doctor. run.\" into \"See the doctor run. Run doctor. Run.\""},
		},
	}
	Zt_MakeTitleCase = typeinfo.Flow{
		Name: "make_title_case",
		Lede: "title",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to copy and then transform.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Copy some text, making every word start with a capital letter.", "For example, turns \"empire apple\" into \"Empire Apple\"."},
		},
	}
	Zt_MakeUppercase = typeinfo.Flow{
		Name: "make_uppercase",
		Lede: "upper",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to copy and then transform into uppercase.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Copy some text, changing every letter into uppercase.", "For example, transforms \"loud\" into \"LOUD\"."},
		},
	}
	Zt_Pluralize = typeinfo.Flow{
		Name: "pluralize",
		Lede: "plural",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "of",
			Markup: map[string]any{
				"comment": "The text to pluralize.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Pluralize a word.", "The singular form of a word can have more than one plural form.", "For example: \"person\" can be \"people\" or \"persons\".", "If more than one exists, this chooses arbitrarily.", "", "Note, The transformation uses predefined rules and some explicit mappings.", "The story command [DefinePlural] can add new mappings."},
		},
	}
	Zt_Singularize = typeinfo.Flow{
		Name: "singularize",
		Lede: "singular",
		Terms: []typeinfo.Term{{
			Name:  "text",
			Label: "of",
			Markup: map[string]any{
				"comment": "The text to turn into its singular form.",
			},
			Type: &rtti.Zt_TextEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Change a plural word into its singular form.", "A plural word only has one singular form.", "For example, given the word \"people\", return \"person\".", "See [pluralize] for more information."},
		},
	}
	Zt_PrintNumDigits = typeinfo.Flow{
		Name: "print_num_digits",
		Lede: "numeral",
		Terms: []typeinfo.Term{{
			Name:  "num",
			Label: "digits",
			Markup: map[string]any{
				"comment": "The number to change into text, or to print.",
			},
			Type: &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Express a number using numerals.", "For example, given the number `1` return the text \"1\".", "", "The [story.Execute] version prints the text for the player."},
		},
	}
	Zt_PrintNumWords = typeinfo.Flow{
		Name: "print_num_words",
		Lede: "numeral",
		Terms: []typeinfo.Term{{
			Name:  "num",
			Label: "words",
			Markup: map[string]any{
				"comment": "The number to change into words, or to print.",
			},
			Type: &rtti.Zt_NumEval,
		}},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_TextEval,
			&rtti.Zt_Execute,
		},
		Markup: map[string]any{
			"comment": []interface{}{"Express an integer in plain english.", "For example, given the number `1` return the text \"one\".", "It converts floating point numbers to integer by truncating:", "given `1.6`, it returns \"one\".", "", "The [story.Execute] version prints the text for the player."},
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "text",
	Comment: []string{
		"Text manipulation and transformation.",
	},

	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_TextLen,
	&Zt_FindText,
	&Zt_TextStartsWith,
	&Zt_TextEndsWith,
	&Zt_IsEmpty,
	&Zt_Matches,
	&Zt_Capitalize,
	&Zt_Join,
	&Zt_MakeLowercase,
	&Zt_MakeReversed,
	&Zt_MakeSentenceCase,
	&Zt_MakeTitleCase,
	&Zt_MakeUppercase,
	&Zt_Pluralize,
	&Zt_Singularize,
	&Zt_PrintNumDigits,
	&Zt_PrintNumWords,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	8695677004499439692:  (*Capitalize)(nil),       /* text_eval=Capitalize: */
	7190419643383427713:  (*FindText)(nil),         /* bool_eval=Find:text: */
	2922486476891061679:  (*FindText)(nil),         /* num_eval=Find:text: */
	10867951538760575464: (*IsEmpty)(nil),          /* bool_eval=Is empty: */
	6139101839500298168:  (*Matches)(nil),          /* bool_eval=Is text:expression: */
	43416298232103202:    (*TextStartsWith)(nil),   /* bool_eval=Is text:prefix: */
	14194170362800670601: (*TextEndsWith)(nil),     /* bool_eval=Is text:suffix: */
	10106284345457008764: (*Join)(nil),             /* text_eval=Join parts: */
	16037301925772243654: (*Join)(nil),             /* text_eval=Join:parts: */
	11334467785012784241: (*MakeLowercase)(nil),    /* text_eval=Lower: */
	4721393964025254579:  (*PrintNumDigits)(nil),   /* execute=Numeral digits: */
	14515844015968836994: (*PrintNumDigits)(nil),   /* text_eval=Numeral digits: */
	9655583796217513308:  (*PrintNumWords)(nil),    /* execute=Numeral words: */
	18009133328614046007: (*PrintNumWords)(nil),    /* text_eval=Numeral words: */
	11420921600352749983: (*Pluralize)(nil),        /* text_eval=Plural of: */
	12963686195606417453: (*MakeReversed)(nil),     /* text_eval=Reverse text: */
	10747671703915852065: (*MakeSentenceCase)(nil), /* text_eval=Sentence: */
	2397382738676796596:  (*Singularize)(nil),      /* text_eval=Singular of: */
	5968278258259888036:  (*TextLen)(nil),          /* num_eval=Text length: */
	10878271994667616824: (*MakeTitleCase)(nil),    /* text_eval=Title: */
	5481656653805454214:  (*MakeUppercase)(nil),    /* text_eval=Upper: */
}
