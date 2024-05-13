package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *DefineScene) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if a, b, e := op.GetSceneReqs(run); e != nil {
			err = e
		} else if e := cat.DomainStart(a, b); e != nil {
			err = e
		} else if e := Weave(cat, op.Statements); e != nil {
			err = e
		} else {
			err = cat.DomainEnd()
		}
		return
	})
}

// errs if there's no scene name
func (op *DefineScene) GetSceneReqs(run rt.Runtime) (retScene string, retReqs []string, err error) {
	if name, e := safe.GetText(run, op.SceneName); e != nil {
		err = e
	} else if dependsOn, e := safe.GetOptionalTexts(run, op.RequiredSceneNames, nil); e != nil {
		err = e
	} else {
		retScene = name.String()
		retReqs = dependsOn.Strings()
	}
	return
}
