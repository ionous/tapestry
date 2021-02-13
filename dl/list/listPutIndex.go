package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atIndex: numEval.
 */
type PutIndex struct {
	From    core.Assignment `if:"selector"`
	Into    ListTarget      `if:"selector"`
	AtIndex rt.NumberEval
}

func (*PutIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put at index: replace one value in a list with another",
	}
}

func (op *PutIndex) Execute(run rt.Runtime) (err error) {
	return errutil.New("not implemented")
}
