package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// ex. "cats are a kind of animal"
func (op *DefineKind) Weave(cat *weave.Catalog) error {
	return cat.ScheduleCmd(op, weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kind, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else if ancestor, e := safe.GetText(run, op.AncestorKindName); e != nil {
			err = e
		} else {
			kind := inflect.Normalize(kind.String())
			ancestor := inflect.Normalize(ancestor.String())
			err = w.AddKind(kind, ancestor)
		}
		return
	})
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.PropertyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kind, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else {
			k := inflect.Normalize(kind.String())
			err = w.AddKindFields(k, reduceFields(run, op.FieldNames))
		}
		return
	})
}
