package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

// note: there's only one kind of hook now: the activity
func (op *ProgramHook) ImportProgram(k *Importer) (ret core.Activity, err error) {
	if opt, ok := op.Value.(*core.Activity); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		ret = *opt
	}
	return
}
