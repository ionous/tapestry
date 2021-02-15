package check

import (
	"bytes"
	"log"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type CheckOutput struct {
	Name, Expect string
	Test         *core.Activity
}

func (t *CheckOutput) RunTest(run rt.Runtime) (err error) {
	var buf bytes.Buffer
	prev := run.SetWriter(print.NewAutoWriter(&buf))
	run.ActivateDomain(t.Name, true)
	{
		if e := safe.Run(run, t.Test); e != nil {
			err = errutil.Fmt("ng! %s test encountered error: %s", t.Name, e)
		} else if res := buf.String(); res != t.Expect {
			err = errutil.Fmt("ng! %s got:  %q, want: %q", t.Name, res, t.Expect)
		} else {
			log.Printf("ok. test %s got %q", t.Name, res)
		}
	}
	run.ActivateDomain(t.Name, false)
	run.SetWriter(prev)
	return
}
