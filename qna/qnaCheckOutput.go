package qna

import (
	"bytes"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

type CheckOutput struct {
	Name   string
	Domain string
	Expect string // all tests generate text right now; fix: need to handle comparison of literal values
	Test   []rt.Execute
}

func (t *CheckOutput) RunTest(run rt.Runtime) (err error) {
	var buf bytes.Buffer
	prevWriter := run.SetWriter(print.NewLineSentences(markup.ToText(&buf)))
	//
	if prevDomain, e := run.ActivateDomain(t.Domain); e != nil {
		err = e
	} else {
		if e := safe.RunAll(&checker{run, &buf, 0}, t.Test); e != nil {
			err = errutil.Fmt("ng! %s test encountered error: %s", t.Name, e)
		} else if res := buf.String(); res != t.Expect && len(t.Expect) > 0 {
			if eol := '\n'; strings.ContainsRune(res, eol) || strings.ContainsRune(t.Expect, eol) {
				err = errutil.Fmt("ng! %s have:>>>\n%s<<<\nwant:>>>\n%s<<<", t.Name, res, t.Expect)
			} else {
				err = errutil.Fmt("ng! %s have: %q want: %q", t.Name, res, t.Expect)
			}
		}
		// restore even on a test mismatch
		if _, e := run.ActivateDomain(prevDomain); e != nil {
			err = errutil.Append(err, errutil.New("couldnt restore domain", e))
		}
	}
	run.SetWriter(prevWriter)
	return
}

type checker struct {
	rt.Runtime
	buf     *bytes.Buffer
	lastOut int
}

func (c *checker) GetAccumulatedOutput() (ret []string) {
	b := c.buf.Bytes()
	if cnt := len(b); cnt > c.lastOut {
		ret = strings.FieldsFunc(string(b[c.lastOut:]), func(r rune) bool { return r == newline })
		c.lastOut = cnt
	}
	return
}

const newline = '\n'
