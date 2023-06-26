package eph

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type Aliases struct {
	ShortName string   `if:"label=understand,type=text"`
	Aliases   []string `if:"label=as,type=text"`
}

func (op *Aliases) Assert(k assert.Assertions) (err error) {
	return k.AssertAlias(op.ShortName, op.Aliases...)
}

// Aspects A set of related object states such that exactly one member of the set is true for a given object at a single time.
// Generates an implicit kind of 'aspect' where every field of the kind is a boolean property.
type Aspects struct {
	Aspects string   `if:"label=aspects,type=text"`
	Traits  []string `if:"label=traits,type=text"`
}

func (op *Aspects) Assert(k assert.Assertions) (err error) {
	return k.AssertAspectTraits(op.Aspects, op.Traits)
}

type BeginDomain struct {
	Name     string   `if:"label=domain,type=text"`
	Requires []string `if:"label=requires,type=text"`
}

func (op *BeginDomain) Assert(k assert.Assertions) (err error) {
	return k.AssertDomainStart(op.Name, op.Requires)
}

// Checks
// type Checks struct {
// 	Name   string               `if:"label=check,type=text"`
// 	Expect literal.LiteralValue `if:"label=expect,optional"`
// 	Exe    []rt.Execute         `if:"label=does"`
// }

// func (op *Checks) Assert(k assert.Assertions) (err error) {
// 	return k.AssertCheck(op.Name, op.Exe, op.Expect)
// }

// type Definition struct {
// 	Path  []string
// 	Value string
// }

// func (op *Definition) Assert(k assert.Assertions) (err error) {
// 	path := append(op.Path, op.Value)
// 	return k.AssertDefinition(path...)
// }

// Directives
type Directives struct {
	Name      string            `if:"label=go,type=text"`
	Directive grammar.Directive `if:"label=parse"`
}

func (op *Directives) Assert(k assert.Assertions) (err error) {
	return k.AssertGrammar(op.Name, &op.Directive)
}

// EndDomain
type EndDomain struct {
	Name string `if:"label=domain,type=text"`
}

func (op *EndDomain) Assert(k assert.Assertions) (err error) {
	return k.AssertDomainEnd()
}

// Kinds A new type deriving from another existing type.
// The new kind has all of the properties of all of its ancestor kinds
// and it can be used wherever one of its ancestor kinds is needed.
// ( The reverse isn't true because the new kind can have its own unique properties not available to its ancestors. )
type Kinds struct {
	Kind     string   `if:"label=kinds,type=text"`
	Ancestor string   `if:"label=from,type=text"`
	Contain  []Params `if:"label=contain"`
}

func (op *Kinds) Assert(k assert.Assertions) (err error) {
	ps := op.Contain
	if len(ps) > 0 {
		err = assertFields(op.Kind, ps, k.AssertField)
	}
	// tbd: are we supposed to be able to declare a kind and its fields in one fell swoop?
	// the tests seem to indicate no.
	// except for TestValuePaths which seemed to indicate yes:
	// ( i had to add an explicit Kind )
	// doesn't really matter since eph is only used for testing.
	if len(ps) == 0 || len(op.Ancestor) > 0 {
		err = k.AssertAncestor(op.Kind, op.Ancestor)
	}
	return
}

// EphMacro - hijacks pattern registration for use with macros
// type Macro struct {
// 	Patterns
// 	MacroStatements []rt.Execute
// }

// Nouns
type Nouns struct {
	Noun string `if:"label=noun,type=text"`
	Kind string `if:"label=kind,type=text"`
}

func (op *Nouns) Assert(k assert.Assertions) (err error) {
	return k.AssertNounKind(op.Noun, op.Kind)
}

// Opposites Rules for transforming plural text to singular text and back again.
// Used by the assembler to help interpret author definitions,
// and at runtime to help the parser interpret user input.
type Opposites struct {
	Opposite string `if:"label=opposite,type=text"`
	Word     string `if:"label=word,type=text"`
}

func (op *Opposites) Assert(k assert.Assertions) (err error) {
	return k.AssertOpposite(op.Opposite, op.Word)
}

