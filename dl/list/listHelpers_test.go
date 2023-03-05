package list_test

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
)

var (
	B = assign.B
	F = assign.F
	I = assign.I
	N = assign.N
	P = assign.P
	T = assign.T
	W = assign.W
)

func FromTs(vs []string) (ret assign.Assignment) {
	if len(vs) == 1 {
		ret = &assign.FromText{Value: assign.T(vs[0])}
	} else {
		ret = &assign.FromTextList{Value: assign.Ts(vs...)}
	}
	return
}
