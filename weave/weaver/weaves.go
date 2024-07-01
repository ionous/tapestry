package weaver

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// when the definition would contradict existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
var ErrConflict = mdl.ErrConflict

// when the definition would repeat existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
var ErrDuplicate = mdl.ErrDuplicate

// when the definition can't find some required information:
// the returned error wraps this tag. errors.Is can be used to detect it.
var ErrMissing = mdl.ErrMissing

type Weaves interface {
	// set the possible traits of an aspect
	AddAspectTraits(aspect string, traits []string) error
	AddCheck(name string, value literal.LiteralValue, prog []rt.Execute) error
	AddFact(key string, partsAndValue ...string) error
	AddGrammar(name string, prog *grammar.Directive) error
	AddKind(kind, ancestor string) error
	AddKindFields(kind string, fields []mdl.FieldInfo) error
	// apply a default trait to a kind
	AddKindTrait(kind, trait string) error
	AddNounKind(noun, kind string) error
	AddNounName(noun, name string, rank int) error
	AddNounPair(rel, noun, otherNoun string) error
	AddNounPath(noun string, path []string, val literal.LiteralValue) error
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error
	AddPattern(mdl.Pattern) error
	AddPlural(many, one string) error
	AddRelation(name, oneKind, otherKind string, amany bool, bmany bool) error
	ExtendPattern(mdl.Pattern) error
	// generates a unique name with a passed category
	GenerateUniqueName(category string) string
}
