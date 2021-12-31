package core

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestCompareNumbers(t *testing.T) {
	test := func(a float64, op Comparator, b float64, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareNum{A: F(a), Is: op, B: F(b)}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test(10, &GreaterThan{}, 1, true)
	test(1, &GreaterThan{}, 10, false)
	test(8, &GreaterThan{}, 8, false)
	//
	test(10, &LessThan{}, 1, false)
	test(1, &LessThan{}, 10, true)
	test(8, &LessThan{}, 8, false)
	//
	test(10, &Equal{}, 1, false)
	test(1, &Equal{}, 10, false)
	test(8, &Equal{}, 8, true)
}

func TestCompareText(t *testing.T) {
	test := func(a string, op Comparator, b string, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareText{A: T(a), Is: op, B: T(b)}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test("Z", &GreaterThan{}, "A", true)
	test("A", &GreaterThan{}, "Z", false)
	//
	test("marzip", &LessThan{}, "marzipan", true)
	test("marzipan", &LessThan{}, "marzip", false)
	//
	test("bobby", &Equal{}, "bobby", true)
	test("bobby", &Equal{}, "phillipa", false)
}
