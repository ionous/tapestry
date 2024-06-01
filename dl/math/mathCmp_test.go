package math

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestCompareNumbers(t *testing.T) {
	test := func(a float64, op Comparison, b float64, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareNum{A: literal.F(a), Is: op, B: literal.F(b)}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test(10, C_Comparison_GreaterThan, 1, true)
	test(1, C_Comparison_GreaterThan, 10, false)
	test(8, C_Comparison_GreaterThan, 8, false)
	//
	test(10, C_Comparison_LessThan, 1, false)
	test(1, C_Comparison_LessThan, 10, true)
	test(8, C_Comparison_LessThan, 8, false)
	//
	test(10, C_Comparison_EqualTo, 1, false)
	test(1, C_Comparison_EqualTo, 10, false)
	test(8, C_Comparison_EqualTo, 8, true)
}

func TestCompareText(t *testing.T) {
	test := func(a string, op Comparison, b string, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareText{A: literal.T(a), Is: op, B: literal.T(b)}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test("Z", C_Comparison_GreaterThan, "A", true)
	test("A", C_Comparison_GreaterThan, "Z", false)
	//
	test("marzip", C_Comparison_LessThan, "marzipan", true)
	test("marzipan", C_Comparison_LessThan, "marzip", false)
	//
	test("bobby", C_Comparison_EqualTo, "bobby", true)
	test("bobby", C_Comparison_EqualTo, "phillipa", false)
}
