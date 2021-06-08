// Code generated by "makeops"; edit at your own risk.
package story

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/dl/reader"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

// AbstractAction requires a user-specified string.
type AbstractAction struct {
	Str string
}

func (op *AbstractAction) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const AbstractAction_Nothing = "$NOTHING"

func (*AbstractAction) Choices() (choices map[string]string) {
	return map[string]string{
		AbstractAction_Nothing: "nothing",
	}
}

func (*AbstractAction) Compose() composer.Spec {
	return composer.Spec{
		Name: "abstract_action",
		Strings: []string{
			"nothing",
		},
	}
}

// ActionContext
type ActionContext struct {
	At   reader.Position `if:"internal"`
	Kind SingularKind    `if:"label=kind"`
}

func (*ActionContext) Compose() composer.Spec {
	return composer.Spec{
		Name: "action_context",
	}
}

// ActionDecl
type ActionDecl struct {
	At           reader.Position `if:"internal"`
	Event        EventName       `if:"label=act"`
	Action       ActionName      `if:"label=acting"`
	ActionParams ActionParams    `if:"label=action_params"`
}

var _ StoryStatement = (*ActionDecl)(nil)

func (*ActionDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "action_decl",
	}
}

// ActionName requires a user-specified string.
type ActionName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *ActionName) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*ActionName) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*ActionName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "action_name",
		OpenStrings: true,
	}
}

// ActionParams swaps between various options
type ActionParams struct {
	Opt interface{}
}

func (*ActionParams) Compose() composer.Spec {
	return composer.Spec{
		Name: "action_params",
	}
}

func (*ActionParams) Choices() map[string]interface{} {
	return map[string]interface{}{
		"common": (*CommonAction)(nil),
		"dual":   (*PairedAction)(nil),
		"none":   (*AbstractAction)(nil),
	}
}

// AreAn requires a user-specified string.
type AreAn struct {
	Str string
}

func (op *AreAn) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const AreAn_Are = "$ARE"
const AreAn_Area = "$AREA"
const AreAn_Arean = "$AREAN"
const AreAn_Is = "$IS"
const AreAn_Isa = "$ISA"
const AreAn_Isan = "$ISAN"

func (*AreAn) Choices() (choices map[string]string) {
	return map[string]string{
		AreAn_Are: "are", AreAn_Area: "area", AreAn_Arean: "arean", AreAn_Is: "is", AreAn_Isa: "isa", AreAn_Isan: "isan",
	}
}

func (*AreAn) Compose() composer.Spec {
	return composer.Spec{
		Name: "are_an",
		Strings: []string{
			"are", "area", "arean", "is", "isa", "isan",
		},
	}
}

// Argument
type Argument struct {
	At   reader.Position `if:"internal"`
	Name value.Text      `if:"label=_"`
	From rt.Assignment   `if:"label=from"`
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Name: "argument",
	}
}

// Arguments
type Arguments struct {
	At   reader.Position `if:"internal"`
	Args []Argument      `if:"label=_"`
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Name: "arguments",
	}
}

// Aspect requires a user-specified string.
type Aspect struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *Aspect) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*Aspect) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*Aspect) Compose() composer.Spec {
	return composer.Spec{
		Name:        "aspect",
		OpenStrings: true,
	}
}

// AspectTraits
type AspectTraits struct {
	Aspect      Aspect      `if:"label=aspect"`
	TraitPhrase TraitPhrase `if:"label=trait_phrase"`
}

var _ StoryStatement = (*AspectTraits)(nil)

func (*AspectTraits) Compose() composer.Spec {
	return composer.Spec{
		Name: "aspect_traits",
	}
}

// BoxedNumber
type BoxedNumber struct {
	Number float64 `if:"label=number"`
}

func (*BoxedNumber) Compose() composer.Spec {
	return composer.Spec{
		Name: "boxed_number",
	}
}

// BoxedText
type BoxedText struct {
	Text value.Text `if:"label=text"`
}

func (*BoxedText) Compose() composer.Spec {
	return composer.Spec{
		Name: "boxed_text",
	}
}

// Certainties
type Certainties struct {
	PluralKinds PluralKinds `if:"label=plural_kinds"`
	AreBeing    bool        `if:"label=are_being"`
	Certainty   Certainty   `if:"label=certainty"`
	Trait       Trait       `if:"label=trait"`
}

var _ StoryStatement = (*Certainties)(nil)

func (*Certainties) Compose() composer.Spec {
	return composer.Spec{
		Name: "certainties",
	}
}

