package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type Weaver struct {
	d     *Domain
	at    string
	phase assert.Phase
	rt.Runtime
}
