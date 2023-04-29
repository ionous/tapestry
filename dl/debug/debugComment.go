package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Schedule - comment does nothing when imported.
func (*Comment) Schedule(*weave.Catalog) (_ error) {
	return
}

// PreImport turns a comment statement into a debug log.
func (op *Comment) PreImport(k *weave.Catalog) (ret interface{}, err error) {
	if !k.Env().InProgram() {
		ret = op
	} else {
		ret = &DebugLog{
			Value: &assign.FromText{Value: core.T(op.Lines.String())},
			LogLevel: LoggingLevel{
				Str: LoggingLevel_Note,
			},
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
// since Schedule doesnt do anything, neither do we.
func (*Comment) Execute(rt.Runtime) (_ error) {
	return
}
