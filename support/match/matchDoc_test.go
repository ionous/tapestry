package match_test

import (
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
	"github.com/kr/pretty"
)

func TestSubDocument(t *testing.T) {
	test := func(name, str string, expect any) {
		var doc any
		p := charm.MakeParser(strings.NewReader(str))
		if e := p.ParseEof(
			match.DecodeTestDoc(func(a match.AsyncDoc) charm.State {
				doc = a.Content
				return charm.Finished()
			})); e != nil {
			t.Logf("failed %s with %s", name, e)
			t.Fail()
		} else if !reflect.DeepEqual(doc, expect) {
			pretty.Println(name, "mismatch, got:\n", doc)
			t.Fail()
		} else {
			t.Log("ok", name)
		}
	}
	for k, v := range subTests {
		test(k, v.str, v.expect)
	}
}

type subTest struct {
	str    string
	expect any
}

// just the one
var subTests = map[string]subTest{
	"assignText": {str: `
  FromText: "text"
Plain text again.`,
		expect: map[string]any{
			"FromText:": "text",
		}},
}
