package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

func (op *HasDominion) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	return run.GetField(meta.Domain, op.Name)
}
