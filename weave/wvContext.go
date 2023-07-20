package weave

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Weaver struct {
	Catalog *Catalog
	Domain  *Domain
	At      string
	Phase   Phase
	rt.Runtime
}

func (w *Weaver) GrokSpan(p grok.Span) (grok.Results, error) {
	return w.Catalog.gdb.GrokSpan(w.Domain.name, p)
}

func (w *Weaver) MatchArticle(ws []string) (ret int, err error) {
	return w.Catalog.gdb.MatchArticle(ws)
}

func (w *Weaver) AddNoun(name, kind string) (err error) {
	long, short, kind := name, lang.Normalize(name), lang.Normalize(kind)
	_, err = w.Domain.AddNoun(long, short, kind)
	return
}

func (w *Weaver) GetClosestNoun(name string) (ret string, err error) {
	if bare, e := grok.StripArticle(name); e != nil {
		err = e
	} else if n, e := w.Domain.GetClosestNoun(lang.Normalize(bare)); e != nil {
		err = e
	} else {
		ret = n.name
	}
	return
}

func (w *Weaver) Pin() *mdl.Pen {
	return w.Catalog.Modeler.Pin(w.Domain.name, w.At)
}
