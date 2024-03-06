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
		// understand "..." as .....
		&op.Understand,
		// names are "a kind of"/"kinds of" [traits] kind:any.
		&op.KindsOf,
		// kind:objects are "usually" traits.
		&op.KindsAreTraits,
		// kinds:records|objects "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
		&op.KindsHaveProperties,
		// kinds:objects ("can be"|"are either") new_trait [or new_trait...]
		&op.KindsAreEither,
		// kind:aspects are names
		&op.AspectsAreTraits,
		&op.VerbNamesAreNames,
		&op.NamesVerbNames,
		&op.NamesAreLikeVerbs,
		&op.PropertyNounValue,
		&op.NounPropertyValue,
		&op.MapLocations,
		&op.MapDirections,
		&op.MapConnections,
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
