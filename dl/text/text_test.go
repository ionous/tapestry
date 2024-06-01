package text_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	type baseRuntime struct {
		testutil.PanicRuntime
	}
	var run baseRuntime

	t.Run("is", func(t *testing.T) {
		if e := testTrue(t, &run, literal.B(true)); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: literal.B(false)}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("isEmpty", func(t *testing.T) {
		if e := testTrue(t, &run, &text.IsEmpty{Text: literal.T("")}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &text.IsEmpty{Text: literal.T("xxx")}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("includes", func(t *testing.T) {
		if e := testTrue(t, &run, &text.Includes{
			Text: literal.T("full"),
			Part: literal.T("ll"),
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &logic.Not{Test: &text.Includes{
			Text: literal.T("full"),
			Part: literal.T("bull"),
		}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("join", func(t *testing.T) {
		if e := testTrue(t, &run, &math.CompareText{
			A: &text.Join{Parts: []rt.TextEval{
				literal.T("one"), literal.T("two"), literal.T("three"),
			}},
			Compare: math.C_Comparison_EqualTo,
			B:       literal.T("onetwothree"),
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &math.CompareText{
			A: &text.Join{Sep: literal.T(" "), Parts: []rt.TextEval{
				literal.T("one"), literal.T("two"), literal.T("three"),
			}},
			Compare: math.C_Comparison_EqualTo,
			B:       literal.T("one two three"),
		}); e != nil {
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
