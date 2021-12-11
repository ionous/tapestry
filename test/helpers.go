package test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/value"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{n} }
func T(s string) *literal.TextValue { return &literal.TextValue{s} }

func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func V(i string) *core.GetVar       { return &core.GetVar{N(i)} }
func W(v string) string             { return v }
