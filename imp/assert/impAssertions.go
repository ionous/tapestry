package assert

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

type EventTiming int

const (
	//
	DefaultTiming EventTiming = iota
	Before
	During
	After
	Later     // 100
	RunAlways = 1 << 3
)

type Phase int

//go:generate stringer -type=Phase
const (
	DomainPhase Phase = iota
	PluralPhase
	AncestryPhase
	PropertyPhase // collect the properties of kinds
	AspectPhase   // traits of kinds
	FieldPhase    // actually assemble those fields
	NounPhase     // instances ( of kinds )
	ValuePhase
	RelativePhase // initial relations between nouns
	PatternPhase
	AliasPhase
	DirectivePhase // more grammar
	RefPhase
	NumPhases
)

// assert.Facts?
// fix: more gl like?
// in that case would ancestor be primary or kind?
// and would domain requires be declarative or a series of push?
type Assertions interface {
	BeginDomain(name string, requires []string) error
	EndDomain() error

	AssertAlias(short string, names ...string) error
	AssertAncestor(kind, ancestor string) error
	AssertAspectTraits(aspect string, traits []string) error

	AssertField(kind, name, class string, aff affine.Affinity, init assign.Assignment) error
	AssertLocal(kind, name, class string, aff affine.Affinity, init assign.Assignment) error
	AssertParam(kind, name, class string, aff affine.Affinity, init assign.Assignment) error
	AssertResult(kind, name, class string, aff affine.Affinity, init assign.Assignment) error

	AssertGrammar(name string, directive *grammar.Directive) error
	// AssertMacro(a, b string) error
	AssertNounKind(noun, kind string) error
	AssertNounPhrase() error
	// fix should val be assign?
	// should path be a DottedPath?
	AssertNounValue(nounName, fieldName string, path []string, val literal.LiteralValue) error
	AssertOpposite(a, b string) error
	AssertPlural(singluar, plural string) error
	// AssertRef(a, b string) error
	AssertRelation(rel, a, b string, amany, bmany bool) error
	AssertRelative(rel, noun, otherNoun string) error
	// can this be a rule or something?
	AssertCheck(name string, do []rt.Execute) error
	// fix: target should become part of the guard.
	// and/or rule should be wrapped up more like "grammar.Directive"
	AssertRule(name string, target string, guard rt.BoolEval, flags EventTiming, do []rt.Execute) error

	AssertDefinition(path ...string) error

	// fix: this is i think the only reason Phase is in assert instead of ... say... weave
	// is schedule necessary? [ can it be supplanted by commands, even faux ones? ]
	Schedule(Phase, func(World, Assertions) error) error
}

// fix: this should eventually be a runtime if at all possible
type World interface {
	PluralOf(single string) string
	SingularOf(plural string) string
	OppositeOf(string) string
}

// // helper: fix: maybe move to story -- make part of importer?
func AssertNounValue(a Assertions, val literal.LiteralValue, noun string, path ...string) error {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return a.AssertNounValue(noun, field, parts, val)
}
