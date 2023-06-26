package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *Test) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

func (op *Test) Schedule(cat *weave.Catalog) (err error) {
	if name := op.TestName.String(); len(name) == 0 {
		errutil.New("test has empty name")
	} else {
		var req []string
		if n := op.DependsOn.String(); len(n) > 0 {
			req = []string{n}
		}
		if e := cat.AssertDomainStart(name, req); e != nil {
			err = e
		} else {
			if e := story.ScheduleStatements(cat, op.TestStatements); e != nil {
				err = e
			} else if len(op.Do) > 0 {
				if e := cat.AssertCheck(name, op.Do, nil); e != nil {
					err = e
				}
			}

			if err == nil {
				err = cat.AssertDomainEnd()
			}
		}
	}
	return
}
