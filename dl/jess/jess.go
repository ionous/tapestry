package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches
type Matched = grok.Matched

type Interpreter interface {
	Match(Query, InputState) int
}
