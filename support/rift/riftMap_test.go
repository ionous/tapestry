package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
)

// keys that start with t or f need special handleing
func TestBoolKeys(t *testing.T) {
	if e := match(t,
		"nested value",
		testValue(`
true: false
false: true
`),
		rift.MapValues{
			{"true:", false},
			{"false:", true},
		}); e != nil {
		t.Fatal(e)
	}
}

func TestMap(t *testing.T) {
	if e := match(t,
		"single value",
		testMap(
			`name: "Sammy Sosa"`),
		rift.MapValues{{
			"name:", "Sammy Sosa",
		}}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"several values",
		testMap(
			`name: "Sammy Sosa"
hr:   63
avg:  true`),
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"nested map",
		testMap(
			`name: "Sammy Sosa"
hr:   63
avg:  true`),
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		}); e != nil {
		t.Fatal(e)
	}
	// make sure we can also parse that block as a value
	if e := match(t,
		"nested value",
		testValue(`
name: "Sammy Sosa"
hr:   63
avg:  true`),
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		}); e != nil {
		t.Fatal(e)
	}
	// make sure we can also parse that block as a value
	if e := match(t,
		"nested value",
		testValue(`
true: false
false: true`),
		rift.MapValues{
			{"true:", false},
			{"false:", true},
		}); e != nil {
		t.Fatal(e)
	}
}

func testMap(str string) func() any {
	return func() (ret any) {
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
}
