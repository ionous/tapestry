package story

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Execute - called by the macro runtime during weave.
func (op *MakePlural) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *MakePlural) PostImport(k *imp.Importer) error {
	return k.AssertPlural(op.Singular, op.Plural)
}

// Execute - called by the macro runtime during weave.
func (op *MakeOpposite) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *MakeOpposite) PostImport(k *imp.Importer) error {
	return k.AssertOpposite(op.Opposite, op.Word)
}