// Certainty requires a user-specified string.
type Certainty struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *Certainty) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const Certainty_Usually = "$USUALLY"
const Certainty_Always = "$ALWAYS"
const Certainty_Seldom = "$SELDOM"
const Certainty_Never = "$NEVER"

func (*Certainty) Choices() (choices map[string]string) {
	return map[string]string{
		Certainty_Usually: "usually", Certainty_Always: "always", Certainty_Seldom: "seldom", Certainty_Never: "never",
	}
}

func (*Certainty) Compose() composer.Spec {
	return composer.Spec{
		Name: "certainty",
		Strings: []string{
			"usually", "always", "seldom", "never",
		},
	}
}

// Comment Information about the story for you and other authors.
type Comment struct {
	Lines value.Lines `if:"label=comment"`
}

var _ StoryStatement = (*Comment)(nil)
var _ rt.Execute = (*Comment)(nil)

func (*Comment) Compose() composer.Spec {
	return composer.Spec{
		Name: "comment",
	}
}

// CommonAction
type CommonAction struct {
	At            reader.Position `if:"internal"`
	Kind          SingularKind    `if:"label=kind"`
	ActionContext *ActionContext  `if:"label=action_context,optional"`
}

func (*CommonAction) Compose() composer.Spec {
	return composer.Spec{
		Name: "common_action",
	}
}

// CountOf A guard which returns true based on a counter.
type CountOf struct {
	At      reader.Position `if:"internal"`
	Trigger core.Trigger    `if:"label=_"`
	Num     rt.NumberEval   `if:"label=num"`
}

var _ rt.BoolEval = (*CountOf)(nil)

func (*CountOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "count_of",
	}
}

// CycleText
type CycleText struct {
	At    reader.Position `if:"internal"`
	Parts []rt.TextEval   `if:"label=_"`
}

var _ rt.TextEval = (*CycleText)(nil)

func (*CycleText) Compose() composer.Spec {
	return composer.Spec{
		Name: "cycle_text",
	}
}

// Determine
type Determine struct {
	Name      value.PatternName `if:"label=_"`
	Arguments *Arguments        `if:"label=args,optional"`
}

var _ rt.Execute = (*Determine)(nil)
var _ rt.BoolEval = (*Determine)(nil)
var _ rt.NumberEval = (*Determine)(nil)
var _ rt.TextEval = (*Determine)(nil)
var _ rt.RecordEval = (*Determine)(nil)
var _ rt.NumListEval = (*Determine)(nil)
var _ rt.TextListEval = (*Determine)(nil)
var _ rt.RecordListEval = (*Determine)(nil)

func (*Determine) Compose() composer.Spec {
	return composer.Spec{
		Name: "determine",
	}
}

// Determiner requires a user-specified string.
type Determiner struct {
	Str string
}

func (op *Determiner) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const Determiner_A = "$A"
const Determiner_An = "$AN"
const Determiner_The = "$THE"
const Determiner_Our = "$OUR"

func (*Determiner) Choices() (choices map[string]string) {
	return map[string]string{
		Determiner_A: "a", Determiner_An: "an", Determiner_The: "the", Determiner_Our: "our",
	}
}

func (*Determiner) Compose() composer.Spec {
	return composer.Spec{
		Name:        "determiner",
		OpenStrings: true,
		Strings: []string{
			"a", "an", "the", "our",
		},
	}
}

// EventBlock Listeners let objects in the game world react to changes before, during, or after they happen.
type EventBlock struct {
	At       reader.Position `if:"internal"`
	Target   EventTarget     `if:"label=the_target"`
	Handlers []EventHandler  `if:"label=handlers"`
}

var _ StoryStatement = (*EventBlock)(nil)

func (*EventBlock) Compose() composer.Spec {
	return composer.Spec{
		Name: "event_block",
	}
}

// EventHandler
type EventHandler struct {
	EventPhase   EventPhase     `if:"label=event_phase"`
	Event        EventName      `if:"label=the_event"`
	Locals       *PatternLocals `if:"label=with_locals,optional"`
	PatternRules PatternRules   `if:"label=pattern_rules"`
}

func (*EventHandler) Compose() composer.Spec {
	return composer.Spec{
		Name: "event_handler",
	}
}

// EventName requires a user-specified string.
type EventName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *EventName) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*EventName) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*EventName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "event_name",
		OpenStrings: true,
	}
}

// EventPhase requires a user-specified string.
type EventPhase struct {
	Str string
}

func (op *EventPhase) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const EventPhase_Before = "$BEFORE"
const EventPhase_While = "$WHILE"
const EventPhase_After = "$AFTER"

