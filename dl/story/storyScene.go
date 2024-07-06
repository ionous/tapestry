package story

import (
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
)

func (op *DefineScene) Weave(cat *weave.Catalog) (err error) {
	// some of the tests ( ex. TestImportSequence ) dont establish a top level scene
	// so they cant schedule, and have to use the annoying cat.GetRuntime()
	// maybe there could be a "current schedule" list
	// and push stacks that ( rather than push creating the first schedule )
	// return cat.ScheduleCmd(op, weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
	if domain, reqs, e := op.GetSceneReqs(cat.GetRuntime()); e != nil {
		err = e
	} else {
		domain := cat.EnsureScene(domain)
		pos := compact.MakeSource(op.GetMarkup(false))
		if pen, e := cat.SceneBegin(domain, pos); e != nil {
			err = e
		} else {
			defer cat.SceneEnd()
			if e := pen.AddDependency(reqs...); e != nil {
				err = e
			} else if e := Weave(cat, op.Statements); e != nil {
				err = e // add all the statements that are a part of this domain.
			} else {
				domain.AddStartup(op.Exe)
			}
		}
	}
	return
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
