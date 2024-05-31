package story

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/printer"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/weave"
)

// ensure that a valid counter exists
func (op *CycleText) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	ret = &printer.CallCycle{Parts: op.Parts, Name: cat.NewCounter("seq")}
	return
}

// ensure that a valid counter exists
func (op *ShuffleText) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	ret = &printer.CallShuffle{Parts: op.Parts, Name: cat.NewCounter("seq")}
	return
}

// ensure that a valid counter exists
func (op *StoppingText) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	ret = &printer.CallTerminal{Parts: op.Parts, Name: cat.NewCounter("seq")}
	return
}

// ensure that a valid counter exists
func (op *CountOf) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	ret = &call.CallTrigger{Num: op.Num, Trigger: op.Trigger, Name: cat.NewCounter("seq")}
	return
}
