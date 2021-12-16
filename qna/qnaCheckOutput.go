package check

import (
	"bytes"
	"log"
	"strings"

	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (t *CheckOutput) RunTest(run rt.Runtime) (err error) {
	var buf bytes.Buffer
	prev := run.SetWriter(print.NewAutoWriter(&buf))
	run.ActivateDomain(t.Name, true)
	{
		if e := safe.Run(run, &t.Test); e != nil {
			err = errutil.Fmt("ng! %s test encountered error: %s", t.Name, e)
		} else if res := buf.String(); res != t.Expect {
			if eol := '\n'; strings.ContainsRune(res, eol) || strings.ContainsRune(t.Expect, eol) {
				err = errutil.Fmt("ng! %s have:>>>\n%s<<<\nwant:>>>\n%s<<<", t.Name, res, t.Expect)
			} else {
				err = errutil.Fmt("ng! %s have: %q want: %q", t.Name, res, t.Expect)
			}
		} else {
			log.Printf("ok. test %s got %q", t.Name, res)
		}
	}
	run.ActivateDomain(t.Name, false)
	run.SetWriter(prev)
	return
}
