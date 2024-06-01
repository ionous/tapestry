package story

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/weave"
)

// ensure that a valid counter exists
func (op *CountOf) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	ret = &call.CallTrigger{Num: op.Num, Trigger: op.Trigger, Name: cat.NewCounter("seq")}
	return
}
