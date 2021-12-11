package list_test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{n} }
func T(s string) *literal.TextValue { return &literal.TextValue{s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(v string) core.VariableName { return core.VariableName{Str: v} }
func V(i string) *core.GetVar      { return &core.GetVar{N(i)} }
func W(v string) string            { return v }

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&literal.TextValue{vs[0]}}
	} else {
		ret = &core.FromTexts{&literal.TextValues{vs}}
	}
	return
}
