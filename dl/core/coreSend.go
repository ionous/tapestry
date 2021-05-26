package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"

	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// GetBool returns the first matching bool evaluation.
func (op *Send) GetBool(run rt.Runtime) (g.Value, error) {
	return op.send(run, affine.Bool)
}

func (op *Send) send(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if path, e := safe.GetTextList(run, op.Path); e != nil {
		err = e
	} else {
		var args []rt.Arg
		for _, a := range op.Arguments.Args {
			args = append(args, rt.Arg{a.Name.Value(), a.From})
		}
		name, up := op.Event.Value(), path.Strings()
		if v, e := run.Send(name, up, args); e != nil {
			err = cmdErrorCtx(op, name, e)
		} else {
			ret = v
		}
	}
	return
}
