package debug

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// an optional interface runtimes can implement for testing
type GreatExpectations interface {
	// gets and resets the most recent game output
	GetAccumulatedOutput() []string
}

func (op *ExpectOutput) Execute(run rt.Runtime) (err error) {
	return compareOutput(run, op.Output.GetLines())
}

func (op *ExpectText) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		err = compareOutput(run, []string{v.String()})
	}
	return
}

func (op *Expect) Execute(run rt.Runtime) (err error) {
	if condition, e := safe.GetBool(run, op.Value); e != nil {
		err = e
	} else if !condition.Bool() {
		err = errutil.New("expectation failed")
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

func compareLines(match, input []string) (err error) {
	if want, have := len(match), len(input); want != have {
		err = errutil.New("expected", want, "lines, have", have,
			"wanted:", strings.Join(match, "; "), "have:", strings.Join(input, "; "))
	} else {
		for i, w := range match {
			if h := input[i]; w != h {
				e := errutil.Fmt("line %v mismatched. wanted '%v' have '%v'", i, w, h)
				err = errutil.Append(err, e)
			} else if LogLevel.Str <= LoggingLevel_Debug {
				log.Println("~ ", w)
			}
		}
	}
	return
}
