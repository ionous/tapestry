package jess

import (
	"fmt"
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// different phases (z) can match different phrases (ws)
// should a match occur, return true; and set 'out' to the matched phrase.
func matchSentence(q Query, z weaver.Phase, ws match.Span, out *bestMatch) (okay bool) {
	var op MatchingPhrases
	next := MakeInput(ws)
	switch z {
	case weaver.LanguagePhase:
		// "understand" {quoted text} as .....
		okay = matchPhrase(q, next, &op.Understand, out)

	case weaver.AncestryPhase:
		// FIX -- KindsOf needs TRAITS to *match*
		// to do the idea of match once, generate often;
		// this would have to delay parsing the trailing phrase.
		// probably part of the phase has to be scheduled; while the basic naming does not.
		// ---
		// names {are} "a kind of"/"kinds of" [traits] kind.
		okay = matchPhrase(q, next, &op.KindsOf, out) ||
			// The colors are black and blue.
			matchPhrase(q, next, schedule(&op.AspectsAreTraits), out)

	case weaver.PropertyPhase:
		// kinds {are} "usually"
		okay = matchPhrase(q, next, schedule(&op.KindsAreTraits), out) ||
			// kinds(of records|objects, out) "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
			matchPhrase(q, next, schedule(&op.KindsHaveProperties), out) ||
			// kinds(of objects, out) ("can be"|"are either", out) new_trait [or new_trait...]
			matchPhrase(q, next, schedule(&op.KindsAreEither), out)

	case weaver.NounPhase:
		// note: the direction phrases have to be first here
		// so that the verbs phrases don't incorporate the names of directions
		// into the names of nouns ( ex. so that "west of summit" isn't a place name. )
		//
		// "through" door {is} place.
		okay = matchPhrase(q, next, &op.MapConnections, out) ||
			// direction "of/from" place {is} place.
			matchPhrase(q, next, &op.MapDirections, out) ||
			// place {is} direction "of/from" places.
			matchPhrase(q, next, &op.MapLocations, out) ||
			// verb nouns {are} nouns
			matchPhrase(q, next, &op.VerbNamesAreNames, out) ||
			// nouns {are} verbing nouns
			matchPhrase(q, next, &op.NamesVerbNames, out) ||
			// nouns {are} adjectives [verb nouns]
			matchPhrase(q, next, &op.NamesAreLikeVerbs, out) ||
			// property "of" noun {are} value
			matchPhrase(q, next, &op.PropertyNounValue, out) ||
			// noun "has" property value
			matchPhrase(q, next, &op.NounPropertyValue, out)
	}
	return
}

type bestMatch struct {
	match    Generator
	numWords int
}

type matcher interface {
	Interpreter
	Generator
	typeinfo.Instance // for logging
}

type schedulee interface {
	Interpreter
	typeinfo.Instance // for logging
	Phase() weaver.Phase
	Weave(weaver.Weaves, rt.Runtime) error
}

// phases that can weave immediately, without needing to schedule more phases
// can use this to define a Generate method
func schedule(s schedulee) genericSchedule {
	return genericSchedule{s}
}

type genericSchedule struct{ schedulee }

func (op genericSchedule) Generate(ctx Context) (err error) {
	return ctx.Schedule(op.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		if e := op.Weave(w, run); e != nil {
			err = fmt.Errorf("%w during %q", e, op.TypeInfo().TypeName())
		}
		return
	})
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
