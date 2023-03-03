package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// helps patterns handle values when the affinity isn't known in advance
type RenderEval interface {
	RenderEval(run rt.Runtime, hint affine.Affinity) (g.Value, error)
}
