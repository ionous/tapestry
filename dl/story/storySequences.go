package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/imp"
)

// fix? are these primports stubs really needed anymore?
// perhaps we can use reserved fields of the counter computation?
// the advantage of the stubs is that it keeps knowledge of the importer out of the "core" which is otherwise just runtime
// an alternative is to put a post-processing hook inside "importStory" like ResponseType, etc. do.
// the only reason against that is the big switch statement all in one place
// maybe, there's a way to stub the signature/s in certain contexts to say --
// all of these things you think are type X are actually type y during import [ even if its a type alias ]

// ensure that a valid counter exists
func (op *CycleText) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = &core.CallCycle{Parts: op.Parts, Name: k.NewCounter("seq", op.Markup)}
	return
}

// ensure that a valid counter exists
func (op *ShuffleText) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = &core.CallShuffle{Parts: op.Parts, Name: k.NewCounter("seq", op.Markup)}
	return
}

// ensure that a valid counter exists
func (op *StoppingText) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = &core.CallTerminal{Parts: op.Parts, Name: k.NewCounter("seq", op.Markup)}
	return
}

// ensure that a valid counter exists
func (op *CountOf) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = &core.CallTrigger{Num: op.Num, Trigger: op.Trigger, Name: k.NewCounter("seq", op.Markup)}
	return
}
