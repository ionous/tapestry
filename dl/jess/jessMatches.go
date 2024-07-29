package jess

import (
	"fmt"
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// the top level matching interface implemented by members of MatchedPhrase.
// "Line" here means "sentence of a paragraph"
type LineMatcher interface {
	MatchLine(Query, InputState) (InputState, bool)
	typeinfo.Instance // for logging
	typeinfo.Markup
}

// different phases (z) can match different phrases (ws)
// should a match occur, return true; and set 'out' to the matched phrase.
func matchSentence(z weaver.Phase, q Query, line InputState, out *bestMatch) (okay bool) {
	var op MatchedPhrase
	switch z {
	case weaver.LanguagePhase:
		// "understand" {quoted text} as .....
		okay = matchLine(q, line, &op.Understand, out)

	case weaver.AncestryPhase:
		// FIX -- KindsAreKind needs TRAITS to *match*
		// to do the idea of match once, generate often;
		// this would have to delay parsing the trailing phrase.
		// probably part of the phase has to be scheduled; while the basic naming does not.
		// ---
		// names {are} "a kind of"/"kinds of" [traits] kind.
		okay = matchLine(q, line, &op.KindsAreKind, out) ||
			// The colors are black and blue.
			matchLine(q, line, schedule(line, &op.AspectsAreTraits), out)

	case weaver.PropertyPhase:
		// fix? combine these to speed matching?
		// kinds {are} "usually"
		okay = matchLine(q, line, schedule(line, &op.KindsAreTraits), out) ||
			// kinds(of records|objects, out) "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
			matchLine(q, line, schedule(line, &op.KindsHaveProperties), out) ||
			// kinds(of objects, out) ("can be"|"are either", out) new_trait [or new_trait...]
			matchLine(q, line, schedule(line, &op.KindsAreEither), out)

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
			matchLine(q, line, &op.MapConnections, out) ||
				// direction "of/from" place {is} place.
				matchLine(q, line, &op.MapDirections, out) ||
				// place {is} direction "of/from" places.
				matchLine(q, line, &op.MapLocations, out) ||

				// field "of" noun {are} value
				matchLine(q, line, &op.PropertyNounValue, out) ||
				// noun "has" field value
				matchLine(q, line, &op.NounPropertyValue, out) ||

				// verb nouns {are} nouns
				matchLine(q, line, &op.VerbNamesAreNames, out) ||
				// nouns {are} verbing nouns
				matchLine(q, line, &op.NamesVerbNames, out) ||
				// nouns {are} adjectives [verb nouns]
				matchLine(q, line, &op.NamesAreLikeVerbs, out)

	case weaver.VerbPhase:
		okay = matchLine(q, line, &op.TimedRule, out)
	}
	return
}

type bestMatch struct {
	match    matchGenerator
	numWords int
	pronouns pronounSource
}

// each LineMatchers is either a Generator, a schedulee.
// ( to manually schedule pieces of its data )
type matchGenerator interface {
	LineMatcher
	Generator
}

// a variation that can weaves on a specific phase.
// it gets wrapped with a call to "schedule()"
// which turns it into a standard Generator.
type schedulee interface {
	LineMatcher
	Phase() weaver.Phase
	Weave(weaver.Weaves, rt.Runtime) error
}

// phases that can weave immediately, without needing to schedule more phases
// can use this to define a Generate method
func schedule(line InputState, s schedulee) genericSchedule {
	return genericSchedule{line, s}
}

type genericSchedule struct {
	line InputState
	schedulee
}

func (op genericSchedule) Generate(ctx Context) (err error) {
	return ctx.Schedule(op.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		if e := op.Weave(w, run); e != nil {
			err = schedulingError{op, e, op.Phase()}
		}
		return
	})
}

type schedulingError struct {
	op    genericSchedule
	err   error
	phase weaver.Phase
}

func (e schedulingError) Unwrap() error {
	return e.err
}

func (e schedulingError) Error() string {
	src := e.op.line.Source().ErrorString()
	return fmt.Sprintf("%s during %s at %s for %s", e.err, e.phase, src, e.op.TypeInfo().TypeName())
}

// match the input against the passed parse tree.
// passes out an object which can create nouns, define kinds, set properties, and so on.
// returns the number of words *not* matched
func matchLine(q Query, line InputState, op matchGenerator, out *bestMatch) (okay bool) {
	// "understand" {quoted text} as .....
	if next, ok := op.MatchLine(q, line); ok {
		// was the phrase only partially matched?
		if remaining := next.Len(); remaining > 0 {
			wordsMatched := line.Len() - remaining
			if wordsMatched > out.numWords {
				out.numWords = wordsMatched
			}
		} else {
			if useLogging(q) {
				m := Matched(line.words)
				log.Printf("matched %s %q\n", op.TypeInfo().TypeName(), m.DebugString())
			}
			(*out) = bestMatch{match: op, pronouns: next.pronouns}
			okay = true
		}
	}
	return
}
