package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
)

// A normal reduce would return a value, instead we accumulate into a variable
func (op *ListReverse) Execute(run rt.Runtime) (err error) {
	if e := op.reverse(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListReverse) reverse(run rt.Runtime) (err error) {
	if root, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if els, e := root.GetList(run); e != nil {
		err = e
	} else {
		cnt := els.Len()
		for i := cnt/2 - 1; i >= 0; i-- {
			j := cnt - 1 - i
			eli, elj := els.Index(i), els.Index(j)
			// while technically SetIndex returns error,
			// because we are setting to ourself, it should be fine.
			els.SetIndex(i, elj)
			els.SetIndex(j, eli)
		}
	}
	return
}
