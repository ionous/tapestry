package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
)

// ensure that a valid counter exists
func (op *CycleText) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = &core.CallCycle{Parts: op.Parts, At: k.newCounter("seq", op.At)}
	return
}

// ensure that a valid counter exists
func (op *ShuffleText) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = &core.CallShuffle{Parts: op.Parts, At: k.newCounter("seq", op.At)}
	return
}

// ensure that a valid counter exists
func (op *StoppingText) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = &core.CallTerminal{Parts: op.Parts, At: k.newCounter("seq", op.At)}
	return
}

// ensure that a valid counter exists
func (op *CountOf) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = &core.CallTrigger{Num: op.Num, Trigger: op.Trigger, At: k.newCounter("seq", op.At)}
	return
}

// generate a unique name for the counter --
// for stability's sake, preferring an existing id in the source to an autogenerated id.
func (k *Importer) newCounter(name string, at reader.Position) (ret reader.Position) {
	if at.IsValid() {
		ret = at
	} else {
		ret.Offset = k.autoCounter.Next(name)
	}
	return
}
