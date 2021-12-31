package express

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/rt"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Bool: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Num: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Num: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Text: s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(s string) core.VariableName { return core.VariableName{Str: s} }
func W(v string) string            { return v }

func O(n string, exact bool) (ret rt.TextEval) {
	if !exact {
		ret = &render.RenderRef{N(n), render.RenderFlags{Str: render.RenderFlags_RenderAsAny}}
	} else {
		ret = T(n)
	}
	return ret
}
