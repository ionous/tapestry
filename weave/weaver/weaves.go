package weaver

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// when the definition would contradict existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Conflict = errutil.Error("Conflict")

// when the definition would repeat existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Duplicate = errutil.NoPanicError("Duplicate")

// when the definition can't find some required information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Missing = errutil.NoPanicError("Missing")

type Weaves interface {
	// set the possible traits of an aspect
	AddAspectTraits(aspect string, traits []string) error
	AddFact(key string, partsAndValue ...string) error
	AddGrammar(name string, prog *grammar.Directive) error
	AddKind(kind, ancestor string) error
	// fix: shouldn't depend on mdl
	AddKindFields(kind string, fields []FieldInfo) error
	// apply a default trait to a kind
	AddKindTrait(kind, trait string) error
	AddNounKind(noun, kind string) error
	AddNounName(noun, name string, rank int) error
	AddNounPair(rel, noun, otherNoun string) error
	AddNounPath(noun string, path []string, val literal.LiteralValue) error
	AddNounTrait(noun, trait string) error
	AddNounValue(noun, prop string, val rt.Assignment) error
	AddPlural(many, one string) error
	// generates a unique name with a passed category
	GenerateUniqueName(category string) string
}
