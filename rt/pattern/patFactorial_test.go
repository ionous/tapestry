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
						&core.Assign{Var: N("num"),
							From: &core.FromNum{
								Val: &core.ProductOf{
									A: V("num"),
									B: &core.CallPattern{
										Pattern: P("factorial"),
										Arguments: core.NamedArgs(
											"num", &core.FromNum{
												Val: &core.DiffOf{
													A: V("num"),
													B: I(1),
												},
											},
										)}}}}),
				}, {
					Filter: &core.CompareNum{
						A:  V("num"),
						Is: &core.Equal{},
						B:  I(0),
					},
					Execute: core.MakeActivity(
						&core.Assign{Var: N("num"),
							From: &core.FromNum{
								Val: I(1),
							}},
					),
				}}},
		}}
	// determine the factorial of the number 3
	det := core.CallPattern{
		Pattern: P("factorial"),
		Arguments: core.NamedArgs(
			"num", &core.FromNum{
				Val: I(3),
			}),
	}
	if v, e := safe.GetNumber(&run, &det); e != nil {
		t.Fatal(e)
	} else if got, want := v.Int(), 3*(2*(1*1)); got != want {
		t.Fatal("mismatch: expected:", want, "have:", got)
	} else {
		t.Log("factorial okay", got)
	}
}