func (*EventPhase) Choices() (choices map[string]string) {
	return map[string]string{
		EventPhase_Before: "before", EventPhase_While: "while", EventPhase_After: "after",
	}
}

func (*EventPhase) Compose() composer.Spec {
	return composer.Spec{
		Name: "event_phase",
		Strings: []string{
			"before", "while", "after",
		},
	}
}

// EventTarget swaps between various options
type EventTarget struct {
	Opt interface{}
}

func (*EventTarget) Compose() composer.Spec {
	return composer.Spec{
		Name: "event_target",
	}
}

func (*EventTarget) Choices() map[string]interface{} {
	return map[string]interface{}{
		"kinds":      (*PluralKinds)(nil),
		"named_noun": (*NamedNoun)(nil),
	}
}

// ExtType swaps between various options
type ExtType struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*ExtType) Compose() composer.Spec {
	return composer.Spec{
		Name: "ext_type",
	}
}

func (*ExtType) Choices() map[string]interface{} {
	return map[string]interface{}{
		"numbers":   (*NumberList)(nil),
		"text_list": (*TextList)(nil),
		"record":    (*RecordType)(nil),
		"records":   (*RecordList)(nil),
	}
}

// GrammarDecl
type GrammarDecl struct {
	Grammar grammar.GrammarMaker `if:"label=_"`
}

func (*GrammarDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "grammar_decl",
	}
}

// KindOfNoun
type KindOfNoun struct {
	AreAn        AreAn         `if:"label=are_an"`
	Trait        []Trait       `if:"label=trait,optional"`
	Kind         SingularKind  `if:"label=kind"`
	NounRelation *NounRelation `if:"label=noun_relation,optional"`
}

func (*KindOfNoun) Compose() composer.Spec {
	return composer.Spec{
		Name: "kind_of_noun",
	}
}

// KindOfRelation
type KindOfRelation struct {
	Relation            value.RelationName  `if:"label=relation"`
	RelationCardinality RelationCardinality `if:"label=relation_cardinality"`
}

var _ StoryStatement = (*KindOfRelation)(nil)

func (*KindOfRelation) Compose() composer.Spec {
	return composer.Spec{
		Name: "kind_of_relation",
	}
}

// KindsOfAspect
type KindsOfAspect struct {
	Aspect Aspect `if:"label=aspect"`
}

var _ StoryStatement = (*KindsOfAspect)(nil)

func (*KindsOfAspect) Compose() composer.Spec {
	return composer.Spec{
		Name: "kinds_of_aspect",
	}
}

// KindsOfKind
type KindsOfKind struct {
	PluralKinds  PluralKinds  `if:"label=plural_kinds"`
	SingularKind SingularKind `if:"label=singular_kind"`
}

var _ StoryStatement = (*KindsOfKind)(nil)

func (*KindsOfKind) Compose() composer.Spec {
	return composer.Spec{
		Name: "kinds_of_kind",
	}
}

// KindsOfRecord
type KindsOfRecord struct {
	RecordPlural RecordPlural `if:"label=records"`
}

var _ StoryStatement = (*KindsOfRecord)(nil)

func (*KindsOfRecord) Compose() composer.Spec {
	return composer.Spec{
		Name: "kinds_of_record",
	}
}

// KindsPossessProperties
type KindsPossessProperties struct {
	PluralKinds  PluralKinds    `if:"label=plural_kinds"`
	PropertyDecl []PropertyDecl `if:"label=property_decl"`
}

var _ StoryStatement = (*KindsPossessProperties)(nil)

func (*KindsPossessProperties) Compose() composer.Spec {
	return composer.Spec{
		Name: "kinds_possess_properties",
	}
}

// Lede Describes one or more nouns.
type Lede struct {
	Nouns      []NamedNoun `if:"label=nouns"`
	NounPhrase NounPhrase  `if:"label=noun_phrase"`
}

func (*Lede) Compose() composer.Spec {
	return composer.Spec{
		Name: "lede",
	}
}

// LocalDecl
type LocalDecl struct {
	VariableDecl VariableDecl `if:"label=variable_decl"`
	Value        *LocalInit   `if:"label=starting_as,optional"`
}

func (*LocalDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "local_decl",
	}
}

// LocalInit
type LocalInit struct {
	Value rt.Assignment `if:"label=value"`
}

func (*LocalInit) Compose() composer.Spec {
	return composer.Spec{
		Name: "local_init",
	}
}

// Make
type Make struct {
	Name      value.Text `if:"label=_"`
	Arguments *Arguments `if:"label=args,optional"`
}

var _ rt.RecordEval = (*Make)(nil)

