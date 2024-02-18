package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/match"
)

type Span = match.Span

// Matched - generic interface so implementations can track backchannel data.
// ex. a row id for kinds.
type Matched interface {
	String() string
	NumWords() int // should match strings.Fields(Str)
}

type Query interface {
	// if the passed words starts with a kind,
	// return the number of words that matched.
	FindKind(Span) (Matched, int)

	// if the passed words starts with a trait,
	// return the number of words that matched.
	FindTrait(Span) (Matched, int)

	// if the passed words starts with a macro,
	// return information about that match.
	FindMacro(Span) (Macro, int)
}

type Applicant interface {
	Interpreter
	Apply(Registrar) error
}

type Interpreter interface {
	Match(Query, *InputState) bool
}

// the root of a sentence matching tree
// can produce full results when matched,
type Matches interface {
	Interpreter
	GetResults() (localResults, error)
}

func Match(q Query, ws match.Span) (ret Applicant, err error) {
	input := InputState(ws)
	if m, ok := matchit(q, &input); !ok {
		err = errors.New("failed to match phrase")
	} else if cnt := len(input); cnt != 0 {
		err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
	} else {
		ret = m
	}
	return
}

func matchit(q Query, input *InputState) (ret Applicant, okay bool) {
	var m MatchingPhrases
	return m.Match(q, input)
}

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Applicant, okay bool) {
	// fix? could change to reflect ( or expand type info ) to walk generically
	var best InputState
	for _, m := range []Applicant{
		&op.KindsAreTraits,
		&op.KindsOf,
		&op.VerbLinks,
		&op.LinksVerb,
		&op.LinksAdjectives,
	} {
		if next := *input; //
		m.Match(q, &next) /* && len(next) == 0 */ {
			if !okay || len(next) < len(best) {
				best = next
				ret, okay = m, true
				if len(best) == 0 {
					break
				}
			}
		}
	}
	if okay {
		*input = best
	}
	return
}

func makeResult(macro Macro, reverse bool, a, b []resultName) localResults {
	if reverse {
		a, b = b, a
	}
	return localResults{
		Primary:   a,
		Secondary: b,
		Macro:     macro,
	}
}
