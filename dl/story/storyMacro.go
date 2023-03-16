package story

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Execute - called by the macro runtime during weave.
func (op *DefineMacro) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *DefineMacro) PostImport(k *imp.Importer) (err error) {
	return nil
}
