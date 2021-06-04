package core

import (
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

func (op *HasDominion) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return run.GetField(object.Domain, op.Name.String())
}
