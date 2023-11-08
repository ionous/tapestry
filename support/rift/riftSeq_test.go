package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
)

func xTestSeq(t *testing.T) {
	testSeq(t,
		// --------------
		"test single value", `
- 5`,
		[]any{5.0},

		// --------------
		"fail without dash", `
-false`,
		errutil.New(""),

		// --------------
		"test value with newline", `
- 5
`, []any{5.0},

		// --------------
		"test split line", `
-
  5
`, []any{5.0},

		// --------------
		"test several values", `
- 5
- 10
- 12`,
		[]any{5.0, 10.0, 12.0},

		// --------------
		"test nested sub sequence", `
- - 5`,
		[]any{[]any{5.0}},

		// --------------
		"test new line sub sequence", `
-
  - 5
`, []any{[]any{5.0}},

		// --------------
		"test nil values", `
-
-
-`,
		[]any{nil, nil, nil},

		// --------------
		"test nil value trailing newline", `
-
-
-
`,
		[]any{nil, nil, nil},

		// --------------
		"test continuing sub sequence ", `
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
			seq := rift.NewSequence(&doc, "", doc.Col)
			if e := doc.ParseLines(str, seq); e != nil {
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
