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

type NameCalled interface {
	GetName(traits, kinds []Matched) grok.Name
}
