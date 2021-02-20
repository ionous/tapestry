package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type At struct {
	List  rt.Assignment
	Index rt.NumberEval
}

// future: lists of lists? probably through lists of records containing lists.
func (*At) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_at",
		Group: "list",
		Spec:  "list {list:assignment} at {index:number_eval}",
		Desc:  "Value of List: Get a value from a list. The first element is is index 1.",
	}
}

func (op *At) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.NumList)
}

func (op *At) GetText(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.TextList)
}

func (op *At) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.RecordList)
}

func (op *At) getAt(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
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
