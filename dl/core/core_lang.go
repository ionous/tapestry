// Code generated by "makeops"; edit at your own risk.
package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/reader"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

// Activity
type Activity struct {
	Exe []rt.Execute `if:"label=do"`
}

func (*Activity) Compose() composer.Spec {
	return composer.Spec{
		Name: "activity",
		Lede: "act",
	}
}

// AllTrue Returns true if all of the evaluations are true.
type AllTrue struct {
	Test []rt.BoolEval `if:"label=_"`
}

func (*AllTrue) Compose() composer.Spec {
	return composer.Spec{
		Name: "all_true",
		Lede: "all_of",
	}
}

// Always Returns true always.
type Always struct {
}

func (*Always) Compose() composer.Spec {
	return composer.Spec{
		Name: "always",
	}
}

// AnyTrue Returns true if any of the evaluations are true.
type AnyTrue struct {
	Test []rt.BoolEval `if:"label=_"`
}

func (*AnyTrue) Compose() composer.Spec {
	return composer.Spec{
		Name: "any_true",
		Lede: "any_of",
	}
}

// Argument
type Argument struct {
	Name value.Text    `if:"label=_"`
	From rt.Assignment `if:"label=from"`
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Name: "argument",
		Lede: "arg",
	}
}

// Arguments
type Arguments struct {
	Args []Argument `if:"label=_"`
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Name: "arguments",
	}
}

// Assign Assigns a variable to a value.
type Assign struct {
	Var  string        `if:"label=_"`
	From rt.Assignment `if:"label=be"`
}

func (*Assign) Compose() composer.Spec {
	return composer.Spec{
		Name: "assign",
		Lede: "let",
	}
}

// Blankline Add a single blank line following some text.
type Blankline struct {
}

func (*Blankline) Compose() composer.Spec {
	return composer.Spec{
		Name: "blankline",
		Lede: "p",
	}
}

// BoolValue Specify an explicit true or false value.
type BoolValue struct {
	Bool bool `if:"label=_"`
}

func (*BoolValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "bool_value",
		Lede: "bool",
	}
}

// Bracket Sandwiches text printed during a block and puts them inside parenthesis &#x27;()&#x27;.
type Bracket struct {
	Do Activity `if:"label=_"`
}

func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Name: "bracket",
		Lede: "bracket_text",
	}
}

// Break In a repeating loop, exit the loop.
type Break struct {
}

func (*Break) Compose() composer.Spec {
	return composer.Spec{
		Name: "break",
	}
}

// Buffer
type Buffer struct {
	Do Activity `if:"label=_"`
}

func (*Buffer) Compose() composer.Spec {
	return composer.Spec{
		Name: "buffer",
		Lede: "buffer_text",
	}
}

// Capitalize Returns new text, with the first letter turned into uppercase.
type Capitalize struct {
	Text rt.TextEval `if:"label=_"`
}

func (*Capitalize) Compose() composer.Spec {
	return composer.Spec{
		Name: "capitalize",
	}
}

// ChooseAction An if statement.
type ChooseAction struct {
	If   rt.BoolEval `if:"label=_"`
	Do   Activity    `if:"label=do"`
	Else Brancher    `if:"label=else,optional"`
}

func (*ChooseAction) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_action",
		Lede: "if",
	}
}

// ChooseMore
type ChooseMore struct {
	If   rt.BoolEval `if:"label=_"`
	Do   Activity    `if:"label=do"`
	Else Brancher    `if:"label=else,optional"`
}

func (*ChooseMore) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_more",
		Lede: "else_if",
	}
}

// ChooseMoreValue
type ChooseMoreValue struct {
	Assign value.Text    `if:"label=_"`
	From   rt.Assignment `if:"label=from"`
	Filter rt.BoolEval   `if:"label=and"`
	Do     Activity      `if:"label=do"`
	Else   Brancher      `if:"label=else,optional"`
}

func (*ChooseMoreValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_more_value",
		Lede: "else_if",
	}
}

// ChooseNothingElse
type ChooseNothingElse struct {
	Do Activity `if:"label=_"`
}

