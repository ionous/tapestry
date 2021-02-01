package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

type programImporter interface {
	ImportProgram(k *Importer) (ret programHook, err error)
}

func (op *ProgramHook) ImportProgram(k *Importer) (ret programHook, err error) {
	if opt, ok := op.Opt.(programImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		ret, err = opt.ImportProgram(k)
	}
	return
}

func (op *Activity) ImportProgram(k *Importer) (ret programHook, err error) {
	ret = &executeSlot{&core.Activity{op.Exe}}
	return
}
