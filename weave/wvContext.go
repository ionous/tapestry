package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type Weaver struct {
	Domain *Domain
	At     string
	Phase  assert.Phase
	rt.Runtime
}

func (w *Weaver) Grok(p string) (grok.Results, error) {
	cat := w.Domain.cat
	return cat.gdb.Grok(w.Domain.name, p)
}
