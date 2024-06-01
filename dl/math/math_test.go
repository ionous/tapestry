package math_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func TestMath(t *testing.T) {
	match := func(v float64, eval rt.NumEval) (err error) {
		var run testutil.PanicRuntime
		if n, e := eval.GetNum(&run); e != nil {
			err = e
		} else if n.Float() != v {
			err = errutil.Fmt("%v != %v (have != want)", n, v)
		}
		return
	}
	t.Run("Add", func(t *testing.T) {
		if e := match(11, &math.AddValue{A: literal.I(1), B: literal.I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Sub", func(t *testing.T) {
		if e := match(-9, &math.SubtractValue{A: literal.I(1), B: literal.I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Mul", func(t *testing.T) {
		if e := match(200, &math.MultiplyValue{A: literal.I(20), B: literal.I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div", func(t *testing.T) {
		if e := match(2, &math.DivideValue{A: literal.I(20), B: literal.I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div By Zero", func(t *testing.T) {
		if e := match(0, &math.DivideValue{A: literal.I(20), B: literal.I(0)}); e == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("Mod", func(t *testing.T) {
		if e := match(1, &math.ModValue{A: literal.I(3), B: literal.I(2)}); e != nil {
			t.Fatal(e)
		}
	})
}
