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

	// find the name of the kind which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindKind(match.Span, *kindsOf.Kinds) (string, int)

	// find the name of the trait which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindTrait(match.Span) (string, int)

	// find the name of the field which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindField(match.Span) (string, int)

	// find the macro which best matches the passed span.
	// return the number of words that matched ( if any. )
	FindMacro(match.Span) (Macro, int)

	// find the name of the noun which best matches the passed span.
	// return the number of words that matched ( if any. )
	// the kind, if specified, will ensure the noun is of that kind;
	// [ so that the caller doesn't have to validate materialized kind paths ]
	FindNoun(name match.Span, kind string) (string, int)
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

// match one or more sentences and use the registrar
// to create  nouns, define kinds, set properties, and so on.
func Generate(q Query, w Registrar, paragraph string) (err error) {
	if spans, e := match.MakeSpans(paragraph); e != nil {
		err = fmt.Errorf("%w reading %s", e, paragraph)
	} else {
		for _, span := range spans {
			if m, e := Match(q, span); e != nil {
				err = fmt.Errorf("%w matching %s", e, span)
				break
			} else if e := m.Generate(w); e != nil {
				err = fmt.Errorf("%w generating %s", e, paragraph)
				break
			}
		}
	}
	return
}

// matches an english like sentence against jess's parse trees.
// returns an object which can create nouns, define kinds, set properties, and so on.
func Match(q Query, ws match.Span) (ret Generator, err error) {
	var m MatchingPhrases
	input := MakeInput(ws)
	if m, ok := m.Match(q, &input); !ok {
		err = errors.New("failed to match phrase")
	} else if cnt := input.Len(); cnt != 0 {
		err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
	} else {
		ret = m
	}
	return
}
