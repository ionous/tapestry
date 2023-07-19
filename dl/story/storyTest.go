package story

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *Test) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *Test) Weave(cat *weave.Catalog) (err error) {
	if name := lang.Normalize(op.TestName.String()); len(name) == 0 {
		errutil.New("test has empty name")
	} else {
		var req []string
		if n := op.DependsOn.String(); len(n) > 0 {
			req = []string{n}
		}
		if e := cat.DomainStart(name, req); e != nil {
			err = e
		} else {
			if e := WeaveStatements(cat, op.TestStatements); e != nil {
				err = e
			} else if len(op.Do) > 0 {
				err = cat.Schedule(weave.RequireAll, func(w *weave.Weaver) error {
					return w.Pin().AddCheck(name, nil, op.Do)
				})
			}

			if err == nil {
				err = cat.DomainEnd()
			}
		}
	}
	return
}
