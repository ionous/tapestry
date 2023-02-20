package list_test

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
)

var (
	B           = core.B
	F           = core.F
	I           = core.I
	N           = core.N
	P           = core.P
	T           = core.T
	W           = core.W
	GetVariable = core.GetVariable
)

func FromTs(vs []string) (ret assign.Assignment) {
	if len(vs) == 1 {
		ret = &assign.FromText{Value: core.T(vs[0])}
	} else {
		ret = &assign.FromTextList{Value: core.Ts(vs)}
	}
	return
}
