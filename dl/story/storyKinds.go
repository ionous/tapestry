package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// ex. "cats are a kind of animal"
func (op *DefineKinds) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kinds, e := safe.GetTextList(run, op.Kinds); e != nil {
			err = e
		} else if ancestor, e := safe.GetText(run, op.Ancestor); e != nil {
			err = e
		} else {
			ancestor := inflect.Normalize(ancestor.String())
			for _, kind := range kinds.Strings() {
				kind := match.StripArticle(inflect.Normalize(kind))
				if e := w.AddKind(kind, ancestor); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		return
	})
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.PropertyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kind, e := safe.GetText(run, op.Kind); e != nil {
			err = e
		} else {
			k := inflect.Normalize(kind.String())
			err = w.AddKindFields(k, reduceFields(run, op.Fields))
		}
		return
	})
}
