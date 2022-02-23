package qna

import (
	"bytes"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

type CheckOutput struct {
	Name   string
	Domain string
	Expect string // all tests generate text right now; fix: need to handle comparison of literal values
	Test   []rt.Execute
}

func (t *CheckOutput) RunTest(run rt.Runtime) (err error) {
	log.Println("-- Checking:", t.Name, t.Domain)
	var buf bytes.Buffer
	prev := run.SetWriter(print.NewAutoWriter(&buf))
	if prev, e := run.ActivateDomain(t.Domain); e != nil {
		err = e
	} else {
		if e := safe.RunAll(run, t.Test); e != nil {
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

		if _, e := run.ActivateDomain(prev); e != nil {
			err = errutil.Append(err, errutil.New("couldnt restore domain", e))
		}

	}
	run.SetWriter(prev)
	return
}
