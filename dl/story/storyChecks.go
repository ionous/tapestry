package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"github.com/ionous/errutil"
)

func (op *TestRule) ImportPhrase(k *Importer) (err error) {
	if n, e := makeTestName(k, op.TestName); e != nil {
		err = e
	} else if exe, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else {
		k.Write(&eph.EphChecks{Name: n, Exe: exe})
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
	} else if n, e := makeTestName(k, op.TestName); e != nil {
		err = e
	} else {
		err = t.ImportTest(k, n)
	}
	return
}

type Testing interface {
	ImportTest(k *Importer, testName string) (err error)
}

func (op *TestOutput) ImportTest(k *Importer, testName string) (err error) {
	// note: we use the raw lines here, we don't expect the text output to be a template.
	k.Write(&eph.EphChecks{Name: testName, Expect: &literal.TextValue{op.Lines.Str}})
	return
}

func makeTestName(k *Importer, n TestName) (ret string, err error) {
	if n.Str == TestName_CurrentTest {
		ret = k.Env().Recent.Test
	} else {
		ret = n.Str
	}
	return
}
