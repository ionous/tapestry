package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type Gather struct {
	Var   core.Variable `if:"selector"`
	From  ListSource    `if:"selector"`
	Using pattern.PatternName
}

func (*Gather) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_gather",
		Fluent: &composer.Fluid{Name: "gather", Role: composer.Command},
		Group:  "list",
		Desc: `Gather list: Transform the values from a list.
		The named pattern gets called once for each value in the list.
		It get called with two parameters: 'in' as each value from the list, 
		and 'out' as the var passed to the gather.`,
	}
}

func (op *Gather) Execute(rt.Runtime) (ret g.Value, err error) {
	err = errutil.New("not implemented")
	return
}
