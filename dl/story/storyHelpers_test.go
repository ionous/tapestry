package story_test

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
)

var (
	B = literal.B
	F = literal.F
	I = literal.I
	T = literal.T
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}
