package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// different phases (z) can match different phrases (ws)
// should a match occur, return true; and set 'out' to the matched phrase.
func matchSentence(q Query, z Phase, ws match.Span, out *bestMatch) (okay bool) {
	var op MatchingPhrases
	next := MakeInput(ws)
	switch z {
	case mdl.LanguagePhase:
		// "understand" {quoted text} as .....
		okay = matchPhrase(q, next, &op.Understand, out)

	case mdl.AncestryPhase:
		// FIX -- these need TRAITS to *match*
		// to do the idea of match once, generate often;
		// this would have to delay parsing the trailing phrase.
		// probably part of the phase has to be scheduled; while the basic naming does not.
		// ---
		// names {are} "a kind of"/"kinds of" [traits] kind.
		okay = matchPhrase(q, next, &op.KindsOf, out)

	case mdl.PropertyPhase:
		// fix: it'd be nice to re-factor these so they dont each have to match kinds
		// ex. kinds(of aspects, out) {are} names
		// The colors are red, blue, and greasy green
		okay = matchPhrase(q, next, &op.AspectsAreTraits, out) ||
			// kinds(of objects, out) {are} "usually" traits.
			matchPhrase(q, next, &op.KindsAreTraits, out) ||
			// kinds(of records|objects, out) "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
			matchPhrase(q, next, &op.KindsHaveProperties, out) ||
			// kinds(of objects, out) ("can be"|"are either", out) new_trait [or new_trait...]
			matchPhrase(q, next, &op.KindsAreEither, out)

	case mdl.MappingPhase:
		// "through" door {is} place.
		okay = matchPhrase(q, next, &op.MapConnections, out) ||
			// direction "of/from" place {is} place.
			matchPhrase(q, next, &op.MapDirections, out) ||
			// place {is} direction "of/from" places.
			matchPhrase(q, next, &op.MapLocations, out)

	case mdl.NounPhase:
		// verb nouns {are} nouns
		okay = matchPhrase(q, next, &op.VerbNamesAreNames, out) ||
			// nouns {are} verbing nouns
			matchPhrase(q, next, &op.NamesVerbNames, out) ||
			// nouns {are} adjectives [verb nouns]
			matchPhrase(q, next, &op.NamesAreLikeVerbs, out)

	case mdl.FallbackPhase:
		// property "of" noun {are} value
		okay = matchPhrase(q, next, &op.PropertyNounValue, out) ||
			// noun "has" property value
			matchPhrase(q, next, &op.NounPropertyValue, out)
	}
	return
}

// fix: actually declare all the members with this interface as a slot?
type matcher interface {
	Generator
	Interpreter
	typeinfo.Instance
	Phase() mdl.Phase
}

type bestMatch struct {
	match    Generator
	numWords int
}

// match the input against the passed parse tree.
// passes out an object which can create nouns, define kinds, set properties, and so on.
// returns the number of words *not* matched
func matchPhrase(q Query, input InputState, op matcher, out *bestMatch) (okay bool) {
	// "understand" {quoted text} as .....
	if next := input; op.Match(q, &next) {
		if remaining := next.Len(); remaining > 0 {
			if cnt := input.Len() - remaining; cnt > out.numWords {
				out.numWords = cnt
			}
		} else {
			if useLogging(q) {
				log.Printf("matched %s %q\n", op.TypeInfo().TypeName(), match.Span(input.Words()).String())
			}
			(*out) = bestMatch{match: op}
			okay = true
		}
	}
	return
}
