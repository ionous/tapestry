package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
)

func TestSeq(t *testing.T) {
	if e := match(t,
		"a single value",
		testSeq(
			`- 5`), []any{5.0}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"a single with newline",
		testSeq(
			`- 5
`), []any{5.0}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"consistent indents",
		testSeq(`
      - 5
      - 10
      - 12`),
		[]any{5.0, 10.0, 12.0}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"nested sub sequence",
		testSeq(
			`- - 5
`), []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"new line sub sequence",
		testSeq(`
    -
      - 5
    `), []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"nil values",
		testSeq(`
      -
      -
      -`),
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"nil value trailing newline",
		testSeq(`
  -
  -
  -
`),
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	}
	if e := match(t,
		"continuing sub sequence ",
		testSeq(`
    -
      - 5
    - 6
    `), []any{[]any{5.0}, 6.0}); e != nil {
		t.Fatal(e)
	}
}

func testSeq(str string) func() (ret any) {
	return func() (ret any) {
		var h rift.History
		if e := charm.Parse(str, rift.OptionalSpaces("test", 0, func(indent int) charm.State {
			return rift.NewSequence(&h, indent, func(vs []any) (_ error) {
				ret = vs
				return
			})
		})); e != nil {
			ret = e
		} else if e := h.PopAll(); e != nil {
			ret = e
		}
		return
	}
}
