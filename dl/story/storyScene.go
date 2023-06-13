package story

import (
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

func (op *DefineScene) Schedule(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireAll, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetOptionalText(w, op.Scene, ""); e != nil {
			err = e
		} else if dependsOn, e := safe.GetOptionalTexts(w, op.DependsOn, nil); e != nil {
			err = e
		} else {
			if e := cat.BeginDomain(name.String(), dependsOn.Strings()); e != nil {
				err = e
			} else if e := ScheduleStatements(cat, op.With); e != nil {
				err = e
			} else {
				err = cat.EndDomain()
			}
		}
		return
	})
}
