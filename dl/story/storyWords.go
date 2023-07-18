package story

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *MakePlural) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (op *MakePlural) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
		if plural := lang.Normalize(op.Plural); len(plural) < 0 {
			err = errutil.New("no plural specified")
		} else if singular := lang.Normalize(op.Singular); len(singular) < 0 {
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
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
		if a := lang.Normalize(op.Word); len(a) < 0 {
			err = errutil.New("no word for opposite specified")
		} else if b := lang.Normalize(op.Opposite); len(b) < 0 {
			err = errutil.New("no opposite for word specified")
		} else {
			err = w.Pin().AddOpposite(a, b)
		}
		return
	})
}
