package pattern_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

// TestFactorial of the number 3 to verify pattern recursion works.
func TestFactorial(t *testing.T) {
	type Factorial struct {
		Num float64
	}
	var kinds testutil.Kinds
	kinds.AddKinds((*Factorial)(nil))
	// rules are run in reverse order.
	run := testutil.Runtime{
		Kinds: &kinds,

		PatternMap: testutil.PatternMap{
			"factorial": &pattern.Pattern{
				Name:   "factorial",
				Labels: []string{"num"},
				Return: "num",
				Fields: []g.Field{
					{Name: "num", Affinity: affine.Number},
				},
				Rules: []*pattern.Rule{{
					Execute: core.NewActivity(
						&core.Assign{Var: N("num"),
							From: &core.FromNum{
								&core.ProductOf{
									A: V("num"),
									B: &pattern.Determine{
										Pattern: "factorial",
										Arguments: core.NamedArgs(
											"num", &core.FromNum{
												&core.DiffOf{
													V("num"),
													I(1),
												},
											},
										)}}}}),
				}, {
					Filter: &core.CompareNum{
						V("num"),
						&core.EqualTo{},
						I(0),
					},
					Execute: core.NewActivity(
						&core.Assign{Var: N("num"),
							From: &core.FromNum{
								I(1),
							}},
					),
				}}},
		}}
	// determine the factorial of the number 3
	det := pattern.Determine{
		Pattern: "factorial",
		Arguments: core.NamedArgs(
			"num", &core.FromNum{
				I(3),
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
