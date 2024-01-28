package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineTest) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineTest) Weave(cat *weave.Catalog) (err error) {
	if name := inflect.Normalize(op.TestName); len(name) == 0 {
		errutil.New("test has empty name")
	} else if dependsOn, e := safe.GetOptionalTexts(cat.Runtime(), op.SceneNames, nil); e != nil {
		err = e
	} else {
		if e := cat.DomainStart(name, dependsOn.Strings()); e != nil {
			err = e
		} else {
			if e := WeaveStatements(cat, op.Statements); e != nil {
				err = e
			} else if len(op.Exe) > 0 {
				err = cat.Schedule(weave.RequireAll, func(w *weave.Weaver) error {
					return w.Pin().AddCheck(name, nil, op.Exe)
				})
			}
			if err == nil {
				err = cat.DomainEnd()
			}
		}
	}
	return
}
