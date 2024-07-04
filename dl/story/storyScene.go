package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *DefineScene) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if domain, reqs, e := op.GetSceneReqs(run); e != nil {
			err = e
		} else if pen, e := cat.PushScene(cat.EnsureScene(domain), mdl.MakeSource(op)); e != nil {
			err = e
		} else {
			defer cat.PopScene()
			if e := pen.AddDependency(reqs...); e != nil {
				err = e
			} else if e := Weave(cat, op.Statements); e != nil {
				err = e // add all the statements that are a part of this domain.
			}
		}
		return
	})
}

func (op *DefineScene) GetSceneReqs(run rt.Runtime) (string, []string, error) {
	return GetSceneReqs(run, op.SceneName, op.RequiredSceneNames)
}

// errs if there's no scene name
func GetSceneReqs(run rt.Runtime, scene rt.TextEval, reqs rt.TextListEval) (retScene string, retReqs []string, err error) {
	if name, e := safe.GetText(run, scene); e != nil {
		err = e
	} else if dependsOn, e := safe.GetOptionalTexts(run, reqs, nil); e != nil {
		err = e
	} else {
		retScene = inflect.Normalize(name.String())
		retReqs = dependsOn.Strings()
		for i, req := range retReqs {
			retReqs[i] = inflect.Normalize(req)
		}
	}
	return
}
