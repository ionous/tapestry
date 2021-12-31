package pattern_test

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Bool: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Num: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Num: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Text: s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(v string) core.VariableName { return core.VariableName{Str: v} }
func V(i string) *core.GetVar      { return &core.GetVar{N(i)} }
func W(v string) string            { return v }
