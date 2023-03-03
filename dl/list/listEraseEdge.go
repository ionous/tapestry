package list

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *EraseEdge) Execute(run rt.Runtime) (err error) {
	if _, e := eraseEdge(run, op.Target, op.AtEdge); e != nil {
		err = CmdError(op, e)
	}
	return
}

func eraseEdge(run rt.Runtime, target assign.Address, atFront rt.BoolEval) (ret g.Value, err error) {
	if root, e := assign.GetRootValue(run, target); e != nil {
		err = e
	} else if els, e := root.GetList(run); e != nil {
		err = e
	} else if cnt := els.Len(); cnt > 0 {
		var at int
		if atFront, e := safe.GetOptionalBool(run, atFront, false); e != nil {
			err = e
		} else {
			if !atFront.Bool() {
				at = cnt - 1
			}
			if v, e := els.Splice(at, at+1, nil); e != nil {
				err = e
			} else {
				root.SetDirty(run)
				ret = v
			}
		}
	}
	return
}
