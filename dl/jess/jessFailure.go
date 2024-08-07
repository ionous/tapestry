package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/match"
)

type FailedMatch struct {
	reason string
	input  InputState
}

func (m FailedMatch) Error() string {
	// fix: add source
	return fmt.Sprintf("%s matching %s", m.reason, match.DebugStringify(m.input.words))
}
