package story

import (
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func (op *DefineTest) Weave(cat *weave.Catalog) (err error) {
	// some of the tests ( ex. TestImportSequence ) dont establish a top level scene
	// so they cant schedule. maybe we could unroll it so "Push" saves the current schedule
	// ( rather than push creating the first schedule )
	// return cat.Schedule(weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
	if domain, reqs, e := GetSceneReqs(cat.GetRuntime(), op.TestName, op.RequiredSceneNames); e != nil {
		err = e
	} else if pen, e := cat.PushScene(cat.EnsureScene(domain), mdl.MakeSource(op)); e != nil {
		err = e
	} else {
		defer cat.PopScene()
		if e := pen.AddDependency(reqs...); e != nil {
			err = e
		} else if e := Weave(cat, op.Statements); e != nil {
			err = e
		} else if len(op.Exe) > 0 {
			err = pen.AddCheck(domain, nil, op.Exe)
		}
	}
	return
	// })
}
