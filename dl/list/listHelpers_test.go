package list_test

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

var (
	B = literal.B
	F = literal.F
	I = literal.I
	T = literal.T
	W = core.W
)

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &assign.FromText{Value: literal.T(vs[0])}
	} else {
		ret = &assign.FromTextList{Value: literal.Ts(vs...)}
	}
	return
}
