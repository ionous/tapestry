package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Write a story to a story file.
func Encode(src *StoryFile) (interface{}, error) {
	return cout.Encode(src, CompactEncoder)
}

// customized writer of compact data
var CompactEncoder = core.CompactEncoder

// story break is an empty do nothing statement, used as a paragraph marker.
func (op *StoryBreak) Schedule(k *weave.Catalog) error { return nil }

// Execute - called by the macro runtime during weave.
func (op *StoryBreak) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}
