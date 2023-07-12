package story

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// helper to simplify setting the values of nouns
func assertNounValue(a assert.Assertions, val literal.LiteralValue, noun string, path ...string) error {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return a.AssertNounValue(noun, field, parts, val)
}
