package list

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *EraseEdge) Execute(run rt.Runtime) (err error) {
	if _, e := eraseEdge(run, op.From, op.AtEdge); e != nil {
		err = cmdError(op, e)
	}
	return
}

func eraseEdge(run rt.Runtime,
	from ListSource,
	atFront rt.BoolEval,
) (ret g.Value, err error) {
	if vs, e := GetListSource(run, from); e != nil {
		err = e
	} else {
		if cnt := vs.Len(); cnt > 0 {
			var at int
			if atFront, e := safe.GetOptionalBool(run, atFront, false); e != nil {
				err = e
			} else if !atFront.Bool() {
				at = cnt - 1
			}
			ret, err = vs.Splice(at, at+1, nil)
		}
	}
	return
}
