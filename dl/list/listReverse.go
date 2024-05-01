package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// A normal reduce would return a value, instead we accumulate into a variable
func (op *ListReverse) Execute(run rt.Runtime) (err error) {
	if e := op.reverse(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListReverse) reverse(run rt.Runtime) (err error) {
	if at, e := assign.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else {
		cnt := vs.Len()
		for i := cnt/2 - 1; i >= 0; i-- {
			j := cnt - 1 - i
			eli, elj := vs.Index(i), vs.Index(j)
			// while technically SetIndex returns error,
			// because we are setting to ourself, it should be fine.
			vs.SetIndex(i, elj)
			vs.SetIndex(j, eli)
		}
	}
	return
}
