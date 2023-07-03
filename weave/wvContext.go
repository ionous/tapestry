package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type Weaver struct {
	Domain *Domain
	At     string
	Phase  assert.Phase
	rt.Runtime
}