func (*ChooseNothingElse) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_nothing_else",
		Lede: "else_do",
	}
}

// ChooseNum Pick one of two numbers based on a boolean test.
type ChooseNum struct {
	If    rt.BoolEval   `if:"label=if"`
	True  rt.NumberEval `if:"label=then"`
	False rt.NumberEval `if:"label=else"`
}

func (*ChooseNum) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_num",
		Lede: "num",
	}
}

// ChooseText Pick one of two strings based on a boolean test.
type ChooseText struct {
	If    rt.BoolEval `if:"label=if"`
	True  rt.TextEval `if:"label=then"`
	False rt.TextEval `if:"label=else"`
}

func (*ChooseText) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_text",
		Lede: "txt",
	}
}

// ChooseValue An if statement with local assignment.
type ChooseValue struct {
	Assign value.Text    `if:"label=_"`
	From   rt.Assignment `if:"label=from"`
	Filter rt.BoolEval   `if:"label=and"`
	Do     Activity      `if:"label=do"`
	Else   Brancher      `if:"label=else,optional"`
}

func (*ChooseValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose_value",
		Lede: "if",
	}
}

// Commas Separates words with commas, and &#x27;and&#x27;.
type Commas struct {
	Do Activity `if:"label=_"`
}

func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Name: "commas",
		Lede: "comma_text",
	}
}

// CompareNum True if eq,ne,gt,lt,ge,le two numbers.
type CompareNum struct {
	A  rt.NumberEval `if:"label=_"`
	Is Comparator    `if:"label=is"`
	B  rt.NumberEval `if:"label=num"`
}

func (*CompareNum) Compose() composer.Spec {
	return composer.Spec{
		Name: "compare_num",
		Lede: "cmp",
	}
}

// CompareText True if eq,ne,gt,lt,ge,le two strings ( lexical. )
type CompareText struct {
	A  rt.TextEval `if:"label=_"`
	Is Comparator  `if:"label=is"`
	B  rt.TextEval `if:"label=txt"`
}

func (*CompareText) Compose() composer.Spec {
	return composer.Spec{
		Name: "compare_text",
		Lede: "cmp",
	}
}

// CountOf A guard which returns true based on a counter. Counters start at zero and are incremented every time the guard gets checked.
type CountOf struct {
	At      reader.Position `if:"internal"`
	Trigger Trigger         `if:"label=_"`
	Num     rt.NumberEval   `if:"label=num"`
}

func (*CountOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "count_of",
		Lede: "trigger",
	}
}

// CycleText When called multiple times, returns each of its inputs in turn.
type CycleText struct {
	At    reader.Position `if:"internal"`
	Parts []rt.TextEval   `if:"label=_"`
}

func (*CycleText) Compose() composer.Spec {
	return composer.Spec{
		Name: "cycle_text",
		Lede: "cycle",
	}
}

// Determine Runs a pattern, and potentially returns a value.
type Determine struct {
	Pattern   string    `if:"label=_"`
	Arguments Arguments `if:"label=arguments"`
}

func (*Determine) Compose() composer.Spec {
	return composer.Spec{
		Name: "determine",
	}
}

// DiffOf Subtract two numbers.
type DiffOf struct {
	A rt.NumberEval `if:"label=_"`
	B rt.NumberEval `if:"label=by,optional"`
}

func (*DiffOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "diff_of",
		Lede: "dec",
	}
}

// During Decide whether a pattern is running.
type During struct {
	Pattern string `if:"label=_"`
}

func (*During) Compose() composer.Spec {
	return composer.Spec{
		Name: "during",
	}
}

// EqualTo Two values exactly match.
type EqualTo struct {
}

func (*EqualTo) Compose() composer.Spec {
	return composer.Spec{
		Name: "equal_to",
	}
}

// FromBool Assigns the calculated boolean value.
type FromBool struct {
	Val rt.BoolEval `if:"label=_"`
}

func (*FromBool) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_bool",
		Lede: "bool",
	}
}

// FromNum Assigns the calculated number.
type FromNum struct {
	Val rt.NumberEval `if:"label=_"`
}

