package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Execute - called by the macro runtime during weave.
func (op *MakePlural) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

func (op *MakePlural) Schedule(cat *weave.Catalog) error {
	return cat.AssertPlural(op.Singular, op.Plural)
}

// Execute - called by the macro runtime during weave.
func (op *MakeOpposite) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

func (op *MakeOpposite) Schedule(cat *weave.Catalog) error {
	return cat.AssertOpposite(op.Opposite, op.Word)
}
