package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/debug"
	"git.sr.ht/~ionous/iffy/dl/literal"
)

func (op *Comment) ImportStub(k *Importer) (ret interface{}, err error) {
	if !k.InProgram() {
		ret = op
	} else {
		ret = &debug.DebugLog{
			&core.FromText{
				&literal.TextValue{op.Lines.String()},
			},
			debug.LoggingLevel{
				debug.LoggingLevel_Note,
			},
		}
	}
	return
}
