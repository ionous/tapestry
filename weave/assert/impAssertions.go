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
	RequireDependencies Phase = iota + 1
	RequirePlurals
	RequireDeterminers // kinds and kinds of kinds
	RequireAncestry    // parameters precede results
	RequireParameters  // results precede normal fields
	RequireResults     // properties of kinds and records
	RequireFields      // initial values of fields
	RequireDefaults
	RequireNouns
	RequireNames
	RequireRelatives
	RequireRules
	RequireAll
)

type Assertions interface {
	AssertDomainStart(name string, requires []string) error
	AssertDomainEnd() error

	AssertAlias(short string, names ...string) error
	AssertAncestor(kind, ancestor string) error
	AssertAspectTraits(aspect string, traits []string) error

	AssertField(kind, name, class string, aff affine.Affinity, init assign.Assignment) error
	AssertParam(kind, name, class string, aff affine.Affinity, init assign.Assignment) error
	AssertResult(kind, name, class string, aff affine.Affinity, init assign.Assignment) error

	AssertGrammar(name string, directive *grammar.Directive) error
	AssertNounKind(noun, kind string) error
	// tbd: should path be a DottedPath?
	AssertNounValue(nounName, fieldName string, path []string, val literal.LiteralValue) error
	AssertOpposite(a, b string) error
	AssertPlural(singluar, plural string) error
	AssertRelation(rel, a, b string, amany, bmany bool) error
	AssertRelative(rel, noun, otherNoun string) error

	// tbd: can this be replaced by a rule or something?
	AssertCheck(name string, do []rt.Execute, expect literal.LiteralValue) error

	// fix: target should become part of the guard.
	// and/or rule should be wrapped up more like "grammar.Directive"
	AssertRule(name string, target string, guard rt.BoolEval, flags EventTiming, do []rt.Execute) error

	// any application defined key-value pair
	// ( the last element is the value, the prefix is the key )
	AssertDefinition(path ...string) error
}
