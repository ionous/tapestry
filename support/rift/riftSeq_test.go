package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

func TestSeq(t *testing.T) {
	testSeq(t,
		// --------------
		"a single value", `
- 5`,
		[]any{5.0},

		// --------------
		"a single with newline", `
- 5
`, []any{5.0},

		// --------------
		"split line", `
- 
    5
`, []any{5.0},

		// --------------
		"several values", `
- 5
- 10
- 12`,
		[]any{5.0, 10.0, 12.0},

		// --------------
		"nested sub sequence", `
- - 5`,
		[]any{[]any{5.0}},

		// --------------
		"new line sub sequence", `
-
  - 5
`, []any{[]any{5.0}},

		// --------------
		"nil values", `
-
-
-`,
		[]any{nil, nil, nil},

		// --------------
		"nil value trailing newline", `
-
-
-
`,
		[]any{nil, nil, nil},

		// --------------
		"continuing sub sequence ", `
-
  - 5
- 6`,
		[]any{[]any{5.0}, 6.0})
}

// name of test, input string, expected result
// leading whitespace is trimmed
func testSeq(t *testing.T, nameInputExpect ...any) {
	for i, cnt := 0, len(nameInputExpect); i < cnt; i += 3 {
		name, input, expect := nameInputExpect[0+i].(string), nameInputExpect[1+i].(string), nameInputExpect[2+i]
		if strings.HasPrefix(name, `x `) {
			// commenting out tests causes go fmt to replace spaces with tabs. *sigh*
			t.Log("skipping", name)
		} else {
			var res any
			var doc rift.Document
			str := strings.TrimLeftFunc(input, unicode.IsSpace)
			if e := doc.Parse(str, rift.NewSequence(&doc, 0, func(vs []any) (_ error) {
				res = vs
				return
			})); e != nil {
				res = e
			}
			if e := compare(res, expect); e != nil {
				t.Fatal("ng:", name, e)
			} else {
				t.Log("ok:", name)
			}

		}
	}
}
