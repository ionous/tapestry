package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/debug"
)

func (op *Comment) ImportStub(k *Importer) (ret interface{}, err error) {
	if !inProg(k) {
		ret = op
	} else {
		ret = &debug.DebugLog{
			&core.FromText{
				&core.TextValue{op.Lines.String()},
			},
			debug.LoggingLevel{
				debug.LoggingLevel_Note,
			},
		}
	}
	return
}

// a hopefully temporary hack
func inProg(k *Importer) (ret bool) {
	for _, k := range k.decoder.Path {
		if k == "core.Activity" {
			ret = true
			break
		}
	}
	return
}
