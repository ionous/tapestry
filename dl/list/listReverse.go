package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

// A normal reduce would return a value, instead we accumulate into a variable
type Reverse struct {
	List ListSource `if:"selector"`
}

func (*Reverse) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_reverse",
		Group:  "list",
		Desc:   `Reverse List: Reverse a list.`,
		Fluent: &composer.Fluid{Role: composer.Command},
	}
}

func (op *Reverse) Execute(run rt.Runtime) (err error) {
	if e := op.reverse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Reverse) reverse(run rt.Runtime) (err error) {
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
