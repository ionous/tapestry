package core

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	var run baseRuntime

	t.Run("is", func(t *testing.T) {
		if e := testTrue(t, &run, B(true)); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &Not{Test: B(false)}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("isEmpty", func(t *testing.T) {
		if e := testTrue(t, &run, &IsEmpty{Text: T("")}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &Not{Test: &IsEmpty{Text: T("xxx")}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("includes", func(t *testing.T) {
		if e := testTrue(t, &run, &Includes{
			Text: T("full"),
			Part: T("ll"),
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &Not{Test: &Includes{
			Text: T("full"),
			Part: T("bull"),
		}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("join", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{
			A: &Join{Parts: []rt.TextEval{
				T("one"), T("two"), T("three"),
			}},
			Is: C_Comparison_EqualTo,
			B:  T("onetwothree"),
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{
			A: &Join{Sep: T(" "), Parts: []rt.TextEval{
				T("one"), T("two"), T("three"),
			}},
			Is: C_Comparison_EqualTo,
			B:  T("one two three"),
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
