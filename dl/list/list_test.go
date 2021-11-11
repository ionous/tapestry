package list_test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func F(n float64) rt.NumberEval     { return &core.NumValue{n} }
func I(n int) rt.NumberEval         { return &core.NumValue{float64(n)} }
func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func T(s string) *core.TextValue    { return &core.TextValue{W(s)} }
func V(i string) *core.GetVar       { return &core.GetVar{N(i)} }
func W(v string) string             { return v }

func FromTs(vs []string) (ret rt.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&core.TextValue{vs[0]}}
	} else {
		ret = &core.FromTexts{&core.Texts{vs}}
	}
	return
}