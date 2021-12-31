package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

func (op *HasDominion) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return run.GetField(meta.Domain, op.Name)
}
