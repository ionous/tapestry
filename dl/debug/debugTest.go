package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
)

func (op *Test) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = op // provisionally
	var req []string
	if n := op.DependsOn.String(); len(n) > 0 {
		req = []string{n}
	}
	// everything between this and "EndDomain" in the Post/PostImport will be in this test domain.
	k.WriteEphemera(&eph.EphBeginDomain{Name: op.TestName.String(), Requires: req})
	return
}

func (op *Test) PostImport(k *imp.Importer) (err error) {
	k.WriteEphemera(&eph.EphEndDomain{})
	k.WriteEphemera(&eph.EphChecks{Name: op.TestName.String(), Exe: op.Do})
	return
}
