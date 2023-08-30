package cmdcheck

import (
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
	var buf strings.Builder
	prevWriter := run.SetWriter(print.NewLineSentences(markup.ToText(&buf)))
	//
	if prevDomain, e := run.ActivateDomain(t.Domain); e != nil {
		err = e
	} else {
		if e := safe.RunAll(&checker{run, &buf}, t.Test); e != nil {
			err = errutil.Fmt("NG!  %s test encountered error: %s", t.Name, e)
		} else if res := buf.String(); res != t.Expect && len(t.Expect) > 0 {
			if eol := '\n'; strings.ContainsRune(res, eol) || strings.ContainsRune(t.Expect, eol) {
				err = errutil.Fmt("NG!  %s have:>>>\n%s<<<\nwant:>>>\n%s<<<", t.Name, res, t.Expect)
			} else {
				err = errutil.Fmt("NG!  %s have: %q want: %q", t.Name, res, t.Expect)
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
	buf *strings.Builder
}

func (c *checker) GetAccumulatedOutput() (ret []string) {
	str := c.buf.String()
	if cnt := len(str); cnt > 0 {
		ret = strings.FieldsFunc(str, func(r rune) bool { return r == newline })
		c.buf.Reset()
	}
	return
}

const newline = '\n'
