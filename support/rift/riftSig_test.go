package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
)

func TestSig(t *testing.T) {
	// returns point of failure
	test := func(str string) (ret string, err error) {
		var doc rift.Document
		if e := doc.Parse(str, rift.NewSignature(&doc, 0, func(str string) (_ error) {
			ret = str
			return
		})); e != nil {
			err = e
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
			err = errutil.New(str, "unexpected result", res)
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
