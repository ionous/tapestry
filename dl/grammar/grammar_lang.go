// Code generated by "makeops"; edit at your own risk.
package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/export/jsn"
)

// Action makes a parser scanner producing a script defined action.
type Action struct {
	Action string `if:"label=_,type=text"`
}

func (*Action) Compose() composer.Spec {
	return composer.Spec{
		Name: Action_Type,
		Uses: composer.Type_Flow,
		Lede: "as",
	}
}

const Action_Type = "action"

const Action_Field_Action = "$ACTION"

func (op *Action) Marshal(n jsn.Marshaler) {
	Action_Marshal(n, op)
}

func Action_Repeats_Marshal(n jsn.Marshaler, vals *[]Action) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Action_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Action_Optional_Marshal(n jsn.Marshaler, val **Action) {
	if *val != nil {
		Action_Marshal(n, *val)
	}
}

func Action_Marshal(n jsn.Marshaler, val *Action) {
	n.MapValues("as", Action_Type)
	n.MapKey("", Action_Field_Action)
	/* */ value.Text_Unboxed_Marshal(n, &val.Action)
	n.EndValues()
	return
}

// Alias allows the user to refer to a noun by one or more other terms.
type Alias struct {
	Names  []string `if:"label=_,type=text"`
	AsNoun string   `if:"label=as_noun,type=text"`
}

func (*Alias) Compose() composer.Spec {
	return composer.Spec{
		Name: Alias_Type,
		Uses: composer.Type_Flow,
	}
}

const Alias_Type = "alias"

const Alias_Field_Names = "$NAMES"
const Alias_Field_AsNoun = "$AS_NOUN"

func (op *Alias) Marshal(n jsn.Marshaler) {
	Alias_Marshal(n, op)
}

func Alias_Repeats_Marshal(n jsn.Marshaler, vals *[]Alias) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Alias_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Alias_Optional_Marshal(n jsn.Marshaler, val **Alias) {
	if *val != nil {
		Alias_Marshal(n, *val)
	}
}

func Alias_Marshal(n jsn.Marshaler, val *Alias) {
	n.MapValues(Alias_Type, Alias_Type)
	n.MapKey("", Alias_Field_Names)
	/* */ value.Text_Unboxed_Repeats_Marshal(n, &val.Names)
	n.MapKey("as_noun", Alias_Field_AsNoun)
	/* */ value.Text_Unboxed_Marshal(n, &val.AsNoun)
	n.EndValues()
	return
}

// AllOf makes a parser scanner
type AllOf struct {
	Series []ScannerMaker `if:"label=_"`
}

func (*AllOf) Compose() composer.Spec {
	return composer.Spec{
		Name: AllOf_Type,
		Uses: composer.Type_Flow,
	}
}

const AllOf_Type = "all_of"

const AllOf_Field_Series = "$SERIES"

func (op *AllOf) Marshal(n jsn.Marshaler) {
	AllOf_Marshal(n, op)
}

func AllOf_Repeats_Marshal(n jsn.Marshaler, vals *[]AllOf) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			AllOf_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func AllOf_Optional_Marshal(n jsn.Marshaler, val **AllOf) {
	if *val != nil {
		AllOf_Marshal(n, *val)
	}
}

func AllOf_Marshal(n jsn.Marshaler, val *AllOf) {
	n.MapValues(AllOf_Type, AllOf_Type)
	n.MapKey("", AllOf_Field_Series)
	/* */ ScannerMaker_Repeats_Marshal(n, &val.Series)
	n.EndValues()
	return
}

// AnyOf makes a parser scanner
type AnyOf struct {
	Options []ScannerMaker `if:"label=_"`
}

func (*AnyOf) Compose() composer.Spec {
	return composer.Spec{
		Name: AnyOf_Type,
		Uses: composer.Type_Flow,
	}
}

const AnyOf_Type = "any_of"

const AnyOf_Field_Options = "$OPTIONS"

func (op *AnyOf) Marshal(n jsn.Marshaler) {
	AnyOf_Marshal(n, op)
}

func AnyOf_Repeats_Marshal(n jsn.Marshaler, vals *[]AnyOf) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			AnyOf_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func AnyOf_Optional_Marshal(n jsn.Marshaler, val **AnyOf) {
	if *val != nil {
		AnyOf_Marshal(n, *val)
	}
}

