package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *DefineScene) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if name, e := safe.GetOptionalText(run, op.Scene, ""); e != nil {
			err = e
		} else if dependsOn, e := safe.GetOptionalTexts(run, op.RequireScenes, nil); e != nil {
			err = e
		} else {
			if e := cat.DomainStart(name.String(), dependsOn.Strings()); e != nil {
				err = e
			} else if e := WeaveStatements(cat, op.Statements); e != nil {
				err = e
			} else {
				err = cat.DomainEnd()
			}
		}
		return
	})
}