func (*Make) Compose() composer.Spec {
	return composer.Spec{
		Name: "make",
	}
}

// ManyToMany
type ManyToMany struct {
	Kinds      PluralKinds `if:"label=kinds"`
	OtherKinds PluralKinds `if:"label=other_kinds"`
}

func (*ManyToMany) Compose() composer.Spec {
	return composer.Spec{
		Name: "many_to_many",
	}
}

// ManyToOne
type ManyToOne struct {
	Kinds PluralKinds  `if:"label=kinds"`
	Kind  SingularKind `if:"label=kind"`
}

func (*ManyToOne) Compose() composer.Spec {
	return composer.Spec{
		Name: "many_to_one",
	}
}

// NamedNoun
type NamedNoun struct {
	Determiner Determiner `if:"label=determiner"`
	Name       NounName   `if:"label=name"`
}

func (*NamedNoun) Compose() composer.Spec {
	return composer.Spec{
		Name: "named_noun",
	}
}

// NounAssignment Assign text.
type NounAssignment struct {
	Property Property    `if:"label=property"`
	Nouns    []NamedNoun `if:"label=nouns"`
	Lines    value.Lines `if:"label=the_text"`
}

var _ StoryStatement = (*NounAssignment)(nil)

func (*NounAssignment) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun_assignment",
	}
}

// NounName requires a user-specified string.
type NounName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *NounName) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*NounName) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*NounName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "noun_name",
		OpenStrings: true,
	}
}

// NounPhrase swaps between various options
type NounPhrase struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*NounPhrase) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun_phrase",
	}
}

func (*NounPhrase) Choices() map[string]interface{} {
	return map[string]interface{}{
		"kind_of_noun":  (*KindOfNoun)(nil),
		"noun_traits":   (*NounTraits)(nil),
		"noun_relation": (*NounRelation)(nil),
	}
}

// NounRelation
type NounRelation struct {
	AreBeing bool               `if:"label=are_being,optional"`
	Relation value.RelationName `if:"label=relation"`
	Nouns    []NamedNoun        `if:"label=nouns"`
}

func (*NounRelation) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun_relation",
	}
}

// NounStatement Describes people, places, or things.
type NounStatement struct {
	Lede    Lede     `if:"label=lede"`
	Tail    []Tail   `if:"label=tail,optional"`
	Summary *Summary `if:"label=summary,optional"`
}

var _ StoryStatement = (*NounStatement)(nil)

func (*NounStatement) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun_statement",
	}
}

// NounTraits
type NounTraits struct {
	AreBeing bool    `if:"label=are_being"`
	Trait    []Trait `if:"label=trait"`
}

func (*NounTraits) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun_traits",
	}
}

// NumberList requires a user-specified string.
type NumberList struct {
	Str string
}

func (op *NumberList) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const NumberList_List = "$LIST"

func (*NumberList) Choices() (choices map[string]string) {
	return map[string]string{
		NumberList_List: "list",
	}
}

func (*NumberList) Compose() composer.Spec {
	return composer.Spec{
		Name: "number_list",
		Strings: []string{
			"list",
		},
	}
}

// ObjectType
type ObjectType struct {
	An   bool         `if:"label=an"`
	Kind SingularKind `if:"label=kind_of"`
}

func (*ObjectType) Compose() composer.Spec {
	return composer.Spec{
		Name: "object_type",
	}
}

// OneToMany
type OneToMany struct {
	Kind  SingularKind `if:"label=kind"`
	Kinds PluralKinds  `if:"label=kinds"`
}

func (*OneToMany) Compose() composer.Spec {
	return composer.Spec{
		Name: "one_to_many",
	}
}

// OneToOne
type OneToOne struct {
	Kind      SingularKind `if:"label=kind"`
	OtherKind SingularKind `if:"label=other_kind"`
}

func (*OneToOne) Compose() composer.Spec {
	return composer.Spec{
		Name: "one_to_one",
	}
}

// PairedAction
type PairedAction struct {
	At    reader.Position `if:"internal"`
	Kinds PluralKinds     `if:"label=kinds"`
}

func (*PairedAction) Compose() composer.Spec {
	return composer.Spec{
		Name: "paired_action",
	}
}

// Paragraph
type Paragraph struct {
	StoryStatement []StoryStatement `if:"label=story_statement,optional"`
}

func (*Paragraph) Compose() composer.Spec {
	return composer.Spec{
		Name: "paragraph",
	}
}

// PatternActions Actions to take when using a pattern.
type PatternActions struct {
	Name          value.PatternName `if:"label=pattern_name"`
	PatternLocals *PatternLocals    `if:"label=pattern_locals,optional"`
	PatternReturn *PatternReturn    `if:"label=pattern_return,optional"`
	PatternRules  PatternRules      `if:"label=pattern_rules"`
}