func AnyOf_Marshal(n jsn.Marshaler, val *AnyOf) {
	n.MapValues(AnyOf_Type, AnyOf_Type)
	n.MapKey("", AnyOf_Field_Options)
	/* */ ScannerMaker_Repeats_Marshal(n, &val.Options)
	n.EndValues()
	return
}

// Directive starts a parser scanner
type Directive struct {
	Lede  []string       `if:"label=_,type=text"`
	Scans []ScannerMaker `if:"label=scans"`
}

func (*Directive) Compose() composer.Spec {
	return composer.Spec{
		Name: Directive_Type,
		Uses: composer.Type_Flow,
	}
}

const Directive_Type = "directive"

const Directive_Field_Lede = "$LEDE"
const Directive_Field_Scans = "$SCANS"

func (op *Directive) Marshal(n jsn.Marshaler) {
	Directive_Marshal(n, op)
}

func Directive_Repeats_Marshal(n jsn.Marshaler, vals *[]Directive) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Directive_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Directive_Optional_Marshal(n jsn.Marshaler, val **Directive) {
	if *val != nil {
		Directive_Marshal(n, *val)
	}
}

func Directive_Marshal(n jsn.Marshaler, val *Directive) {
	n.MapValues(Directive_Type, Directive_Type)
	n.MapKey("", Directive_Field_Lede)
	/* */ value.Text_Unboxed_Repeats_Marshal(n, &val.Lede)
	n.MapKey("scans", Directive_Field_Scans)
	/* */ ScannerMaker_Repeats_Marshal(n, &val.Scans)
	n.EndValues()
	return
}

// Grammar Read what the player types and turn it into actions.
type Grammar struct {
	Grammar GrammarMaker `if:"label=_"`
}

func (*Grammar) Compose() composer.Spec {
	return composer.Spec{
		Name: Grammar_Type,
		Uses: composer.Type_Flow,
	}
}

const Grammar_Type = "grammar"

const Grammar_Field_Grammar = "$GRAMMAR"

func (op *Grammar) Marshal(n jsn.Marshaler) {
	Grammar_Marshal(n, op)
}

func Grammar_Repeats_Marshal(n jsn.Marshaler, vals *[]Grammar) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Grammar_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Grammar_Optional_Marshal(n jsn.Marshaler, val **Grammar) {
	if *val != nil {
		Grammar_Marshal(n, *val)
	}
}

func Grammar_Marshal(n jsn.Marshaler, val *Grammar) {
	n.MapValues(Grammar_Type, Grammar_Type)
	n.MapKey("", Grammar_Field_Grammar)
	/* */ GrammarMaker_Marshal(n, &val.Grammar)
	n.EndValues()
	return
}

const GrammarMaker_Type = "grammar_maker"

var GrammarMaker_Optional_Marshal = GrammarMaker_Marshal

func GrammarMaker_Marshal(n jsn.Marshaler, ptr *GrammarMaker) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func GrammarMaker_Repeats_Marshal(n jsn.Marshaler, vals *[]GrammarMaker) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			GrammarMaker_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// Noun makes a parser scanner
type Noun struct {
	Kind string `if:"label=_,type=text"`
}

func (*Noun) Compose() composer.Spec {
	return composer.Spec{
		Name: Noun_Type,
		Uses: composer.Type_Flow,
	}
}

const Noun_Type = "noun"

const Noun_Field_Kind = "$KIND"

func (op *Noun) Marshal(n jsn.Marshaler) {
	Noun_Marshal(n, op)
}

func Noun_Repeats_Marshal(n jsn.Marshaler, vals *[]Noun) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Noun_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Noun_Optional_Marshal(n jsn.Marshaler, val **Noun) {
	if *val != nil {
		Noun_Marshal(n, *val)
	}
}

func Noun_Marshal(n jsn.Marshaler, val *Noun) {
	n.MapValues(Noun_Type, Noun_Type)
	n.MapKey("", Noun_Field_Kind)
	/* */ value.Text_Unboxed_Marshal(n, &val.Kind)
	n.EndValues()
	return
}

// Retarget makes a parser scanner
type Retarget struct {
	Span []ScannerMaker `if:"label=_"`
}

func (*Retarget) Compose() composer.Spec {
	return composer.Spec{
		Name: Retarget_Type,
		Uses: composer.Type_Flow,
	}
}

const Retarget_Type = "retarget"

const Retarget_Field_Span = "$SPAN"

func (op *Retarget) Marshal(n jsn.Marshaler) {
	Retarget_Marshal(n, op)
}

