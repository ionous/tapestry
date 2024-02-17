package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

// matches
type Matched = grok.Matched
type Span = grok.Span
type Macro = grok.Macro

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
	GetResults() (grok.Results, error)
}

func Match(g grok.Grokker, ws grok.Span) (Interpreter, error) {
	return MatchLog(g, ws, LogWarning)
}

func MatchLog(g grok.Grokker, ws grok.Span, level LogLevel) (ret Interpreter, err error) {
	query := MakeQueryLog(g, level)
	input := InputState(ws)
	if m, ok := match(query, &input); !ok {
		err = errors.New("failed to match phrase")
	} else if cnt := len(input); cnt != 0 {
		err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
	} else {
		ret = m
	}
	return
}

func match(q Query, input *InputState) (ret Matches, okay bool) {
	var m MatchingPhrases
	return m.Match(q, input)
}

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Matches, okay bool) {
	// fix? could change to reflect ( or expand type info ) to walk generically
	var best InputState
	for _, m := range []Matches{
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

func makeResult(macro grok.Macro, reverse bool, a, b []grok.Name) grok.Results {
	if reverse {
		a, b = b, a
	}
	return grok.Results{
		Primary:   a,
		Secondary: b,
		Macro:     macro,
	}
}
