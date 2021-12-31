package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// a utility, primarily used for testing, which allows values to be passed directly to commands which take parameters
type FromValue struct{ g.Value }

func (op *FromValue) Affinity() affine.Affinity { return "" }

func (op *FromValue) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	ret = op.Value
	return
}
