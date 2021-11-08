package story

import (
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"github.com/ionous/errutil"
)

func (op *TestRule) ImportPhrase(k *Importer) (err error) {
	if n, e := NewTestName(k, op.TestName); e != nil {
		err = e
	} else if hook, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else if prog, e := k.NewProg("execute", hook); e != nil {
		err = e
	} else {
		k.NewTestProgram(n, prog)
	}
	return
}

func (op *TestScene) ImportPhrase(k *Importer) (err error) {
	// handled separately so we can have separate begin/end frames
	return
}

func (op *TestStatement) ImportPhrase(k *Importer) (err error) {
	if t := op.Test; t == nil {
		err = ImportError(op, op.At, errutil.Fmt("%w Test", MissingSlot))
	} else if n, e := NewTestName(k, op.TestName); e != nil {
		err = e
	} else {
		err = t.ImportTest(k, n)
	}
	return
}

type Testing interface {
	ImportTest(k *Importer, testName eph.Named) (err error)
}

func (op *TestOutput) ImportTest(k *Importer, testName eph.Named) (err error) {
	// note: we use the raw lines here, we don't expect the text output to be a template.
	k.NewTestExpectation(testName, "execute", op.Lines.Str)
	return
}
