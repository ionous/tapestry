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

func eraseEdge(run rt.Runtime, tgt assign.Address, atFront rt.BoolEval) (ret g.Value, err error) {
	if at, e := assign.GetReference(run, tgt); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.CheckList(vs); e != nil {
		err = e
	} else if cnt := vs.Len(); cnt > 0 {
		var idx int
		if atFront, e := safe.GetOptionalBool(run, atFront, false); e != nil {
			err = e
		} else {
			if !atFront.Bool() {
				idx = cnt - 1
			}
			if v, e := vs.Splice(idx, idx+1, nil); e != nil {
				err = e
			} else {
				ret = v
			}
		}
	}
	return
}
