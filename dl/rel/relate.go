package rel

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *Relate) Execute(run rt.Runtime) (err error) {
	if e := op.setRelation(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Relate) setRelation(run rt.Runtime) (err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = e
	} else if b, e := safe.ObjectText(run, op.ToObject); e != nil {
		err = e
	} else {
		err = run.RelateTo(a.String(), b.String(), op.Via)
	}
	return
}
