package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

func (op *DefineTest) Weave(cat *weave.Catalog) (err error) {
	if name := inflect.Normalize(op.TestName); len(name) == 0 {
		errutil.New("test has empty name")
	} else if dependsOn, e := safe.GetOptionalTexts(cat.GetRuntime(), op.SceneNames, nil); e != nil {
		err = e
	} else {
		if e := cat.DomainStart(name, dependsOn.Strings()); e != nil {
			err = e
		} else {
			if e := WeaveStatements(cat, op.Statements); e != nil {
				err = e
			} else if len(op.Exe) > 0 {
				err = cat.Schedule(weaver.NextPhase, func(w weaver.Weaves, run rt.Runtime) error {
					return w.AddCheck(name, nil, op.Exe)
				})
			}
			if err == nil {
				err = cat.DomainEnd()
			}
		}
	}
	return
}
