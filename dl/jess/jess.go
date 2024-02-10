package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches
type Matched = grok.Matched
type Span = grok.Span
type Macro = grok.Macro

type Interpreter interface {
	Match(Query, *InputState) bool
}