var _ StoryStatement = (*PatternActions)(nil)

func (*PatternActions) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_actions",
	}
}

// PatternDecl
type PatternDecl struct {
	Type          PatternType           `if:"label=type"`
	Name          value.PatternName     `if:"label=name"`
	Optvars       *PatternVariablesTail `if:"label=parameters,optional"`
	PatternReturn *PatternReturn        `if:"label=pattern_return,optional"`
	About         *Comment              `if:"label=about,optional"`
}

var _ StoryStatement = (*PatternDecl)(nil)

func (*PatternDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_decl",
	}
}

// PatternFlags requires a user-specified string.
type PatternFlags struct {
	Str string
}

func (op *PatternFlags) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const PatternFlags_Before = "$BEFORE"
const PatternFlags_After = "$AFTER"
const PatternFlags_Terminate = "$TERMINATE"

func (*PatternFlags) Choices() (choices map[string]string) {
	return map[string]string{
		PatternFlags_Before: "before", PatternFlags_After: "after", PatternFlags_Terminate: "terminate",
	}
}

func (*PatternFlags) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_flags",
		Strings: []string{
			"before", "after", "terminate",
		},
	}
}

// PatternLocals
type PatternLocals struct {
	LocalDecl []LocalDecl `if:"label=local_decl"`
}

func (*PatternLocals) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_locals",
	}
}

// PatternReturn
type PatternReturn struct {
	Result VariableDecl `if:"label=result"`
}

func (*PatternReturn) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_return",
	}
}

// PatternRule
type PatternRule struct {
	Guard rt.BoolEval  `if:"label=conditions_are met"`
	Flags PatternFlags `if:"label=continue,optional"`
	Hook  ProgramHook  `if:"label=do"`
}

func (*PatternRule) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_rule",
	}
}

// PatternRules
type PatternRules struct {
	PatternRule []PatternRule `if:"label=pattern_rule,optional"`
}

func (*PatternRules) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_rules",
	}
}

// PatternType requires a user-specified string.
type PatternType struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PatternType) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const PatternType_Patterns = "$PATTERNS"
const PatternType_Actions = "$ACTIONS"
const PatternType_Events = "$EVENTS"

func (*PatternType) Choices() (choices map[string]string) {
	return map[string]string{
		PatternType_Patterns: "patterns", PatternType_Actions: "actions", PatternType_Events: "events",
	}
}

func (*PatternType) Compose() composer.Spec {
	return composer.Spec{
		Name:        "pattern_type",
		OpenStrings: true,
		Strings: []string{
			"patterns", "actions", "events",
		},
	}
}

// PatternVariablesDecl Values provided when calling a pattern.
type PatternVariablesDecl struct {
	PatternName  value.PatternName `if:"label=pattern_name"`
	VariableDecl []VariableDecl    `if:"label=variable_decl"`
}

var _ StoryStatement = (*PatternVariablesDecl)(nil)

func (*PatternVariablesDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_variables_decl",
	}
}

// PatternVariablesTail Storage for values used during the execution of a pattern.
type PatternVariablesTail struct {
	VariableDecl []VariableDecl `if:"label=variable_decl"`
}

func (*PatternVariablesTail) Compose() composer.Spec {
	return composer.Spec{
		Name: "pattern_variables_tail",
	}
}

// PluralKinds requires a user-specified string.
type PluralKinds struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PluralKinds) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*PluralKinds) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*PluralKinds) Compose() composer.Spec {
	return composer.Spec{
		Name:        "plural_kinds",
		OpenStrings: true,
	}
}

// PrimitiveType requires a user-specified string.
type PrimitiveType struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PrimitiveType) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const PrimitiveType_Number = "$NUMBER"
const PrimitiveType_Text = "$TEXT"
const PrimitiveType_Bool = "$BOOL"

func (*PrimitiveType) Choices() (choices map[string]string) {
	return map[string]string{
		PrimitiveType_Number: "number", PrimitiveType_Text: "text", PrimitiveType_Bool: "bool",
	}
}

func (*PrimitiveType) Compose() composer.Spec {
	return composer.Spec{
		Name: "primitive_type",
		Strings: []string{
			"number", "text", "bool",
		},
	}
}

// PrimitiveValue swaps between various options
type PrimitiveValue struct {
	Opt interface{}
}

func (*PrimitiveValue) Compose() composer.Spec {
	return composer.Spec{
		Name: "primitive_value",
	}
}

