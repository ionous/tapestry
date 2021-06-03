package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

// note: there's only one kind of hook now: the activity
// though we do have to change from story activity to core activity for ... reasons.
func (op *ProgramHook) ImportProgram(k *Importer) (ret *core.Activity, err error) {
	if opt, ok := op.Opt.(*core.Activity); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		ret = &core.Activity{opt.Exe}
	}
	return
}
