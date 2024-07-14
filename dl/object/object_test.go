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

// test object commands using a mock runtime
func TestObjectCommands(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects
	type Base struct{}
	type Derived struct{ Base }

	kinds.AddKinds((*Base)(nil), (*Derived)(nil))
	objs.AddObjects(kinds.Kind("base"), "this")
	objs.AddObjects(kinds.Kind("derived"), "that")

	run := testutil.Runtime{
		Kinds:   &kinds,
		Objects: objs,
	}
	base := literal.T("base")
	derived := literal.T("derived")
	this := object.Object("this")
	that := object.Object("that")
	nothing := object.Object("nothing")

	t.Run("exists", func(t *testing.T) {
		if e := testTrue(&run, this); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &logic.Not{Test: nothing}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if e := testTrue(&run, &math.CompareText{
			A: &object.KindOf{Target: this}, Compare: math.C_Comparison_EqualTo, B: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &math.CompareText{
			A: &object.KindOf{Target: that}, Compare: math.C_Comparison_EqualTo, B: derived}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(&run, &object.IsKindOf{Target: this, KindName: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &object.IsKindOf{Target: that, KindName: base}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(&run, &object.IsKindOf{Target: that, KindName: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &logic.Not{Test: &object.IsKindOf{Target: this, KindName: derived}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(&run, &object.IsExactKindOf{Target: this, KindName: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &logic.Not{Test: &object.IsExactKindOf{Target: that, KindName: base}}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &object.IsExactKindOf{Target: that, KindName: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(&run, &logic.Not{Test: &object.IsExactKindOf{Target: this, KindName: derived}}); e != nil {
			t.Fatal(e)
		}
	})
}

func testTrue(run rt.Runtime, eval rt.BoolEval) (err error) {
	if ok, e := safe.GetBool(run, eval); e != nil {
		err = e
	} else if !ok.Bool() {
		err = fmt.Errorf("expected true %s", pretty.Sprint(eval))
	}
	return
}
