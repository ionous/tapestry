package core_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/idl/all"
	"git.sr.ht/~ionous/iffy/idl/core"
	"git.sr.ht/~ionous/iffy/idl/rtx"
	"git.sr.ht/~ionous/iffy/test/testutil"
	capnp "zombiezen.com/go/capnproto2"
)

func TestAllTrue(t *testing.T) {
	all.RegisterTypes()
	run := &testutil.PanicRuntime{}
	if _, seg, e := capnp.NewMessage(capnp.SingleSegment(nil)); e != nil {
		t.Fatal(e)
	} else if x, e := core.NewAllTrue(seg); e != nil {
		t.Fatal(e)
	} else {

		for i := 0; i < 3; i++ {
			// lists are new'd with int32, but indexed by int. *sigh*
			if tests, e := x.NewTest(int32(i)); e != nil {
				t.Fatal(e)
			} else {
				for j := 0; j < i; j++ {
					if bptr, e := rtx.NewBoolEval(seg); e != nil {
						t.Fatal(e)
					} else if b, e := core.NewBoolValue(seg); e != nil {
						t.Fatal(e)
					} else {
						b.SetBool(true)
						bptr.SetEvalPtr(b.ToPtr())
						tests.Set(j, bptr)
					}
				}
				//
				if e := x.SetTest(tests); e != nil {
					t.Fatal(e)
				} else if ok, e := x.GetBool(run); e != nil {
					t.Fatal(e)
				} else if !ok.Bool() {
					t.Fatal("expected success")
				}
			}
		}
	}
	// turn one false.
	// l.vals[1] = false
	// test := &AllTrue{evals}
	// if ok, e := safe.GetBool(run, test); e != nil {
	// 	t.Fatal(e)
	// } else if ok.Bool() {
	// 	t.Fatal("expected failure")
	// } else if l.asks != 2 {
	// 	t.Fatal("expected only two got tested", l.asks)
	// }
}