func (*PrimitiveValue) Choices() map[string]interface{} {
	return map[string]interface{}{
		"boxed_text":   (*BoxedText)(nil),
		"boxed_number": (*BoxedNumber)(nil),
	}
}

// ProgramHook swaps between various options
type ProgramHook struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*ProgramHook) Compose() composer.Spec {
	return composer.Spec{
		Name: "program_hook",
	}
}

func (*ProgramHook) Choices() map[string]interface{} {
	return map[string]interface{}{
		"activity": (*core.Activity)(nil),
	}
}

// Pronoun requires a user-specified string.
type Pronoun struct {
	Str string
}

func (op *Pronoun) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const Pronoun_It = "$IT"
const Pronoun_They = "$THEY"

func (*Pronoun) Choices() (choices map[string]string) {
	return map[string]string{
		Pronoun_It: "it", Pronoun_They: "they",
	}
}

func (*Pronoun) Compose() composer.Spec {
	return composer.Spec{
		Name:        "pronoun",
		OpenStrings: true,
		Strings: []string{
			"it", "they",
		},
	}
}

// Property requires a user-specified string.
type Property struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *Property) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*Property) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*Property) Compose() composer.Spec {
	return composer.Spec{
		Name:        "property",
		OpenStrings: true,
	}
}

// PropertyAspect requires a user-specified string.
type PropertyAspect struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PropertyAspect) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const PropertyAspect_Aspect = "$ASPECT"

func (*PropertyAspect) Choices() (choices map[string]string) {
	return map[string]string{
		PropertyAspect_Aspect: "aspect",
	}
}

func (*PropertyAspect) Compose() composer.Spec {
	return composer.Spec{
		Name: "property_aspect",
		Strings: []string{
			"aspect",
		},
	}
}

// PropertyDecl
type PropertyDecl struct {
	An           Determiner   `if:"label=an"`
	Property     Property     `if:"label=property"`
	PropertyType PropertyType `if:"label=property_type"`
	Comment      value.Lines  `if:"label=comment,optional"`
}

func (*PropertyDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "property_decl",
	}
}

// PropertyType swaps between various options
type PropertyType struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*PropertyType) Compose() composer.Spec {
	return composer.Spec{
		Name: "property_type",
	}
}

func (*PropertyType) Choices() map[string]interface{} {
	return map[string]interface{}{
		"property_aspect": (*PropertyAspect)(nil),
		"primitive":       (*PrimitiveType)(nil),
		"ext":             (*ExtType)(nil),
	}
}

// RecordList
type RecordList struct {
	Kind RecordSingular `if:"label=kind"`
}

func (*RecordList) Compose() composer.Spec {
	return composer.Spec{
		Name: "record_list",
	}
}

// RecordPlural requires a user-specified string.
type RecordPlural struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *RecordPlural) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*RecordPlural) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*RecordPlural) Compose() composer.Spec {
	return composer.Spec{
		Name:        "record_plural",
		OpenStrings: true,
	}
}

// RecordSingular requires a user-specified string.
type RecordSingular struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *RecordSingular) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*RecordSingular) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*RecordSingular) Compose() composer.Spec {
	return composer.Spec{
		Name:        "record_singular",
		OpenStrings: true,
	}
}

// RecordType
type RecordType struct {
	Kind RecordSingular `if:"label=kind"`
}

func (*RecordType) Compose() composer.Spec {
	return composer.Spec{
		Name: "record_type",
	}
}

// RecordsPossessProperties
type RecordsPossessProperties struct {
	RecordPlural RecordPlural   `if:"label=records"`
	PropertyDecl []PropertyDecl `if:"label=property_decl"`
}

var _ StoryStatement = (*RecordsPossessProperties)(nil)

func (*RecordsPossessProperties) Compose() composer.Spec {
	return composer.Spec{
		Name: "records_possess_properties",
	}
}

// RelationCardinality swaps between various options
type RelationCardinality struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*RelationCardinality) Compose() composer.Spec {
	return composer.Spec{
		Name: "relation_cardinality",
	}
}

func (*RelationCardinality) Choices() map[string]interface{} {
	return map[string]interface{}{
		"one_to_one":   (*OneToOne)(nil),
		"one_to_many":  (*OneToMany)(nil),
		"many_to_one":  (*ManyToOne)(nil),
		"many_to_many": (*ManyToMany)(nil),
	}
}

// RelativeToNoun
type RelativeToNoun struct {
	Relation value.RelationName `if:"label=relation"`
	Nouns    []NamedNoun        `if:"label=nouns"`
	AreBeing bool               `if:"label=are_being"`
	Nouns1   []NamedNoun        `if:"label=nouns"`
}

