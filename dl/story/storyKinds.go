package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineKinds) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. "cats are a kind of animal"
func (op *DefineKinds) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if kinds, e := safe.GetTextList(w, op.Kinds); e != nil {
			err = e
		} else if ancestor, e := safe.GetText(w, op.Ancestor); e != nil {
			err = e
		} else {
			pen := w.Pin()
			ancestor := inflect.Normalize(ancestor.String())
			for _, kind := range kinds.Strings() {
				kind := match.StripArticle(inflect.Normalize(kind))
				if e := pen.AddKind(kind, ancestor); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineFields) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		run := cat.Runtime()
		if kind, e := safe.GetText(run, op.Kind); e != nil {
			err = e
		} else {
			k := inflect.Normalize(kind.String())
			err = w.Pin().AddFields(k, reduceFields(run, op.Fields))
		}
		return
	})
}
