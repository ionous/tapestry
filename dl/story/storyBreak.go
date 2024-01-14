package story

import (
	"git.sr.ht/~ionous/tapestry/weave"
)

// story break is an empty do nothing statement, used as a paragraph marker.
func (op *StoryBreak) Weave(cat *weave.Catalog) error {
	return nil
}
