package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *MakePlural) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (op *MakePlural) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		if plural, e := safe.GetText(w, op.Plural); e != nil {
			err = e
		} else if singular, e := safe.GetText(w, op.Singular); e != nil {
			err = e
		} else if plural := inflect.Normalize(plural.String()); len(plural) < 0 {
			err = errutil.New("no plural specified")
		} else if singular := inflect.Normalize(singular.String()); len(singular) < 0 {
			err = errutil.New("no singular specified")
		} else {
			err = w.Pin().AddPlural(plural, singular)
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *MakeOpposite) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *MakeOpposite) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		if word, e := safe.GetText(w, op.Word); e != nil {
			err = e
		} else if opposite, e := safe.GetText(w, op.Opposite); e != nil {
			err = e
		} else if a := inflect.Normalize(word.String()); len(a) < 0 {
			err = errutil.New("no word for opposite specified")
		} else if b := inflect.Normalize(opposite.String()); len(b) < 0 {
			err = errutil.New("no opposite for word specified")
		} else {
			err = w.Pin().AddOpposite(a, b)
		}
		return
	})
}