func (*FromNum) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_num",
		Lede: "num",
	}
}

// FromNumbers Assigns the calculated numbers.
type FromNumbers struct {
	Vals rt.NumListEval `if:"label=_"`
}

func (*FromNumbers) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_numbers",
		Lede: "nums",
	}
}

// FromObj Targets an object with a computed name.
type FromObj struct {
	Object rt.TextEval `if:"label=_"`
}

func (*FromObj) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_obj",
		Lede: "obj",
	}
}

// FromRec Targets a record stored in a record.
type FromRec struct {
	Rec rt.RecordEval `if:"label=_"`
}

func (*FromRec) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_rec",
		Lede: "rec",
	}
}

// FromRecord Assigns the calculated record.
type FromRecord struct {
	Val rt.RecordEval `if:"label=_"`
}

func (*FromRecord) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_record",
		Lede: "rec",
	}
}

// FromRecords Assigns the calculated records.
type FromRecords struct {
	Vals rt.RecordListEval `if:"label=_"`
}

func (*FromRecords) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_records",
		Lede: "recs",
	}
}

// FromText Assigns the calculated piece of text.
type FromText struct {
	Val rt.TextEval `if:"label=_"`
}

func (*FromText) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_text",
		Lede: "txt",
	}
}

// FromTexts Assigns the calculated texts.
type FromTexts struct {
	Vals rt.TextListEval `if:"label=_"`
}

func (*FromTexts) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_texts",
		Lede: "txts",
	}
}

// FromVar Targets a record stored in a variable.
type FromVar struct {
	Var string `if:"label=_"`
}

func (*FromVar) Compose() composer.Spec {
	return composer.Spec{
		Name: "from_var",
		Lede: "var",
	}
}

// GetAtField Get a value from a record.
type GetAtField struct {
	Field value.Text       `if:"label=_"`
	From  FromSourceFields `if:"label=from"`
}

func (*GetAtField) Compose() composer.Spec {
	return composer.Spec{
		Name: "get_at_field",
		Lede: "get",
	}
}

// GreaterOrEqual The first value is larger than the second value.
type GreaterOrEqual struct {
}

func (*GreaterOrEqual) Compose() composer.Spec {
	return composer.Spec{
		Name: "greater_or_equal",
		Lede: "at_least",
	}
}

// GreaterThan The first value is larger than the second value.
type GreaterThan struct {
}

func (*GreaterThan) Compose() composer.Spec {
	return composer.Spec{
		Name: "greater_than",
	}
}

// HasDominion
type HasDominion struct {
	Domain value.Text `if:"label=_"`
}

func (*HasDominion) Compose() composer.Spec {
	return composer.Spec{
		Name: "has_dominion",
	}
}

// HasTrait Return true if the object is currently in the requested state.
type HasTrait struct {
	Object rt.TextEval `if:"label=obj"`
	Trait  rt.TextEval `if:"label=trait"`
}

func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Name: "has_trait",
		Lede: "get",
	}
}

// IdOf A unique object identifier.
type IdOf struct {
	Object rt.TextEval `if:"label=_"`
}

func (*IdOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "id_of",
	}
}

// Includes True if text contains text.
type Includes struct {
	Text rt.TextEval `if:"label=_"`
	Part rt.TextEval `if:"label=part"`
}

func (*Includes) Compose() composer.Spec {
	return composer.Spec{
		Name: "includes",
		Lede: "contains",
	}
}

// IntoObj Targets an object with a computed name.
type IntoObj struct {
	Object rt.TextEval `if:"label=_"`
}

func (*IntoObj) Compose() composer.Spec {
	return composer.Spec{
		Name: "into_obj",
		Lede: "obj",
	}
}

// IntoVar Targets an object or record stored in a variable
type IntoVar struct {
	Var string `if:"label=_"`
}

func (*IntoVar) Compose() composer.Spec {
	return composer.Spec{
		Name: "into_var",
		Lede: "var",
	}
}

// IsEmpty True if the text is empty.
type IsEmpty struct {
	Text rt.TextEval `if:"label=empty"`
}

