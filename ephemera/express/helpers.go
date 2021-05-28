package express

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

func T(s string) rt.TextEval {
	return &core.TextValue{value.Text(s)}
}
func N(n float64) rt.NumberEval {
	return &core.NumValue{n}
}
func B(b bool) rt.BoolEval {
	return &core.BoolValue{b}
}
func O(n string, exact bool) (ret rt.TextEval) {
	if !exact {
		ret = &render.RenderRef{n, render.RenderFlags{&render.RenderAsAny{}}}
	} else {
		ret = T(n)
	}
	return ret
}
