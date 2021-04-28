package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type Len struct {
	List rt.Assignment
}

func (*Len) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_len",
		Lede:  "len",
		Group: "list",
		Spec:  "length of {list:assignment}",
		Desc:  "Length of List: Determines the number of values in a list.",
	}
}

func (op *Len) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if !affine.IsList(v.Affinity()) {
		err = cmdError(op, errutil.New("not a list"))
	} else {
		ret = g.IntOf(v.Len())
	}
	return
}
