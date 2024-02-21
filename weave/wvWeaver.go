package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

type Weaver struct {
	Catalog *Catalog
	Domain  string
	At      string
	Phase   Phase
	rt.Runtime
}

func (w *Weaver) MatchSpan(p match.Span) (jess.Generator, error) {
	return w.Catalog.MatchSpan(w.Domain, p)
}

func (w *Weaver) Pin() *mdl.Pen {
	return w.Catalog.Modeler.Pin(w.Domain, w.At)
}

func (w *Weaver) AddInitialValue(pen *mdl.Pen, noun, field string, value rt.Assignment) (err error) {
	// if we are adding an initial value for a different domain
	// then that gets changed into "set value" triggered on "begin domain"
	var u mdl.DomainValueError
	if e := pen.AddInitialValue(noun, field, value); !errors.As(e, &u) {
		err = e // nil or unexpected error.
	} else {
		d := w.Catalog.domains[w.Domain]
		d.initialValues = d.initialValues.add(u.Noun, u.Field, u.Value)
	}
	return
}

func (w *Weaver) GetClosestNoun(name string) (ret string, err error) {
	if bare, e := jess.StripArticle(name); e != nil {
		err = e
	} else if n := inflect.Normalize(bare); len(n) == 0 {
		err = errutil.New("empty name")
	} else if n, e := w.Pin().GetClosestNoun(n); e != nil {
		err = e
	} else {
		ret = n
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
