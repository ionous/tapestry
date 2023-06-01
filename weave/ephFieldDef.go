package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
)

type fieldDef struct {
	name      string
	affinity  affine.Affinity
	class     string
	at        string
	initially assign.Assignment
}
