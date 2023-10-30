package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
)

func TestSig(t *testing.T) {
	// returns point of failure
	test := func(str string) (ret string, err error) {
		var p rift.SigParser
		if e := charm.Parse(&p, str); e != nil {
			err = e
		} else {
			ret, err = p.Signature()
		}
		return
	}
	fails := func(str string) (err error) {
		if v, e := test(str); e != nil {
			t.Log("ok failure:", str, e)
		} else {
			err = errutil.New(str, "expected error", v)
		}
		return
	}
	succeeds := func(str string) (err error) {
		if res, e := test(str); e != nil {
			err = errutil.New(e, "for:", str)
		} else if str != res {
			err = errutil.New(str, "unexpected value", res)
		} else {
			t.Log("ok success:", str)
		}
		return
	}
	if e := fails("a"); e != nil {
		t.Fatal(e)
	}
	if e := fails(" a"); e != nil {
		t.Fatal(e)
	}
	if e := fails("b "); e != nil {
		t.Fatal(e)
	}
	if e := fails("1a"); e != nil {
		t.Fatal(e)
	}
	if e := succeeds("a:"); e != nil {
		t.Fatal(e)
	}
	if e := succeeds("a:b:c:"); e != nil {
		t.Fatal(e)
	}
	if e := succeeds("and:more complex:keys_like_this:"); e != nil {
		t.Fatal(e)
	}
	if e := fails("a:b::c:"); e != nil {
		t.Fatal(e)
	}
}
