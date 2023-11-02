package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
)

func TestMap(t *testing.T) {
	if e := match(t,
		"a single value",
		testMap(
			`name: "Sammy Sosa"`),
		rift.MapValues{{
			"name:", "Sammy Sosa",
		}}); e != nil {
		t.Fatal(e)
	}
}

func testMap(str string) (ret any) {
	var h rift.History
	if e := charm.Parse(str, rift.NewMapping(&h, 0, func(vs rift.MapValues) (_ error) {
		ret = vs
		return
	})); e != nil {
		ret = e
	} else if e := h.PopAll(); e != nil {
		ret = e
	}
	return
}
