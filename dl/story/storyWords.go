package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Execute - called by the macro runtime during weave.
func (op *MakePlural) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *MakePlural) Weave(cat *weave.Catalog) error {
	return cat.AssertPlural(op.Singular, op.Plural)
}

// Execute - called by the macro runtime during weave.
func (op *MakeOpposite) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *MakeOpposite) Weave(cat *weave.Catalog) error {
	return cat.AssertOpposite(op.Opposite, op.Word)
}
