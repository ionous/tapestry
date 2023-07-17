package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Weaver struct {
	Catalog *Catalog
	Domain  *Domain
	At      string
	Phase   assert.Phase
	rt.Runtime
}

func (w *Weaver) Grok(p string) (grok.Results, error) {
	return w.Catalog.gdb.Grok(w.Domain.name, p)
}

func (w *Weaver) MatchArticle(ws []string) (ret int, err error) {
	return w.Catalog.gdb.MatchArticle(ws)
}

func (w *Weaver) Pin() *mdl.Pen {
	return w.Catalog.Modeler.Pin(w.Domain.name, w.At)
}
