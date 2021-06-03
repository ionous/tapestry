package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"

	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// GetBool returns the first matching bool evaluation.
func (op *CallSend) GetBool(run rt.Runtime) (g.Value, error) {
	return op.send(run, affine.Bool)
}

func (op *CallSend) send(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if path, e := safe.GetTextList(run, op.Path); e != nil {
		err = e
	} else {
		name, up := op.Event.String(), path.Strings()
		if v, e := run.Send(name, up, op.Arguments.Pack()); e != nil {
			err = cmdErrorCtx(op, name, e)
		} else {
			ret = v
		}
	}
	return
}
