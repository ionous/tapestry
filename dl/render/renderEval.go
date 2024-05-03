package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// helps patterns handle values when the affinity isn't known in advance
type RenderEval interface {
	RenderEval(run rt.Runtime, hint affine.Affinity) (rt.Value, error)
}
