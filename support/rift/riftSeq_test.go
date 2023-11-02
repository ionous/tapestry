package rift_test

import (
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestSeq(t *testing.T) {
	if e := matchSeq(t,
		"a single value",
		`- 5`, []any{5.0}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"a single with newline",
		`- 5
`, []any{5.0}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"consistent indents", `
      - 5
      - 10
      - 12`,
		[]any{5.0, 10.0, 12.0}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"nested sub sequence",
		`- - 5
`, []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"new line sub sequence", `
    -
      - 5
    `, []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"nil values", `
      -
      -
      -`,
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"nil value trailing newline", `
  -
  -
  -
`,
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	}
	if e := matchSeq(t,
		"continuing sub sequence ", `
    -
      - 5
    - 6
    `, []any{[]any{5.0}, 6.0}); e != nil {
		t.Fatal(e)
	}
}

func matchSeq(t *testing.T, name, str string, want any) (err error) {
	if strings.HasPrefix(name, `x `) {
		// commenting out tests causes go fmt to replace spaces with tabs. *sigh*
		t.Log("skipping", name)
	} else if have, e := testSeq(str); e != nil {
		err = errutil.Fmt("ng failed %q %v", name, e)
	} else if d := pretty.Diff(want, have); len(d) != 0 {
		err = errutil.Fmt("ng mismatched %q want: %v have: %v diff: %v",
			name, want, have, d)
	} else {
		t.Logf("ok success: %q %T %v", name, have, have)
	}
	return
}

func testSeq(str string) (ret []any, err error) {
	var h rift.History
	if e := charm.Parse(str, rift.OptionalSpaces("test", 0, func(indent int) charm.State {
		return rift.NewSequence(&h, indent, func(vs []any) (_ error) {
			ret = vs
			return
		})
	})); e != nil {
		err = e
	} else {
		err = h.PopAll()
	}
	return
}

// could be put in a charm helper package
func init() {
	charm.StateName = func(n charm.State) (ret string) {
		if s, ok := n.(interface{ String() string }); ok {
			ret = s.String()
		} else if n == nil {
			ret = "null"
		} else {
			ret = reflect.TypeOf(n).Elem().Name()
		}
		return
	}
}
