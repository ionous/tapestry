package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestValue(t *testing.T) {
	// returns point of failure
	test := func(str string) (ret any, err error) {
		var p rift.ValueParser
		if e := charm.Parse(&p, str); e != nil {
			err = errutil.Fmt("%v %q", e, str)
		} else if v, e := p.GetValue(); e != nil {
			err = errutil.Fmt("%v %q", e, str)
		} else {
			ret = v.Result
		}
		return
	}

	match := func(str string, want any) (err error) {
		if have, e := test(str); e != nil {
			err = e
		} else if d := pretty.Diff(want, have); len(d) != 0 {
			err = errutil.Fmt("mismatched want: %s have: %s diff: %s", want, have, d)
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
