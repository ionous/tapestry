package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

const activityDepth = "activityDepth"

// Schedule - comment does nothing when imported.
func (*Comment) Weave(*weave.Catalog) (_ error) {
	return
}

// Execute - called by the macro runtime during weave.
// since Schedule doesnt do anything, neither do we.
func (*Comment) Execute(rt.Runtime) (_ error) {
	return
}

// PreImport turns a comment statement into a debug log.
func (op *Comment) PreImport(cat *weave.Catalog) (ret typeinfo.Inspector, err error) {
	if cat.Env.Inc(activityDepth, 0) == 0 {
		ret = op
	} else {
		ret = &debug.DebugLog{
			Value:    &assign.FromText{Value: core.T(op.Lines)},
			LogLevel: debug.C_LoggingLevel_Note,
		}
	}
	return
}
