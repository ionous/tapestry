package story

import (
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
)

func (op *DefineScene) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetOptionalText(w, op.Scene, ""); e != nil {
			err = e
		} else if dependsOn, e := safe.GetOptionalTexts(w, op.DependsOn, nil); e != nil {
			err = e
		} else {
			if e := cat.DomainStart(name.String(), dependsOn.Strings()); e != nil {
				err = e
			} else if e := WeaveStatements(cat, op.With); e != nil {
				err = e
			} else {
				err = cat.DomainEnd()
			}
		}
		return
	})
}