func (*IsEmpty) Compose() composer.Spec {
	return composer.Spec{
		Name: "is_empty",
		Lede: "is",
	}
}

// IsExactKindOf True if the object is exactly the named kind.
type IsExactKindOf struct {
	Object rt.TextEval `if:"label=_"`
	Kind   value.Text  `if:"label=is_exactly"`
}

func (*IsExactKindOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "is_exact_kind_of",
		Lede: "kind_of",
	}
}

// IsKindOf True if the object is compatible with the named kind.
type IsKindOf struct {
	Object rt.TextEval `if:"label=_"`
	Kind   value.Text  `if:"label=is"`
}

func (*IsKindOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "is_kind_of",
		Lede: "kind_of",
	}
}

// Join Returns multiple pieces of text as a single new piece of text.
type Join struct {
	Sep   rt.TextEval   `if:"label=_"`
	Parts []rt.TextEval `if:"label=parts"`
}

func (*Join) Compose() composer.Spec {
	return composer.Spec{
		Name: "join",
	}
}

// KindOf Friendly name of the object&#x27;s kind.
type KindOf struct {
	Object rt.TextEval `if:"label=_"`
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "kind_of",
	}
}

// KindsOf A list of compatible kinds.
type KindsOf struct {
	Kind value.Text `if:"label=_"`
}

func (*KindsOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "kinds_of",
	}
}

// LessOrEqual The first value is larger than the second value.
type LessOrEqual struct {
}

func (*LessOrEqual) Compose() composer.Spec {
	return composer.Spec{
		Name: "less_or_equal",
		Lede: "at_most",
	}
}

// LessThan The first value is less than the second value.
type LessThan struct {
}

func (*LessThan) Compose() composer.Spec {
	return composer.Spec{
		Name: "less_than",
	}
}

// Make
type Make struct {
	Kind      value.Text `if:"label=_"`
	Arguments Arguments  `if:"label=arguments"`
}

func (*Make) Compose() composer.Spec {
	return composer.Spec{
		Name: "make",
	}
}

// MakeLowercase Returns new text, with every letter turned into lowercase. For example, &#x27;shout&#x27; from &#x27;SHOUT&#x27;.
type MakeLowercase struct {
	Text rt.TextEval `if:"label=_"`
}

func (*MakeLowercase) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_lowercase",
		Lede: "lower",
	}
}

// MakePlural Returns the plural form of a singular word. (ex. apples for apple. )
type MakePlural struct {
	Text rt.TextEval `if:"label=of"`
}

func (*MakePlural) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_plural",
		Lede: "plural",
	}
}

// MakeReversed Returns new text flipped back to front. For example, &#x27;elppA&#x27; from &#x27;Apple&#x27;, or &#x27;noon&#x27; from &#x27;noon&#x27;.
type MakeReversed struct {
	Text rt.TextEval `if:"label=_"`
}

func (*MakeReversed) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_reversed",
		Lede: "reverse",
	}
}

// MakeSentenceCase Returns new text, start each sentence with a capital letter. For example, &#x27;Empire Apple.&#x27; from &#x27;Empire apple.&#x27;.
type MakeSentenceCase struct {
	Text rt.TextEval `if:"label=_"`
}

func (*MakeSentenceCase) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_sentence_case",
		Lede: "sentence",
	}
}

// MakeSingular Returns the singular form of a plural word. (ex. apple for apples )
type MakeSingular struct {
	Text rt.TextEval `if:"label=of"`
}

func (*MakeSingular) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_singular",
		Lede: "singular",
	}
}

// MakeTitleCase Returns new text, starting each word with a capital letter. For example, &#x27;Empire Apple&#x27; from &#x27;empire apple&#x27;.
type MakeTitleCase struct {
	Text rt.TextEval `if:"label=_"`
}

func (*MakeTitleCase) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_title_case",
		Lede: "title",
	}
}

// MakeUppercase Returns new text, with every letter turned into uppercase. For example, &#x27;APPLE&#x27; from &#x27;apple&#x27;.
type MakeUppercase struct {
	Text rt.TextEval `if:"label=_"`
}

