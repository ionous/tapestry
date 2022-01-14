package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
)

func (op *Comment) ImportStub(k *Importer) (ret interface{}, err error) {
	if !k.InProgram() {
		ret = op
	} else {
		ret = &debug.DebugLog{
			Value: &core.FromText{
				Val: T(op.Lines.String()),
			},
			LogLevel: debug.LoggingLevel{
				Str: debug.LoggingLevel_Note,
			},
		}
	}
	return
}
