package pattern_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
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
				Fields: []g.Field{
					{Name: "num", Affinity: affine.Number},
				},
				Rules: []rt.Rule{{
					Execute: core.MakeActivity(
						&core.SetValue{
							Target: core.Variable("num"),
							Value: core.AssignFromNumber(&core.ProductOf{
								A: GetVariable("num"),
								B: &core.CallPattern{
									Pattern: P("factorial"),
									Arguments: []core.Arg{{
										Name: "num",
										Value: core.AssignFromNumber(&core.DiffOf{
											A: GetVariable("num"),
											B: I(1),
										})}}}})}),
				}, {
					Filter: &core.CompareNum{
						A:  GetVariable("num"),
						Is: core.Equal,
						B:  I(0),
					},
					Execute: core.MakeActivity(
						&core.SetValue{
							Target: core.Variable("num"),
							Value:  core.AssignFromNumber(I(1))},
					),
				}}},
		}}
	// determine the factorial of the number 3
	det := core.CallPattern{
		Pattern: P("factorial"),
		Arguments: []core.Arg{{
			Name:  "num",
			Value: core.AssignFromNumber(I(3)),
		}},
	}
	if v, e := safe.GetNumber(&run, &det); e != nil {
		t.Fatal(e)
	} else if got, want := v.Int(), 3*(2*(1*1)); got != want {
		t.Fatal("mismatch: expected:", want, "have:", got)
	} else {
		t.Log("factorial okay", got)
	}
}
