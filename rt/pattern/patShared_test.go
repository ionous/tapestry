package pattern_test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func I(n int) rt.NumberEval         { return &literal.NumValue{float64(n)} }
func F(n float64) rt.NumberEval     { return &literal.NumValue{n} }
func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func T(s string) *literal.TextValue { return &literal.TextValue{W(s)} }
func V(i string) *core.GetVar       { return &core.GetVar{N(i)} }
func W(v string) string             { return v }
