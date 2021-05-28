package pattern_test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(i bool) rt.BoolEval   { return &core.BoolValue{i} }
func I(i int) rt.NumberEval  { return &core.NumValue{float64(i)} }
func T(i string) rt.TextEval { return &core.TextValue{i} }
func V(i string) *core.Var   { return &core.Var{Name: i} }
