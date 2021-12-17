package express

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{n} }
func T(s string) *literal.TextValue { return &literal.TextValue{s} }

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
