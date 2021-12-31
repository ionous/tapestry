package core

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func TestMath(t *testing.T) {
	match := func(v float64, eval rt.NumberEval) (err error) {
		var run testutil.PanicRuntime
		if n, e := eval.GetNumber(&run); e != nil {
			err = e
		} else if n.Float() != v {
			err = errutil.Fmt("%v != %v (have != want)", n, v)
		}
		return
	}
	t.Run("Add", func(t *testing.T) {
		if e := match(11, &SumOf{A: I(1), B: I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Sub", func(t *testing.T) {
		if e := match(-9, &DiffOf{A: I(1), B: I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Mul", func(t *testing.T) {
		if e := match(200, &ProductOf{A: I(20), B: I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div", func(t *testing.T) {
		if e := match(2, &QuotientOf{A: I(20), B: I(10)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div By Zero", func(t *testing.T) {
		if e := match(0, &QuotientOf{A: I(20), B: I(0)}); e == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("Mod", func(t *testing.T) {
		if e := match(1, &RemainderOf{A: I(3), B: I(2)}); e != nil {
			t.Fatal(e)
		}
	})
}
