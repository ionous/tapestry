package debug

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// an optional interface runtimes can implement for testing
type GreatExpectations interface {
	// gets and resets the most recent game output
	GetAccumulatedOutput() []string
}

// separate a string into separate lines whenever a linefeed ( \n ) is detected.
func SplitLines(str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool { return r == '\n' })
}

func (op *ExpectText) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		err = compareOutput(run, SplitLines(v.String()))
	}
	return
}

func (op *Expect) Execute(run rt.Runtime) (err error) {
	if condition, e := safe.GetBool(run, op.Value); e != nil {
		err = e
	} else {
		str := compact.JoinComment(op.Markup)
		if !condition.Bool() {
			err = errutil.New("expectation failed", str)
			log.Println("ng:", str)
		} else {
			log.Println("ok:", str)
		}
	}
	return
}

// currently doing a matching of trailing lines rather than all lines output have to match.
func compareOutput(run rt.Runtime, match []string) (err error) {
	if x, ok := run.(GreatExpectations); ok {
		err = compareLines(match, x.GetAccumulatedOutput())
	}
	return
}

func compareLines(want, have []string) (err error) {
	if wcnt, hcnt := len(want), len(have); wcnt > hcnt {
		err = errutil.New("expected", want, "lines, have", have,
			"wanted:", strings.Join(want, "; "), "have:", strings.Join(have, "; "))
	} else {
		//
		var elided bool
		for i, cnt := 0, len(want); i < cnt; i++ {
			w, h := want[i], have[i]
			if elided = strings.HasSuffix(w, "..."); elided {
				// only want to match up to the ellipses
				if max := len(w) - 3; len(h) >= max {
					w, h = w[:max], h[:max]
				}
			}
			if w != h {
				e := errutil.Fmt("line %v mismatched. wanted '%v' have '%v'", i, want[i], have[i])
				err = errutil.Append(err, e)
			} else if LogLevel <= C_LoggingLevel_Debug {
				log.Println("~ ", want[i])
			}
		}
		if hcnt > wcnt && !elided {
			err = errutil.New("unexpected trailing text:", have[wcnt:])
		}
	}
	return
}