var _ StoryStatement = (*RelativeToNoun)(nil)

func (*RelativeToNoun) Compose() composer.Spec {
	return composer.Spec{
		Name: "relative_to_noun",
	}
}

// RenderTemplate Parse text using iffy templates.
type RenderTemplate struct {
	Template value.Lines `if:"label=_"`
}

var _ rt.TextEval = (*RenderTemplate)(nil)

func (*RenderTemplate) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_template",
	}
}

// Send Triggers a event, returns a true/false success value.
type Send struct {
	Event     value.Text      `if:"label=_"`
	Path      rt.TextListEval `if:"label=to"`
	Arguments *Arguments      `if:"label=args,optional"`
}

var _ rt.Execute = (*Send)(nil)
var _ rt.BoolEval = (*Send)(nil)

func (*Send) Compose() composer.Spec {
	return composer.Spec{
		Name: "send",
	}
}

// ShuffleText
type ShuffleText struct {
	At    reader.Position `if:"internal"`
	Parts []rt.TextEval   `if:"label=_"`
}

var _ rt.TextEval = (*ShuffleText)(nil)

func (*ShuffleText) Compose() composer.Spec {
	return composer.Spec{
		Name: "shuffle_text",
	}
}

// SingularKind requires a user-specified string.
type SingularKind struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *SingularKind) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*SingularKind) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*SingularKind) Compose() composer.Spec {
	return composer.Spec{
		Name:        "singular_kind",
		OpenStrings: true,
	}
}

// StoppingText
type StoppingText struct {
	At    reader.Position `if:"internal"`
	Parts []rt.TextEval   `if:"label=_"`
}

var _ rt.TextEval = (*StoppingText)(nil)

func (*StoppingText) Compose() composer.Spec {
	return composer.Spec{
		Name: "stopping_text",
	}
}

// Story
type Story struct {
	Paragraph []Paragraph `if:"label=paragraph,optional"`
}

func (*Story) Compose() composer.Spec {
	return composer.Spec{
		Name: "story",
	}
}

// Summary
type Summary struct {
	At    reader.Position `if:"internal"`
	Lines value.Lines     `if:"label=summary"`
}

func (*Summary) Compose() composer.Spec {
	return composer.Spec{
		Name: "summary",
	}
}

// Tail Adds details about the preceding noun or nouns.
type Tail struct {
	Pronoun    Pronoun    `if:"label=pronoun"`
	NounPhrase NounPhrase `if:"label=noun_phrase"`
}

func (*Tail) Compose() composer.Spec {
	return composer.Spec{
		Name: "tail",
	}
}

// TestName requires a user-specified string.
type TestName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *TestName) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const TestName_CurrentTest = "$CURRENT_TEST"

func (*TestName) Choices() (choices map[string]string) {
	return map[string]string{
		TestName_CurrentTest: "current_test",
	}
}

func (*TestName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "test_name",
		OpenStrings: true,
		Strings: []string{
			"current_test",
		},
	}
}

// TestOutput Expect that a test uses &#x27;Say&#x27; to print some specific text.
type TestOutput struct {
	Lines value.Lines `if:"label=lines"`
}

var _ Testing = (*TestOutput)(nil)

func (*TestOutput) Compose() composer.Spec {
	return composer.Spec{
		Name: "test_output",
	}
}

// TestRule
type TestRule struct {
	TestName TestName    `if:"label=test_name"`
	Hook     ProgramHook `if:"label=do"`
}

var _ StoryStatement = (*TestRule)(nil)

func (*TestRule) Compose() composer.Spec {
	return composer.Spec{
		Name: "test_rule",
	}
}

// TestScene
type TestScene struct {
	TestName TestName `if:"label=test_name"`
	Story    Story    `if:"label=story"`
}

var _ StoryStatement = (*TestScene)(nil)

func (*TestScene) Compose() composer.Spec {
	return composer.Spec{
		Name: "test_scene",
	}
}

// TestStatement
type TestStatement struct {
	At       reader.Position `if:"internal"`
	TestName TestName        `if:"label=test_name"`
	Test     Testing         `if:"label=expectation"`
}

var _ StoryStatement = (*TestStatement)(nil)

func (*TestStatement) Compose() composer.Spec {
	return composer.Spec{
		Name: "test_statement",
	}
}

// TextList requires a user-specified string.
type TextList struct {
	Str string
}

func (op *TextList) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

const TextList_List = "$LIST"

