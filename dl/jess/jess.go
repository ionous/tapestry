package jess

import (
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// Matching requires identifying kinds, traits, and macros.
// tdb: the returned objects are all "string" --
// but it might be more useful if they were "resources"
// and if generator functions required those resources as targets.
type Query interface {
	// provides for customization of individual queries.
	// implementations should return 0;
	// individual parse trees can then wrap the query with their own context specific information..
	GetContext() int

	// find the name of the kind which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindKind(match.Span, *kindsOf.Kinds) (string, int)

	// find the name of the trait which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindTrait(match.Span) (string, int)

	// find the name of the field which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindField(match.Span) (string, int)

	// find the name of the noun which best matches the passed span.
	// return the number of words that matched ( if any. )
	// the kind, if specified, will ensure the noun is of that kind;
	// [ so that the caller doesn't have to validate materialized kind paths ]
	FindNoun(name match.Span, kind string) (string, int)

	// provides some limited access to previous values assigned to nouns.
	// ( assumes the caller already knows the affinity of the field )
	GetNounValue(noun, field string) ([]byte, error)
}

// Matched - generic interface so implementations can track backchannel data.
type Matched interface {
	String() string
}

// implemented by phrases so that they can create story fragments based on
// the english language text they have parsed.
type Generator interface {
	Generate(*Context) error
}

// used internally for matching some kinds of phrases.
type Interpreter interface {
	Match(Query, *InputState) bool
}
