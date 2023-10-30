package rift_test

import (
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
	} else if e := matchSeq(t,
		// fix? trailing space at the end fails
		"a single with newline",
		`- 5
`, []any{5.0}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t,
		"consistent indents", `
  - 5
  - 10
  - 12`,
		[]any{5.0, 10.0, 12.0}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t,
		"nested sub sequence", `
- - 5
`, []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t,
		"new line sub sequence", `
- 
  - 5
`, []any{[]any{5.0}}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t,
		"nil values", `
  -
  -
  -`,
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t,
		"nil value trailing newline", `
  -
  -
  -
`,
		[]any{nil, nil, nil}); e != nil {
		t.Fatal(e)
	}
}

func matchSeq(t *testing.T, name, str string, want any) (err error) {
	if have, e := testSeq(str); e != nil {
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
	var p *rift.Sequence
	ws := rift.OptionalWhitespace()
	if e := charm.Parse(charm.Step(ws, charm.Statement("", func(r rune) (ret charm.State) {
		p = rift.NewSequence(ws.Indent)
		return p.NewRune(r)
	})), str); e != nil {
		err = e
	} else if vs, e := p.GetSequence(); e != nil {
		err = e
	} else {
		ret = condense(vs)
	}
	return
}

func condense(vs []rift.Value) []any {
	out := make([]any, 0, len(vs))
	for _, v := range vs {
		el := v.Result
		if sub, ok := el.([]rift.Value); ok {
			el = condense(sub)
		}
		// splt comment only, but not nil
		out = append(out, el)
	}
	return out
}

// if e := fails("a"); e != nil {
// 	t.Fatal(e)
// }
// if e := fails(" a"); e != nil {
// 	t.Fatal(e)
// }
// if e := fails("b "); e != nil {
// 	t.Fatal(e)
// }
// if e := fails("1a"); e != nil {
// 	t.Fatal(e)
// }
// if e := succeeds("a:"); e != nil {
// 	t.Fatal(e)
// }
// if e := succeeds("a:b:c:"); e != nil {
// 	t.Fatal(e)
// }
// if e := succeeds("and:more complex:keys_like_this:"); e != nil {
// 	t.Fatal(e)
// }
// if e := fails("a:b::c:"); e != nil {
// 	t.Fatal(e)
// }
