package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefinePlural) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (op *DefinePlural) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if plural, e := safe.GetText(run, op.Plural); e != nil {
			err = e
		} else if singular, e := safe.GetText(run, op.Singular); e != nil {
			err = e
		} else if plural := inflect.Normalize(plural.String()); len(plural) < 0 {
			err = errutil.New("no plural specified")
		} else if singular := inflect.Normalize(singular.String()); len(singular) < 0 {
			err = errutil.New("no singular specified")
		} else {
			err = w.AddPlural(plural, singular)
		}
		return
	})
}
