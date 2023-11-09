package rift_test

import (
	"strings"
	"testing"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"github.com/ionous/errutil"
)

func TestScalars(t *testing.T) {
	testValue(t,
		"Test number", `5.4`, 5.4,
		"Test string", `"5.4"`, "5.4",
		"Test bool", `true`, true,

		// ----------
		"Test interpreted",
		`"hello\\world"`,
		`hello\world`,

		// ----------
		"Test raw",
		"`"+`hello\\world`+"`",
		`hello\\world`,

		// -----
		"Test unquoted value",
		"beep",
		errutil.New("signature must end with a colon"),
	)
}

//  name of test, input string, expected result
// leading whitespace is trimmed
func testValue(t *testing.T, nameInputExpect ...any) {
	for i, cnt := 0, len(nameInputExpect); i < cnt; i += 3 {
		name, input, expect := nameInputExpect[0+i].(string), nameInputExpect[1+i].(string), nameInputExpect[2+i]
		if strings.HasPrefix(name, `x `) {
			// commenting out tests causes go fmt to replace spaces with tabs. *sigh*
			t.Log("skipping", name)
		} else {
			var res any
			var doc rift.Document
			str := strings.TrimLeftFunc(input, unicode.IsSpace)
			if e := doc.ParseLines(str, rift.CollectionEntry(&doc, 0)); e != nil {
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
