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
func matchSentence(z weaver.Phase, q Query, next []match.TokenValue, out *bestMatch) (okay bool) {
	var op MatchingPhrases
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
		// fix? combine these to speed matching?
		// kinds {are} "usually"
		okay = matchPhrase(q, next, schedule(&op.KindsAreTraits), out) ||
			// kinds(of records|objects, out) "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
			matchPhrase(q, next, schedule(&op.KindsHaveProperties), out) ||
			// kinds(of objects, out) ("can be"|"are either", out) new_trait [or new_trait...]
			matchPhrase(q, next, schedule(&op.KindsAreEither), out)

	case weaver.NounPhase:
		// note: the direction phrases have to be first
		// so that the verbs phrases don't incorporate the names of directions
		// ( ie. "west of summit" shouldn't be considered a noun. )
		// similarly, we have to check for fields first so field names don't become noun names.
		// ( ie. "the subject of carrying is actors" shouldn't be confused with
		// ... "the bottle is a container." )
		//
		// fix: the confusion with fields seems dangerous...
		// the specific issue is allowing kinds to be used as text.
		// inform doesn't allow the text of kinds ( its a tapestry hack for verbs )
		// -- perhaps a special phrase for unquoted kinds:
		//
		// but also, in inform:
		//   The friend is an animal.
		//   Containers have a container called friend.
		// errors with "'friend' is a nothing valued property, and it is too late to change now."
		// it's globally distinguishing b/t fields and nouns of the same whole name.
		// ( ex. maybe a "name table" that says what type things are )
		// ( and then GetClosestNoun is GetClosestName and it returns a type )
		// note: you can say "magic friend" is an animal; but you cant later refer to it as "friend"
		// the property wins.
		okay = //
			// "through" door {is} place.
			matchPhrase(q, next, &op.MapConnections, out) ||
				// direction "of/from" place {is} place.
				matchPhrase(q, next, &op.MapDirections, out) ||
				// place {is} direction "of/from" places.
				matchPhrase(q, next, &op.MapLocations, out) ||

				// field "of" noun {are} value
				matchPhrase(q, next, &op.PropertyNounValue, out) ||
				// noun "has" field value
				matchPhrase(q, next, &op.NounPropertyValue, out) ||

				// verb nouns {are} nouns
				matchPhrase(q, next, &op.VerbNamesAreNames, out) ||
				// nouns {are} verbing nouns
				matchPhrase(q, next, &op.NamesVerbNames, out) ||
				// nouns {are} adjectives [verb nouns]
				matchPhrase(q, next, &op.NamesAreLikeVerbs, out)

	case weaver.VerbPhase:
		okay = matchPhrase(q, next, &op.TimedRule, out)
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
		// was the whole phrase matched?
		if remaining := next.Len(); remaining > 0 {
			if cnt := input.Len() - remaining; cnt > out.numWords {
				out.numWords = cnt
			}
		} else {
			if useLogging(q) {
				m := Matched(input)
				log.Printf("matched %s %q\n", op.TypeInfo().TypeName(), m.DebugString())
			}
			(*out) = bestMatch{match: op}
			okay = true
		}
	}
	return
}
