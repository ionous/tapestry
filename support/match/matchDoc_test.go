package match_test

import (
	"errors"
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
	"github.com/kr/pretty"
)

func TestSubDocument(t *testing.T) {
	test := func(name, str string, expect any) {
		var doc any
		var err charm.EndpointError
		if e := charm.ParseEof(str,
			match.DecodeDoc(func(q rune, content any) charm.State {
				doc = content
				return charm.Finished()
			}, false)); e != nil && !errors.As(e, &err) {
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

var subTests = map[string]subTest{
	"assignText": {str: `
  FromText: "text"
Plain text again.`,
		expect: map[string]any{
			"FromText:": "text",
		}},

	//		"assignExe": {str: `
	//	  - Say: "hello"
	//	  - Say: "world"
	//
	//	`, expect: []any{
	//			map[string]any{
	//				"Say:": "hello",
	//			},
	//			map[string]any{
	//				"Say:": "world",
	//			},
	//		}},
}
