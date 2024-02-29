package jess

import (
	"errors"
	"fmt"

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

	// find the name of the kind which best matches the passed span ( if any )
	// return the number of words that matched.
	FindKind(match.Span, *kindsOf.Kinds) (string, int)

	// find the name of the trait which best matches the passed span ( if any )
	// return the number of words that matched.
	FindTrait(match.Span) (string, int)

	// find the name of the field which best matches the passed span ( if any )
	// return the number of words that matched.
	FindField(match.Span) (string, int)

	// find the macro which best matches the passed span ( if any )
	// return the number of words that matched.
	FindMacro(match.Span) (Macro, int)

	/// find the name of the kind which best matches the passed span ( if any )
	// return the number of words that matched.
	FindNoun(match.Span) (string, int)
}

// Matched - generic interface so implementations can track backchannel data.
type Matched interface {
	String() string
}

// implemented by phrases so that they can create story fragments based on
// the english language text they have parsed.
type Generator interface {
	Generate(Registrar) error
}

// used internally for matching some kinds of phrases.
type Interpreter interface {
	Match(Query, *InputState) bool
}

// matches an english like sentence against jess's parse trees.
// returns an object which can create nouns, define kinds, set properties, and so on.
func Match(q Query, ws match.Span) (ret Generator, err error) {
	var m MatchingPhrases
	input := InputState(ws)
	if m, ok := m.Match(q, &input); !ok {
		err = errors.New("failed to match phrase")
	} else if cnt := len(input); cnt != 0 {
		err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
	} else {
		ret = m
	}
	return
}