func (*MakeUppercase) Compose() composer.Spec {
	return composer.Spec{
		Name: "make_uppercase",
		Lede: "upper",
	}
}

// Matches Determine whether the specified text is similar to the specified regular expression.
type Matches struct {
	Text    rt.TextEval `if:"label=_"`
	Pattern string      `if:"label=to"`
	Cache   MatchCache  `if:"internal"`
}

func (*Matches) Compose() composer.Spec {
	return composer.Spec{
		Name: "matches",
	}
}

// NameOf Full name of the object.
type NameOf struct {
	Object rt.TextEval `if:"label=_"`
}

func (*NameOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "name_of",
	}
}

// Newline Start a new line.
type Newline struct {
}

func (*Newline) Compose() composer.Spec {
	return composer.Spec{
		Name: "newline",
		Lede: "br",
	}
}

// Next In a repeating loop, try the next iteration of the loop.
type Next struct {
}

func (*Next) Compose() composer.Spec {
	return composer.Spec{
		Name: "next",
	}
}

// Not Returns the opposite value.
type Not struct {
	Test rt.BoolEval `if:"label=_"`
}

func (*Not) Compose() composer.Spec {
	return composer.Spec{
		Name: "not",
	}
}

// NotEqualTo Two values don&#x27;t match exactly.
type NotEqualTo struct {
}

func (*NotEqualTo) Compose() composer.Spec {
	return composer.Spec{
		Name: "not_equal_to",
		Lede: "other_than",
	}
}

// NumList Specify a list of multiple numbers.
type NumList struct {
	Values []float64 `if:"label=_"`
}

func (*NumList) Compose() composer.Spec {
	return composer.Spec{
		Name: "num_list",
		Lede: "nums",
	}
}

// NumValue Specify a particular number.
type NumValue struct {
	Num float64 `if:"label=_"`
}

func (*NumValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "num_value",
		Lede: "num",
	}
}

// ObjectExists Returns whether there is a object of the specified name.
type ObjectExists struct {
	Object rt.TextEval `if:"label=valid"`
}

func (*ObjectExists) Compose() composer.Spec {
	return composer.Spec{
		Name: "object_exists",
		Lede: "is",
	}
}

// PrintNum Writes a number using numerals, eg. &#x27;1&#x27;.
type PrintNum struct {
	Num rt.NumberEval `if:"label=_"`
}

func (*PrintNum) Compose() composer.Spec {
	return composer.Spec{
		Name: "print_num",
		Lede: "numeral",
	}
}

// PrintNumWord Writes a number in plain english: eg. &#x27;one&#x27;
type PrintNumWord struct {
	Num rt.NumberEval `if:"label=words"`
}

func (*PrintNumWord) Compose() composer.Spec {
	return composer.Spec{
		Name: "print_num_word",
		Lede: "numeral",
	}
}

// ProductOf Multiply two numbers.
type ProductOf struct {
	A rt.NumberEval `if:"label=_"`
	B rt.NumberEval `if:"label=by"`
}

func (*ProductOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "product_of",
		Lede: "mul",
	}
}

// PutAtField Put a value into the field of an record or object
type PutAtField struct {
	Into    IntoTargetFields `if:"label=_"`
	From    rt.Assignment    `if:"label=from"`
	AtField value.Text       `if:"label=at"`
}

func (*PutAtField) Compose() composer.Spec {
	return composer.Spec{
		Name: "put_at_field",
		Lede: "put",
	}
}

// QuotientOf Divide one number by another.
type QuotientOf struct {
	A rt.NumberEval `if:"label=_"`
	B rt.NumberEval `if:"label=by"`
}

func (*QuotientOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "quotient_of",
		Lede: "div",
	}
}

// RemainderOf Divide one number by another, and return the remainder.
type RemainderOf struct {
	A rt.NumberEval `if:"label=_"`
	B rt.NumberEval `if:"label=by"`
}

func (*RemainderOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "remainder_of",
		Lede: "mod",
	}
}

// Response Generate text in a replaceable manner.
type Response struct {
	Name value.Text  `if:"label=_"`
	Text rt.TextEval `if:"label=text,optional"`
}

