package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *ListAt) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.NumList)
}

func (op *ListAt) GetText(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.TextList)
}

func (op *ListAt) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.RecordList)
}

func (op *ListAt) getAt(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if vs, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if e := safe.Check(vs, aff); e != nil {
		err = cmdError(op, e)
	} else if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = cmdError(op, e)
	} else if i, e := safe.Range(idx.Int()-1, 0, vs.Len()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs.Index(i)
	}
	return
}
