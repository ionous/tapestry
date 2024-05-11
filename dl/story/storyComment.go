package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

const activityDepth = "activityDepth"

// Schedule - comment does nothing when imported.
func (*Note) Weave(*weave.Catalog) (_ error) {
	return
}

// Execute - panics. PreImport should turn it into a DebugLog.
func (*Note) Execute(rt.Runtime) (_ error) {
	panic("story comment should have been replaced during weave")
}

// PreImport turns a comment statement into a debug log.
func (op *Note) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	if cat.Env.Inc(activityDepth, 0) == 0 {
		ret = op
	} else {
		ret = &debug.DebugLog{
			Value: &assign.FromTextList{
				Value: &literal.TextValues{Values: op.Lines},
			},
			LogLevel: debug.C_LoggingLevel_Note,
		}
	}
	return
}
