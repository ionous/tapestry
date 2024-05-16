package core_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
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
	this := core.T("this")
	that := core.T("that")
	nothing := core.T("nothing")

	t.Run("exists", func(t *testing.T) {
		if e := testTrue(t, &run, &assign.ObjectDot{Name: this}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.Not{Test: &assign.ObjectDot{Name: nothing}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &core.CompareText{
			A: &core.KindOf{Object: this}, Is: core.C_Comparison_EqualTo, B: core.T(base)}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.CompareText{
			A: &core.KindOf{Object: that}, Is: core.C_Comparison_EqualTo, B: core.T(derived)}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &core.IsKindOf{Object: this, Kind: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.IsKindOf{Object: that, Kind: base}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(t, &run, &core.IsKindOf{Object: that, Kind: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.Not{Test: &core.IsKindOf{Object: this, Kind: derived}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &core.IsExactKindOf{Object: this, Kind: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.Not{Test: &core.IsExactKindOf{Object: that, Kind: base}}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.IsExactKindOf{Object: that, Kind: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &core.Not{Test: &core.IsExactKindOf{Object: this, Kind: derived}}); e != nil {
			t.Fatal(e)
		}
	})
}

func testTrue(t *testing.T, run rt.Runtime, eval rt.BoolEval) (err error) {
	if ok, e := safe.GetBool(run, eval); e != nil {
		err = e
	} else if !ok.Bool() {
		err = errutil.New("expected true", pretty.Sprint(eval))
	}
	return
}