func Retarget_Repeats_Marshal(n jsn.Marshaler, vals *[]Retarget) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Retarget_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Retarget_Optional_Marshal(n jsn.Marshaler, val **Retarget) {
	if *val != nil {
		Retarget_Marshal(n, *val)
	}
}

func Retarget_Marshal(n jsn.Marshaler, val *Retarget) {
	n.MapValues(Retarget_Type, Retarget_Type)
	n.MapKey("", Retarget_Field_Span)
	/* */ ScannerMaker_Repeats_Marshal(n, &val.Span)
	n.EndValues()
	return
}

// Reverse makes a parser scanner
type Reverse struct {
	Reverses []ScannerMaker `if:"label=_"`
}

func (*Reverse) Compose() composer.Spec {
	return composer.Spec{
		Name: Reverse_Type,
		Uses: composer.Type_Flow,
	}
}

const Reverse_Type = "reverse"

const Reverse_Field_Reverses = "$REVERSES"

func (op *Reverse) Marshal(n jsn.Marshaler) {
	Reverse_Marshal(n, op)
}

func Reverse_Repeats_Marshal(n jsn.Marshaler, vals *[]Reverse) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Reverse_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Reverse_Optional_Marshal(n jsn.Marshaler, val **Reverse) {
	if *val != nil {
		Reverse_Marshal(n, *val)
	}
}

func Reverse_Marshal(n jsn.Marshaler, val *Reverse) {
	n.MapValues(Reverse_Type, Reverse_Type)
	n.MapKey("", Reverse_Field_Reverses)
	/* */ ScannerMaker_Repeats_Marshal(n, &val.Reverses)
	n.EndValues()
	return
}

const ScannerMaker_Type = "scanner_maker"

var ScannerMaker_Optional_Marshal = ScannerMaker_Marshal

func ScannerMaker_Marshal(n jsn.Marshaler, ptr *ScannerMaker) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func ScannerMaker_Repeats_Marshal(n jsn.Marshaler, vals *[]ScannerMaker) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			ScannerMaker_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

// Self makes a parser scanner which matches the player. ( the player string is just to make the composer happy. )
type Self struct {
	Player string `if:"label=_,type=text"`
}

func (*Self) Compose() composer.Spec {
	return composer.Spec{
		Name: Self_Type,
		Uses: composer.Type_Flow,
	}
}

const Self_Type = "self"

const Self_Field_Player = "$PLAYER"

func (op *Self) Marshal(n jsn.Marshaler) {
	Self_Marshal(n, op)
}

func Self_Repeats_Marshal(n jsn.Marshaler, vals *[]Self) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Self_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Self_Optional_Marshal(n jsn.Marshaler, val **Self) {
	if *val != nil {
		Self_Marshal(n, *val)
	}
}

func Self_Marshal(n jsn.Marshaler, val *Self) {
	n.MapValues(Self_Type, Self_Type)
	n.MapKey("", Self_Field_Player)
	/* */ value.Text_Unboxed_Marshal(n, &val.Player)
	n.EndValues()
	return
}

// Words makes a parser scanner
type Words struct {
	Words []string `if:"label=_,type=text"`
}

func (*Words) Compose() composer.Spec {
	return composer.Spec{
		Name: Words_Type,
		Uses: composer.Type_Flow,
	}
}

const Words_Type = "words"

const Words_Field_Words = "$WORDS"

func (op *Words) Marshal(n jsn.Marshaler) {
	Words_Marshal(n, op)
}

func Words_Repeats_Marshal(n jsn.Marshaler, vals *[]Words) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		n.RepeatValues(cnt)
		for _, el := range *vals {
			Words_Marshal(n, &el)
		}
		n.EndValues()
	}
	return
}

func Words_Optional_Marshal(n jsn.Marshaler, val **Words) {
	if *val != nil {
		Words_Marshal(n, *val)
	}
}

func Words_Marshal(n jsn.Marshaler, val *Words) {
	n.MapValues(Words_Type, Words_Type)
	n.MapKey("", Words_Field_Words)
	/* */ value.Text_Unboxed_Repeats_Marshal(n, &val.Words)
	n.EndValues()
	return
}

var Slots = []interface{}{
	(*GrammarMaker)(nil),
	(*ScannerMaker)(nil),
}

var Slats = []composer.Composer{
	(*Action)(nil),
	(*Alias)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Directive)(nil),
	(*Grammar)(nil),
	(*Noun)(nil),
	(*Retarget)(nil),
	(*Reverse)(nil),
	(*Self)(nil),
	(*Words)(nil),
}
