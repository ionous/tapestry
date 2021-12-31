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
			&core.FromText{
				T(op.Lines.String()),
			},
			debug.LoggingLevel{
				debug.LoggingLevel_Note,
			},
		}
	}
	return
}