func (*Response) Compose() composer.Spec {
	return composer.Spec{
		Name: "response",
	}
}

// Row A single line as part of a group of lines.
type Row struct {
	Do Activity `if:"label=_"`
}

func (*Row) Compose() composer.Spec {
	return composer.Spec{
		Name: "row",
	}
}

// Rows Group text into successive lines.
type Rows struct {
	Do Activity `if:"label=_"`
}

func (*Rows) Compose() composer.Spec {
	return composer.Spec{
		Name: "rows",
	}
}

// Say Print some bit of text to the player.
type Say struct {
	Text rt.TextEval `if:"label=_"`
}

func (*Say) Compose() composer.Spec {
	return composer.Spec{
		Name: "say",
		Lede: "say_text",
	}
}

// Send Triggers a event, returns a true/false success value.
type Send struct {
	Event     value.Text      `if:"label=_"`
	Path      rt.TextListEval `if:"label=to"`
	Arguments Arguments       `if:"label=arguments"`
}

func (*Send) Compose() composer.Spec {
	return composer.Spec{
		Name: "send",
	}
}

// SetTrait Put an object into a particular state.
type SetTrait struct {
	Object rt.TextEval `if:"label=obj"`
	Trait  rt.TextEval `if:"label=trait"`
}

func (*SetTrait) Compose() composer.Spec {
	return composer.Spec{
		Name: "set_trait",
		Lede: "put",
	}
}

// ShuffleText When called multiple times returns its inputs at random.
type ShuffleText struct {
	At      reader.Position `if:"internal"`
	Parts   []rt.TextEval   `if:"label=_"`
	Indices Shuffler        `if:"internal"`
}

func (*ShuffleText) Compose() composer.Spec {
	return composer.Spec{
		Name: "shuffle_text",
		Lede: "shuffle",
	}
}

// Slash Separates words with left-leaning slashes &#x27;/&#x27;.
type Slash struct {
	Do Activity `if:"label=_"`
}

func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Name: "slash",
		Lede: "slash_text",
	}
}

// Softline Start a new line ( if not already at a new line. )
type Softline struct {
}

func (*Softline) Compose() composer.Spec {
	return composer.Spec{
		Name: "softline",
		Lede: "wbr",
	}
}

// Span Writes text with spaces between words.
type Span struct {
	Do Activity `if:"label=_"`
}

func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Name: "span",
		Lede: "span_text",
	}
}

// StoppingText When called multiple times returns each of its inputs in turn, sticking to the last one.
type StoppingText struct {
	At    reader.Position `if:"internal"`
	Parts []rt.TextEval   `if:"label=_"`
}

func (*StoppingText) Compose() composer.Spec {
	return composer.Spec{
		Name: "stopping_text",
		Lede: "stopping",
	}
}

// SumOf Add two numbers.
type SumOf struct {
	A rt.NumberEval `if:"label=_"`
	B rt.NumberEval `if:"label=by,optional"`
}

func (*SumOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "sum_of",
		Lede: "inc",
	}
}

// TextList Specifies multiple string values.
type TextList struct {
	Values []value.Text `if:"label=_"`
}

func (*TextList) Compose() composer.Spec {
	return composer.Spec{
		Name: "text_list",
		Lede: "txts",
	}
}

// TextValue Specify a small bit of text.
type TextValue struct {
	Text value.Text `if:"label=_"`
}

func (*TextValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "text_value",
		Lede: "txt",
	}
}

// TriggerCycle
type TriggerCycle struct {
}

func (*TriggerCycle) Compose() composer.Spec {
	return composer.Spec{
		Name: "trigger_cycle",
		Lede: "every",
	}
}

// TriggerOnce
type TriggerOnce struct {
}

func (*TriggerOnce) Compose() composer.Spec {
	return composer.Spec{
		Name: "trigger_once",
		Lede: "at",
	}
}

// TriggerSwitch
type TriggerSwitch struct {
}

func (*TriggerSwitch) Compose() composer.Spec {
	return composer.Spec{
		Name: "trigger_switch",
		Lede: "after",
	}
}

