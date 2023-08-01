package weave

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Weaver struct {
	Catalog *Catalog
	Domain  string
	At      string
	Phase   Phase
	rt.Runtime
}

func (w *Weaver) GrokSpan(p grok.Span) (grok.Results, error) {
	return w.Catalog.gdb.GrokSpan(w.Domain, p)
}

func (w *Weaver) MatchArticle(ws []string) (ret int, err error) {
	return w.Catalog.gdb.MatchArticle(ws)
}

func (w *Weaver) Pin() *mdl.Pen {
	return w.Catalog.Modeler.Pin(w.Domain, w.At)
}

func (w *Weaver) GetClosestNoun(name string) (ret string, err error) {
	if bare, e := grok.StripArticle(name); e != nil {
		err = e
	} else {
		n := lang.Normalize(bare)
		if n, e := w.Pin().GetClosestNoun(n); e != nil {
			err = e
		} else {
			ret = n
		}
	}
	return
}
