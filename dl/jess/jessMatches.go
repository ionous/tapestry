package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Generator, okay bool) {
	type matcher interface {
		Generator
		Interpreter
		typeinfo.Instance
	}
	var best InputState
	var bestMatch matcher
	// fix? could change to reflect ( or expand type info ) to walk generically
	for _, m := range []matcher{
		//
		&op.Understand, // understand "..." as .....
		// kinds
		&op.KindsOf,             // names {are} "a kind of"/"kinds of" [traits] kind.
		&op.KindsAreTraits,      // objects {are} "usually" traits.
		&op.KindsHaveProperties, // records|objects "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
		&op.KindsAreEither,      // objects ("can be"|"are either") new_trait [or new_trait...]
		&op.AspectsAreTraits,    // aspects {are} names
		// before phrases that might think the directions are part of the names.
		&op.MapConnections, // "through" door {is} place.
		&op.MapDirections,  // direction "of/from" place {is} place.
		&op.MapLocations,   // place {is} direction "of/from" places.
		// noun phrases
		&op.VerbNamesAreNames, // verb nouns {are} nouns
		&op.NamesVerbNames,    // nouns {are} verbing nouns
		&op.NamesAreLikeVerbs, // nouns {are} adjectives [verb nouns]
		&op.PropertyNounValue, // property "of" noun is value
		&op.NounPropertyValue, // noun "has" property value

	} {
		if next := *input; //
		m.Match(q, &next) /* && len(next) == 0 */ {
			if !okay || next.Len() < best.Len() {
				best = next
				bestMatch, okay = m, true
				if best.Len() == 0 {
					break
				}
			}
		}
	}
	if okay {
		if useLogging(q) {
			log.Printf("matched %s %q\n", bestMatch.TypeInfo().TypeName(), match.Span(input.Words()).String())
		}
		ret, *input = bestMatch, best

	}
	return
}