func (*TextList) Choices() (choices map[string]string) {
	return map[string]string{
		TextList_List: "list",
	}
}

func (*TextList) Compose() composer.Spec {
	return composer.Spec{
		Name: "text_list",
		Strings: []string{
			"list",
		},
	}
}

// Trait requires a user-specified string.
type Trait struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *Trait) String() (ret string) {
	if s := op.Str; s != "$EMPTY" {
		ret = s
	}
	return
}

func (*Trait) Choices() (choices map[string]string) {
	return map[string]string{}
}

func (*Trait) Compose() composer.Spec {
	return composer.Spec{
		Name:        "trait",
		OpenStrings: true,
	}
}

// TraitPhrase
type TraitPhrase struct {
	AreEither bool    `if:"label=are_either"`
	Trait     []Trait `if:"label=trait"`
}

func (*TraitPhrase) Compose() composer.Spec {
	return composer.Spec{
		Name: "trait_phrase",
	}
}

// VariableDecl
type VariableDecl struct {
	An      Determiner         `if:"label=an"`
	Name    value.VariableName `if:"label=name"`
	Type    VariableType       `if:"label=type"`
	Comment value.Lines        `if:"label=comment,optional"`
}

func (*VariableDecl) Compose() composer.Spec {
	return composer.Spec{
		Name: "variable_decl",
	}
}

// VariableType swaps between various options
type VariableType struct {
	At  reader.Position `if:"internal"`
	Opt interface{}
}

func (*VariableType) Compose() composer.Spec {
	return composer.Spec{
		Name: "variable_type",
	}
}

func (*VariableType) Choices() map[string]interface{} {
	return map[string]interface{}{
		"primitive": (*PrimitiveType)(nil),
		"object":    (*ObjectType)(nil),
		"ext":       (*ExtType)(nil),
	}
}

var Slots = []interface{}{
	(*StoryStatement)(nil),
	(*Testing)(nil),
}
var Swaps = []interface{}{
	(*ActionParams)(nil),
	(*EventTarget)(nil),
	(*ExtType)(nil),
	(*NounPhrase)(nil),
	(*PrimitiveValue)(nil),
	(*ProgramHook)(nil),
	(*PropertyType)(nil),
	(*RelationCardinality)(nil),
	(*VariableType)(nil),
}
var Slats = []composer.Composer{
	(*ActionContext)(nil),
	(*ActionDecl)(nil),
	(*Argument)(nil),
	(*Arguments)(nil),
	(*AspectTraits)(nil),
	(*BoxedNumber)(nil),
	(*BoxedText)(nil),
	(*Certainties)(nil),
	(*Comment)(nil),
	(*CommonAction)(nil),
	(*CountOf)(nil),
	(*CycleText)(nil),
	(*Determine)(nil),
	(*EventBlock)(nil),
	(*EventHandler)(nil),
	(*GrammarDecl)(nil),
	(*KindOfNoun)(nil),
	(*KindOfRelation)(nil),
	(*KindsOfAspect)(nil),
	(*KindsOfKind)(nil),
	(*KindsOfRecord)(nil),
	(*KindsPossessProperties)(nil),
	(*Lede)(nil),
	(*LocalDecl)(nil),
	(*LocalInit)(nil),
	(*Make)(nil),
	(*ManyToMany)(nil),
	(*ManyToOne)(nil),
	(*NamedNoun)(nil),
	(*NounAssignment)(nil),
	(*NounRelation)(nil),
	(*NounStatement)(nil),
	(*NounTraits)(nil),
	(*ObjectType)(nil),
	(*OneToMany)(nil),
	(*OneToOne)(nil),
	(*PairedAction)(nil),
	(*Paragraph)(nil),
	(*PatternActions)(nil),
	(*PatternDecl)(nil),
	(*PatternLocals)(nil),
	(*PatternReturn)(nil),
	(*PatternRule)(nil),
	(*PatternRules)(nil),
	(*PatternVariablesDecl)(nil),
	(*PatternVariablesTail)(nil),
	(*PropertyDecl)(nil),
	(*RecordList)(nil),
	(*RecordType)(nil),
	(*RecordsPossessProperties)(nil),
	(*RelativeToNoun)(nil),
	(*RenderTemplate)(nil),
	(*Send)(nil),
	(*ShuffleText)(nil),
	(*StoppingText)(nil),
	(*Story)(nil),
	(*Summary)(nil),
	(*Tail)(nil),
	(*TestOutput)(nil),
	(*TestRule)(nil),
	(*TestScene)(nil),
	(*TestStatement)(nil),
	(*TraitPhrase)(nil),
	(*VariableDecl)(nil),
}