// Var Return the value of the named variable.
type Var struct {
	Name string `if:"label=_"`
}

func (*Var) Compose() composer.Spec {
	return composer.Spec{
		Name: "var",
	}
}

// While Keep running a series of actions while a condition is true.
type While struct {
	True rt.BoolEval `if:"label=_"`
	Do   Activity    `if:"label=do"`
}

func (*While) Compose() composer.Spec {
	return composer.Spec{
		Name: "while",
		Lede: "repeating",
	}
}

var Slots = []interface{}{
	(*Brancher)(nil),
	(*Comparator)(nil),
	(*FromSourceFields)(nil),
	(*IntoTargetFields)(nil),
	(*Trigger)(nil),
}
var Flows = []interface{}{
	(*Activity)(nil),
	(*AllTrue)(nil),
	(*Always)(nil),
	(*AnyTrue)(nil),
	(*Argument)(nil),
	(*Arguments)(nil),
	(*Assign)(nil),
	(*Blankline)(nil),
	(*BoolValue)(nil),
	(*Bracket)(nil),
	(*Break)(nil),
	(*Buffer)(nil),
	(*Capitalize)(nil),
	(*ChooseAction)(nil),
	(*ChooseMore)(nil),
	(*ChooseMoreValue)(nil),
	(*ChooseNothingElse)(nil),
	(*ChooseNum)(nil),
	(*ChooseText)(nil),
	(*ChooseValue)(nil),
	(*Commas)(nil),
	(*CompareNum)(nil),
	(*CompareText)(nil),
	(*CountOf)(nil),
	(*CycleText)(nil),
	(*Determine)(nil),
	(*DiffOf)(nil),
	(*During)(nil),
	(*EqualTo)(nil),
	(*FromBool)(nil),
	(*FromNum)(nil),
	(*FromNumbers)(nil),
	(*FromObj)(nil),
	(*FromRec)(nil),
	(*FromRecord)(nil),
	(*FromRecords)(nil),
	(*FromText)(nil),
	(*FromTexts)(nil),
	(*FromVar)(nil),
	(*GetAtField)(nil),
	(*GreaterOrEqual)(nil),
	(*GreaterThan)(nil),
	(*HasDominion)(nil),
	(*HasTrait)(nil),
	(*IdOf)(nil),
	(*Includes)(nil),
	(*IntoObj)(nil),
	(*IntoVar)(nil),
	(*IsEmpty)(nil),
	(*IsExactKindOf)(nil),
	(*IsKindOf)(nil),
	(*Join)(nil),
	(*KindOf)(nil),
	(*KindsOf)(nil),
	(*LessOrEqual)(nil),
	(*LessThan)(nil),
	(*Make)(nil),
	(*MakeLowercase)(nil),
	(*MakePlural)(nil),
	(*MakeReversed)(nil),
	(*MakeSentenceCase)(nil),
	(*MakeSingular)(nil),
	(*MakeTitleCase)(nil),
	(*MakeUppercase)(nil),
	(*Matches)(nil),
	(*NameOf)(nil),
	(*Newline)(nil),
	(*Next)(nil),
	(*Not)(nil),
	(*NotEqualTo)(nil),
	(*NumList)(nil),
	(*NumValue)(nil),
	(*ObjectExists)(nil),
	(*PrintNum)(nil),
	(*PrintNumWord)(nil),
	(*ProductOf)(nil),
	(*PutAtField)(nil),
	(*QuotientOf)(nil),
	(*RemainderOf)(nil),
	(*Response)(nil),
	(*Row)(nil),
	(*Rows)(nil),
	(*Say)(nil),
	(*Send)(nil),
	(*SetTrait)(nil),
	(*ShuffleText)(nil),
	(*Slash)(nil),
	(*Softline)(nil),
	(*Span)(nil),
	(*StoppingText)(nil),
	(*SumOf)(nil),
	(*TextList)(nil),
	(*TextValue)(nil),
	(*TriggerCycle)(nil),
	(*TriggerOnce)(nil),
	(*TriggerSwitch)(nil),
	(*Var)(nil),
	(*While)(nil),
}
