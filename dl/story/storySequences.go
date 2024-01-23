package story

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/weave"
)

// ensure that a valid counter exists
func (op *CycleText) PreImport(cat *weave.Catalog) (ret typeinfo.Inspector, err error) {
	ret = &core.CallCycle{Parts: op.Parts, Name: newCounter(cat, "seq")}
	return
}

// ensure that a valid counter exists
func (op *ShuffleText) PreImport(cat *weave.Catalog) (ret typeinfo.Inspector, err error) {
	ret = &core.CallShuffle{Parts: op.Parts, Name: newCounter(cat, "seq")}
	return
}

// ensure that a valid counter exists
func (op *StoppingText) PreImport(cat *weave.Catalog) (ret typeinfo.Inspector, err error) {
	ret = &core.CallTerminal{Parts: op.Parts, Name: newCounter(cat, "seq")}
	return
}

// ensure that a valid counter exists
func (op *CountOf) PreImport(cat *weave.Catalog) (ret typeinfo.Inspector, err error) {
	ret = &core.CallTrigger{Num: op.Num, Trigger: op.Trigger, Name: newCounter(cat, "seq")}
	return
}

func newCounter(cat *weave.Catalog, name string) (ret string) {
	next := cat.Env.Inc(name, 1)
	return name + "-" + strconv.Itoa(next)
}
