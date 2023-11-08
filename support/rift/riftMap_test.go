package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

// keys that start with t or f need special handleing
func xTestMap(t *testing.T) {
	testMap(t,
		// -----------
		"test keys with boolean names", `
true: false
false: true`,
		rift.MapValues{
			{"true:", false},
			{"false:", true},
		},

		// -----------
		"test single value", `
name: "Sammy Sosa"`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
		},
		// -----------
		"test split line", `
name:
  "Sammy Sosa"`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
		},

		// -----------
		"test several values", `
name: "Sammy Sosa"
hr:   63
avg:  true`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		},
		// -----------
		"test nested map", `
name: "Sammy Sosa"
hr:   63
avg:  true`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		},
	)

	// make sure we can also parse that block as a value
	testValue(t,
		"test nested value", `
name: "Sammy Sosa"
hr:   63
avg:  true`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		})
}

// name of test, source string, expected result
func testMap(t *testing.T, nameInputExpect ...any) {
	for i, cnt := 0, len(nameInputExpect); i < cnt; i += 3 {
		name, input, expect := nameInputExpect[0+i].(string), nameInputExpect[1+i].(string), nameInputExpect[2+i]
		if strings.HasPrefix(name, `x `) {
			// commenting out tests causes go fmt to replace spaces with tabs. *sigh*
			t.Log("skipping", name)
		} else {
			var res any
			var doc rift.Document
			str := strings.TrimLeftFunc(input, unicode.IsSpace)
			mapping := rift.NewMapping(&doc, "", 0)
			if e := doc.ParseLines(str, mapping); e != nil {
				res = e
			} else {
				res = doc.Value
			}
			if e := compare(res, expect); e != nil {
				t.Fatal("ng:", name, e)
			} else {
				t.Log("ok:", name)
			}
		}
	}
}
