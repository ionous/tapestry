package cmdcheck

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

// takes a play time so that it can run parser commands when requeted
func RunTest(run rt.Runtime, t query.CheckData) (err error) {
	var buf strings.Builder
	prevWriter := run.SetWriter(print.NewLineSentences(markup.ToText(&buf)))
	//
	if e := run.ActivateDomain(t.Domain); e != nil {
		err = e
	} else {
		if e := safe.RunAll(&checker{run, &buf}, t.Test); e != nil && !wasQuit(e) && !wasInterrupt(e) {
			err = errutil.Fmt("NG!  %q encountered error: %s", t.Name, e)
		} else if res := buf.String(); res != t.Expect && len(t.Expect) > 0 {
			if eol := '\n'; strings.ContainsRune(res, eol) || strings.ContainsRune(t.Expect, eol) {
				err = errutil.Fmt("NG!  %q have:>>>\n%s<<<\nwant:>>>\n%s<<<", t.Name, res, t.Expect)
			} else {
				err = errutil.Fmt("NG!  %q have: %q want: %q", t.Name, res, t.Expect)
			}
		}
		// restore even on a test mismatch
		if e := run.ActivateDomain(""); e != nil {
			err = errutil.Append(err, errutil.New("couldnt restore domain", e))
		}
	}
	run.SetWriter(prevWriter)
	return
}

// fix: better might be an "Expect error:" that matches partial error return strings
// or "Expect quit/signal:" specifically for this case
func wasQuit(e error) bool {
	var sig game.Signal // if the game was quit, override the error if output remains
	return errors.As(e, &sig) && sig == game.SignalQuit
}

// ignores break and continue statements during tests
// ( some tests only exist to define a scene; and continue in the body of their test )
func wasInterrupt(e error) bool {
	var i logic.DoInterrupt
	return errors.As(e, &i)
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
