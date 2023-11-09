package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"git.sr.ht/~ionous/tapestry/support/rift/imap"
	"github.com/ionous/errutil"
)

func TestSeq(t *testing.T) {
	testSeq(t,
		// --------------
		"Test single value", `
- 5`,
		[]any{5.0},

		// --------------
		"Test fail without dash", `
-false`,
		errutil.New("unexpected character 'f'"),

		// --------------
		"Test value with newline", `
- 5
`, []any{5.0},

		// --------------
		"Test split line", `
-
  5
`, []any{5.0},

		// --------------
		"Test several values", `
- 5
- 10
- 12`,
		[]any{5.0, 10.0, 12.0},

		// --------------
		"Test nested sub sequence", `
- - 5`,
		[]any{[]any{5.0}},

		// --------------
		"Test new line sub sequence", `
-
  - 5
`, []any{[]any{5.0}},
		// --------------
		"Test multiple sub values", `
- -
  - 5
`, []any{[]any{nil, 5.0}},

		// --------------
		"Test nil values", `
-
-
-`,
		[]any{nil, nil, nil},

		// --------------
		"Test nil value trailing newline", `
-
-
-
`,
		[]any{nil, nil, nil},

		// --------------
		"Test continuing sub sequence ", `
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
			doc := rift.Document{MakeMake: imap.MakeBuilder}
			str := strings.TrimLeftFunc(input, unicode.IsSpace)
			seq := rift.NewSequence(&doc, "", doc.Col)
			if e := doc.ParseLines(str, rift.StartSequence(seq)); e != nil {
				res = e
			} else if val, e := seq.FinalizeValue(); e != nil {
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
