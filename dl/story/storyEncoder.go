package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Write a story to a story file.
func Encode(src *StoryFile) (interface{}, error) {
	return cout.Encode(src, CompactEncoder)
}

// customized writer of compact data
var CompactEncoder = core.CompactEncoder

// story break is an empty do nothing statement, used as a paragraph marker.
func (op *StoryBreak) Weave(cat *weave.Catalog) error { return nil }
