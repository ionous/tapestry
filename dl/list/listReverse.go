package list

import (
	"git.sr.ht/~ionous/iffy/rt"
)

// A normal reduce would return a value, instead we accumulate into a variable
func (op *ListReverse) Execute(run rt.Runtime) (err error) {
	if e := op.reverse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ListReverse) reverse(run rt.Runtime) (err error) {
	if els, e := GetListSource(run, op.List); e != nil {
		err = e
	} else {
		cnt := els.Len()
		for i := cnt/2 - 1; i >= 0; i-- {
			j := cnt - 1 - i
			eli, elj := els.Index(i), els.Index(j)
			els.SetIndex(i, elj)
			els.SetIndex(j, eli)
		}
	}
	return
}
