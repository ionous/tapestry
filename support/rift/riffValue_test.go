package rift_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
)

func TestSimpleScalars(t *testing.T) {
	if e := match(t,
		"number",
		testValue(`5.4`),
		5.4); e != nil {
		t.Fatal(e)
	}
	// string
	if e := match(t,
		"string",
		testValue(`"5.4"`),
		"5.4"); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"bool",
		testValue(`true`),
		true); e != nil {
		t.Fatal(e)
	}
	// // something else
	// if e := testValue(`beep`).(error); e != nil {
	// 	t.Log("ok failure:", e)
	// }
}

// returns point of failure
func testValue(str string) func() any {
	str = strings.TrimSpace(str)
	return func() (ret any) {
		var h rift.History
		if e := charm.Parse(str, rift.NewValue(&h, 0, func(v any) (_ error) {
			ret = v
			return
		})); e != nil {
			ret = errutil.Fmt("%v %q", e, str)
		} else if e := h.PopAll(); e != nil {
			ret = e
		}
		return
	}
}
