package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Relate) Execute(run rt.Runtime) (err error) {
	if e := op.setRelation(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *Relate) setRelation(run rt.Runtime) (err error) {
	if a, e := safe.ObjectText(run, op.NounName); e != nil {
		err = e
	} else if b, e := safe.ObjectText(run, op.OtherNounName); e != nil {
		err = e
	} else if rel, e := safe.GetText(run, op.RelationName); e != nil {
		err = e
	} else {
		err = run.RelateTo(a.String(), b.String(), rel.String())
	}
	return
}
