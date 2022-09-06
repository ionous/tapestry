package debug

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// an optional interface runtimes can implement for testing
type GreatExpectations interface {
	// gets and resets the most recent game output
	GetAccumulatedOutput() []string
}

func (op *ExpectLines) Execute(run rt.Runtime) (err error) {
	return compareLast(run, op.Lines.GetLines())
}

func (op *ExpectString) Execute(run rt.Runtime) (err error) {
	return compareLast(run, []string{op.String})
}

func (op *ExpectText) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		err = compareLast(run, []string{v.String()})
	}
	return
}

func (op *ExpectBool) Execute(run rt.Runtime) (err error) {
	if condition, e := safe.GetBool(run, op.Value); e != nil {
		err = e
	} else if !condition.Bool() {
		err = errutil.New("expectation failed")
	}
	return
}

func (op *ExpectNum) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetNumber(run, op.Value); e != nil {
		err = e
	} else {
		tolerance := 1e-3
		if op.Tolerance > 0.0 {
			tolerance = op.Tolerance
		}
		if want, have := op.Result, v.Float(); !op.Is.Compare().CompareFloat(have-want, tolerance) {
			err = errutil.Fmt("expectation failed: wanted '%v', have '%v'", want, have)
		}
	}
	return
}

// currently doing a matching of trailing lines rather than all lines output have to match.
func compareLast(run rt.Runtime, match []string) (err error) {
	if x, ok := run.(GreatExpectations); ok {
		err = compareAtLast(match, x.GetAccumulatedOutput())
	}
	return
}

func compareAtLast(match, input []string) (err error) {
	if want, have := len(match), len(input); want > have {
		err = errutil.New("expected", want, "lines, have", have)
	} else {
		input = input[have-want:]
		for i, w := range match {
			if h := input[i]; w != h {
				err = errutil.Fmt("line %v mismatched. wanted '%v' have '%v'", i, w, h)
				break
			}
		}
	}
	return
}
