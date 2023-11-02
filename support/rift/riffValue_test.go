package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestSimpleScalars(t *testing.T) {
	// returns point of failure
	test := func(str string) (ret any, err error) {
		var h rift.History
		if e := charm.Parse(str, rift.NewValue(&h, 0, func(v any) (_ error) {
			ret = v
			return
		})); e != nil {
			err = errutil.Fmt("%v %q", e, str)
		} else {
			err = h.PopAll()
		}
		return
	}

	match := func(str string, want any) (err error) {
		if have, e := test(str); e != nil {
			err = e
		} else if d := pretty.Diff(want, have); len(d) != 0 {
			err = errutil.Fmt("mismatched want: %v have: %v diff: %s", want, have, d)
		} else {
			t.Logf("ok success: %T %v", have, have)
		}
		return
	}
	// number
	if e := match(`5.4`, 5.4); e != nil {
		t.Fatal(e)
	}
	// string
	if e := match(`"5.4"`, "5.4"); e != nil {
		t.Fatal(e)
	}
	// bool
	if e := match(`true`, true); e != nil {
		t.Fatal(e)
	}
	// something else
	if v, e := test(`beep`); e == nil {
		t.Fatal("expected failure; got:", v)
	} else {
		t.Log("ok failure:", e)
	}
}
