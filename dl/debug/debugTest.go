package debug

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *Test) PreImport(k *imp.Importer) (ret interface{}, err error) {
	ret = op // provisionally
	var req []string
	if n := op.DependsOn.String(); len(n) > 0 {
		req = []string{n}
	}
	// everything between this and "EndDomain" in the Post/PostImport will be in this test domain.
	err = k.BeginDomain(op.TestName.String(), req)
	return
}

// Execute - called by the macro runtime during weave.
func (op *Test) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *Test) PostImport(k *imp.Importer) (err error) {
	if e := k.EndDomain(); e != nil {
		err = e
	} else if e := k.AssertCheck(op.TestName.String(), op.Do, nil); e != nil {
		err = e
	}
	return
}
