package story

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/weave"
)

// helper to simplify setting the values of nouns
func assertNounValue(cat *weave.Catalog, val literal.LiteralValue, noun string, path ...string) error {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return cat.AddNounValue(noun, field, parts, val)
}
