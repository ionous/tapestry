package weave

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
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

func (w *Weaver) GetKindByName(name string) (ret *g.Kind, err error) {
	// fix: we poll the db and once its there ask for more info
	// if we ask for info first, qna returns "Unknown kind" --
	// and even if it returned "missing kind" the error would get cached.
	// option: return "Missing" from qnaKind and implement a runtime config that doest cache?
	if _, e := w.Pin().GetKind(name); e != nil {
		err = e
	} else if k, e := w.Runtime.GetKindByName(name); e != nil {
		err = e
	} else {
		ret = k
	}
	return
}
