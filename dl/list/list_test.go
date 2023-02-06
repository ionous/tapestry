package list_test

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

var B = core.B
var F = core.F
var I = core.I
var N = core.N
var P = core.P
var T = core.T
var V = core.V
var W = core.W

var SetVar = core.SetVar

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{Val: &literal.TextValue{Value: vs[0]}}
	} else {
		ret = &core.FromTexts{Vals: &literal.TextValues{Values: vs}}
	}
	return
}
