package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestSeq(t *testing.T) {
	if e := matchSeq(t, "- 5", []any{5.0}); e != nil {
		t.Fatal(e)
	} else if e := matchSeq(t, list, []any{5.0, 10.0, 12.0}); e != nil {
		t.Fatal(e)
	}
}

var list string = `- 5
- 10
- 12
`

func matchSeq(t *testing.T, str string, want any) (err error) {
	if have, e := testSeq(str); e != nil {
		err = e
	} else if d := pretty.Diff(want, have); len(d) != 0 {
		err = errutil.Fmt("mismatched want: %v have: %v diff: %v", want, have, d)
	} else {
		t.Logf("ok success: %T %v", have, have)
	}
	return
}

func testSeq(str string) (ret []any, err error) {
	var p rift.SeqParser
	if e := charm.Parse(&p, str); e != nil {
		err = e
	} else if vs, e := p.GetValue(); e != nil {
		err = e
	} else {
		ret = condense(vs)
	}
	return
}

func condense(vs []rift.Value) []any {
	out := make([]any, 0, len(vs))
	for _, v := range vs {
		// sjip comment only, but not nil
		out = append(out, v.Result)
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
