package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
)

func (op *TestRule) ImportPhrase(k *Importer) (err error) {
	if n, e := makeTestName(k, op.TestName); e != nil {
		err = e
	} else {
		k.WriteEphemera(&eph.EphChecks{Name: n, Exe: op.Does})
	}
	return
}

func (op *TestScene) ImportPhrase(k *Importer) (err error) {
	// handled separately so we can have separate begin/end frames
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
