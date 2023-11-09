package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

// keys that start with t or f need special handleing
func TestMap(t *testing.T) {
	testMap(t,
		// -----------
		"x test keys with boolean names", `
true: false
false: true`,
		rift.MapValues{
			{"true:", false},
			{"false:", true},
		},

		// -----------
		"x test single value", `
name: "Sammy Sosa"`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
		},
		// -----------
		"x test split line", `
name:
  "Sammy Sosa"`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
		},

		// -----------
		"x test several values", `
name: "Sammy Sosa"
hr:   63
avg:  true`,
		rift.MapValues{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		},
		// -----------------------
		"test map with nil value", `
Field:
Next: 5`,
		[]any{
			rift.MapValues{
				{"Field:", nil},
				{"Next:", 5.0},
			}},

		// -----------------------
		"x test nested maps", `
Field:
  Next: 5`,
		[]any{
			rift.MapValues{
				{"Field:", rift.MapValues{
					{"Next:", 5.0},
				}},
			}},

		// -----------------------
		// in yaml, inline nested maps are invalid
		// should they be here too?
		// to do, i think Value would need to examine history
		// either sniffing prior types or through a flag (ex. require newlines)
		// that it can send into NewMapping
		"x test inline maps", `
Field: Next: 5`,
		[]any{
			rift.MapValues{{
				"Field:", rift.MapValues{{
					"Next:", 5.0,
				}},
			}}},
	)
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
			if e := doc.ParseLines(str, rift.StartMapping(mapping)); e != nil {
				res = e
			} else if val, e := mapping.FinalizeValue(); e != nil {
				res = e // calls finalize directly because the sequence was handled directly to parse,
			} else {
				res = val
			}
			if e := compare(res, expect); e != nil {
				t.Fatal("ng:", name, e)
			} else {
				t.Log("ok:", name)
			}
		}
	}
}
