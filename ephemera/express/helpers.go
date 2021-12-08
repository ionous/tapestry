package express

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func F(n float64) rt.NumberEval     { return &literal.NumValue{n} }
func N(s string) value.VariableName { return value.VariableName{Str: s} }
func P(s string) value.PatternName  { return value.PatternName{Str: s} }
func T(s string) rt.TextEval        { return &literal.TextValue{W(s)} }
func W(v string) string             { return v }

func O(n string, exact bool) (ret rt.TextEval) {
	if !exact {
		ret = &render.RenderRef{N(n), render.RenderFlags{Str: render.RenderFlags_RenderAsAny}}
	} else {
		ret = T(n)
	}
	return ret
}