// Patterns Patterns provide author reusable code.
// The parameters define values provided by the caller.
// Locals provide scratch values for use during pattern processing.
// The result allows the pattern to return a value to the caller of pattern.
// While multiple pattern commands can be used to define a pattern,
// the set of arguments and the return can only be specified once.
type Patterns struct {
	PatternName string   `if:"label=pattern,type=text"`
	Params      []Params `if:"label=with,optional"`
	Locals      []Params `if:"label=locals,optional"`
	Result      *Params  `if:"label=result,optional"`
}

func (op *Patterns) Assert(k assert.Assertions) (err error) {
	kind := op.PatternName
	if e := k.AssertAncestor(kind, kindsOf.Pattern.String()); e != nil {
		err = e
	} else {
		if ps := op.Params; err == nil && len(ps) > 0 {
			err = assertFields(kind, ps, k.AssertParam)
		}
		if p := op.Result; err == nil && p != nil {
			err = k.AssertResult(kind, p.Name, p.Class, p.Affinity, p.Initially)
		}
		if ps := op.Locals; err == nil {
			err = assertFields(kind, ps, k.AssertField)
		}
	}
	return
}

// Plurals Rules for transforming plural text to singular text and back again.
// Used by the assembler to help interpret author definitions,
// and at runtime to help the parser interpret user input.
type Plurals struct {
	Plural   string `if:"label=plural,type=text"`
	Singular string `if:"label=singular,type=text"`
}

func (op *Plurals) Assert(k assert.Assertions) (err error) {
	return k.AssertPlural(op.Singular, op.Plural)
}

// Refs Implies some fact about the world that will be defined elsewhere.
// Reuses the set of ephemera to limit redefinition. Not all are valid.
// type Refs struct {
// 	Refs []Ephemera `if:"label=refs"`
// }

// func (op *Refs) Assert(k assert.Assertions) (err error) {
// 	refsNotImplemented.PrintOnce()
// 	return
// }

// // refs imply some fact about the world that will be defined elsewhere.
// // assembly would verify that the referenced thing really exists
// var refsNotImplemented PrintOnce = "refs not implemented"

// Relations
type Relations struct {
	Rel         string      `if:"label=_,type=text"`
	Cardinality Cardinality `if:"label=relate"`
}

func (op *Relations) Assert(k assert.Assertions) (err error) {
	switch c := op.Cardinality.(type) {
	case *OneOne:
		err = k.AssertRelation(op.Rel, c.Kind, c.OtherKind, false, false)
	case *OneMany:
		err = k.AssertRelation(op.Rel, c.Kind, c.OtherKinds, false, true)
	case *ManyOne:
		err = k.AssertRelation(op.Rel, c.Kinds, c.OtherKind, true, false)
	case *ManyMany:
		err = k.AssertRelation(op.Rel, c.Kinds, c.OtherKinds, true, true)
	}
	return
}

// Relatives
type Relatives struct {
	Rel       string `if:"label=_,type=text"`
	Noun      string `if:"label=relates,type=text"`
	OtherNoun string `if:"label=to,type=text"`
}

func (op *Relatives) Assert(k assert.Assertions) (err error) {
	return k.AssertRelative(op.Rel, op.Noun, op.OtherNoun)
}

// Rules
type Rules struct {
	PatternName string       `if:"label=pattern,type=text"`
	Target      string       `if:"label=target,optional,type=text"`
	Filter      rt.BoolEval  `if:"label=if"`
	When        Timing       `if:"label=when"`
	Exe         []rt.Execute `if:"label=does"`
	Touch       Always       `if:"label=touch,optional"`
}

func (op *Rules) Assert(k assert.Assertions) (err error) {
	flags := toTiming(op.When, op.Touch)
	return k.AssertRule(op.PatternName, op.Target, op.Filter, flags, op.Exe)
}

// Values Give a noun a specific value at startup.
// Initialization is somewhat simplistic:
// 1. Initial values are not scoped to domains, triggers must be used to change values when domains begin and end.
// 2. The values inside of records can be set using a 'path' to find them, however individual values within lists cannot be set.
// Note: when using a path, the path addresses the noun first, the named field - referring to the inner most record - last.
type Values struct {
	Noun  string               `if:"label=noun,type=text"`
	Field string               `if:"label=has,type=text"`
	Path  []string             `if:"label=path,optional,type=text"`
	Value literal.LiteralValue `if:"label=value"`
}

func (op *Values) Assert(k assert.Assertions) (err error) {
	return k.AssertNounValue(op.Noun, op.Field, op.Path, op.Value)
}
