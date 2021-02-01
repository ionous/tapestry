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
	// capture output into bytes
	var buf bytes.Buffer
	auto := run.Writer().(*print.AutoWriter)
	prev := auto.Target
	auto.Target = &buf

	run.ActivateDomain(t.Name, true)
	//
	if e := safe.Run(run, t.Test); e != nil {
		err = errutil.Fmt("ng! %s test encountered error: %s", t.Name, e)
	} else if res := buf.String(); res != t.Expect {
		err = errutil.Fmt("ng! %s got:  %q, want: %q", t.Name, res, t.Expect)
	} else {
		log.Printf("ok. test %s got %q", t.Name, res)
		auto.Target = prev
	}
	run.ActivateDomain(t.Name, false)
	return
}
