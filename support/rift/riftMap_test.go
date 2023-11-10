package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"git.sr.ht/~ionous/tapestry/support/rift/imap"
)

// keys that start with t or f need special handleing
func TestMap(t *testing.T) {
	testMap(t,
		// -----------
		"Test keys with boolean names", `
true: false
false: true`,
		imap.ItemMap{
			{"true:", false},
			{"false:", true},
		},

		// -----------
		"Test single value", `
name: "Sammy Sosa"`,
		imap.ItemMap{
			{"name:", "Sammy Sosa"},
		},
		// -----------
		"Test split line", `
name:
  "Sammy Sosa"`,
		imap.ItemMap{
			{"name:", "Sammy Sosa"},
		},

		// -----------
		"Test several values", `
name: "Sammy Sosa"
hr:   63
avg:  true`,
		imap.ItemMap{
			{"name:", "Sammy Sosa"},
			{"hr:", 63.0},
			{"avg:", true},
		},
		// -----------------------
		"Test map with nil value", `
Field:
Next: 5`,
		imap.ItemMap{
			{"Field:", nil},
			{"Next:", 5.0},
		},

		// -----------------------
		"Test nested maps", `
Field:
  Next: 5`,
		imap.ItemMap{
			{"Field:", imap.ItemMap{
				{"Next:", 5.0},
			}},
		},

		// -----------------------
		// in yaml, inline nested maps are invalid
		// should they be here too?
		// to do, i think Value would need to examine history
		// either sniffing prior types or through a flag (ex. require newlines)
		// that it can send into NewMapping
		"Test inline maps", `
Field: Next: 5`,
		imap.ItemMap{{
			"Field:", imap.ItemMap{{
				"Next:", 5.0,
			}},
		}},
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
			doc := rift.Document{MakeMap: imap.Build, CommentBlock: rift.DiscardCommentWriter()}
			str := strings.TrimLeftFunc(input, unicode.IsSpace)
			mapping := rift.NewMapping(&doc, "", 0)
			if e := doc.ReadLines(strings.NewReader(str), rift.StartMapping(mapping)); e != nil {
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
