package pattern_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

// TestFactorial of the number 3 to verify pattern recursion works.
func TestFactorial(t *testing.T) {
	type Factorial struct {
		Num float64
	}
	var kinds testutil.Kinds
	kinds.AddKinds((*Factorial)(nil))
	// rules are run in reverse order.
	run := testpat.Runtime{
		Runtime: testutil.Runtime{
			Kinds: &kinds,
		},
		Map: testpat.Map{
			"factorial": &testpat.Pattern{
				Name:   "factorial",
				Labels: []string{"num"},
				Return: "num",
				Fields: []rt.Field{
					{Name: "num", Affinity: affine.Num},
				},
				Rules: []rt.Rule{
					{Exe: debug.FactorialDefaultRule},
					{Exe: debug.FactorialDecreaseRule},
				},
			}}}

	// determine the factorial of the number 3
	det := call.CallPattern{
		PatternName: ("factorial"),
		Arguments: []call.Arg{{
			Name:  "num",
			Value: &call.FromNum{Value: literal.I(3)},
		}},
	}
	if v, e := safe.GetNum(&run, &det); e != nil {
		t.Fatal(e)
	} else if got, want := v.Int(), 3*(2*(1*1)); got != want {
		t.Fatal("mismatch: expected:", want, "have:", got)
	} else {
		t.Log("factorial okay", got)
	}
}
