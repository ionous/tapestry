package list_test

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(v string) core.VariableName { return core.VariableName{Str: v} }
func V(i string) *core.GetVar      { return &core.GetVar{Name: N(i)} }
func W(v string) string            { return v }

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{Val: &literal.TextValue{Value: vs[0]}}
	} else {
		ret = &core.FromTexts{Vals: &literal.TextValues{Values: vs}}
	}
	return
}
