package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

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
