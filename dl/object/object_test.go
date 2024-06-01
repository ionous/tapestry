package object_test

import (
	"fmt"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

// test some simple functionality of the object commands using a mock runtime
func TestObjects(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects
	type Base struct{}
	type Derived struct{ Base }

	kinds.AddKinds((*Base)(nil), (*Derived)(nil))
	objs.AddObjects(kinds.Kind("base"), "this")
	objs.AddObjects(kinds.Kind("derived"), "that")

	run := testutil.Runtime{
		Kinds:     &kinds,
		ObjectMap: objs,
	}
	const base, derived = "base", "derived"
	this := literal.T("this")
	that := literal.T("that")
	nothing := literal.T("nothing")

	t.Run("exists", func(t *testing.T) {
		if e := testTrue(t, &run, &object.ObjectDot{Name: this}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &object.ObjectDot{Name: nothing}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &math.CompareText{
			A: &object.KindOf{Object: this}, Compare: math.C_Comparison_EqualTo, B: literal.T(base)}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &math.CompareText{
			A: &object.KindOf{Object: that}, Compare: math.C_Comparison_EqualTo, B: literal.T(derived)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &object.IsKindOf{Object: this, Kind: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &object.IsKindOf{Object: that, Kind: base}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(t, &run, &object.IsKindOf{Object: that, Kind: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &object.IsKindOf{Object: this, Kind: derived}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &object.IsExactKindOf{Object: this, Kind: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &object.IsExactKindOf{Object: that, Kind: base}}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &object.IsExactKindOf{Object: that, Kind: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &object.IsExactKindOf{Object: this, Kind: derived}}); e != nil {
			t.Fatal(e)
		}
	})
}

func testTrue(t *testing.T, run rt.Runtime, eval rt.BoolEval) (err error) {
	if ok, e := safe.GetBool(run, eval); e != nil {
		err = e
	} else if !ok.Bool() {
		err = fmt.Errorf("expected true %s", pretty.Sprint(eval))
	}
	return
}
