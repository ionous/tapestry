package express

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func F(n float64) rt.NumberEval     { return &core.NumValue{n} }
func N(s string) value.VariableName { return value.VariableName{s} }
func P(s string) value.PatternName  { return value.PatternName{s} }
func T(s string) rt.TextEval        { return &core.TextValue{value.Text{s}} }
func W(v string) value.Text         { return value.Text{Str: v} }

func O(n string, exact bool) (ret rt.TextEval) {
	if !exact {
		ret = &render.RenderRef{N(n), render.RenderFlags{&render.RenderAsAny{}}}
	} else {
		ret = T(n)
	}
	return ret
}
