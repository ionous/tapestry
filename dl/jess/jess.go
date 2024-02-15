package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches
type Matched = grok.Matched
type Span = grok.Span
type Macro = grok.Macro

type Interpreter interface {
	Match(Query, *InputState) bool
}

// the root of a sentence matching tree
// can produce full results when matched,
type Matches interface {
	Interpreter
	GetResults(Query) (grok.Results, error)
}

func Match(q Query, input *InputState) (ret Matches, okay bool) {
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

func makeResult(macro grok.Macro, a, b []grok.Name) grok.Results {
	if macro.Reversed {
		a, b = b, a
	}
	return grok.Results{
		Primary:   a,
		Secondary: b,
		Macro:     macro,
	}
}
