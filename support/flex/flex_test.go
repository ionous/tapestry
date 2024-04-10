package flex_test

import (
	"errors"
	"io"
	r "reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/flex"
)

func TestComments(t *testing.T) {
	test := func(name, src string, expect any) {
		t.Log(name)
		if lines, e := flex.ReadComments(strings.NewReader(src)); matchError(t, e, expect) {
			want := expect.([]string)
			if !r.DeepEqual(lines, want) {
				t.Logf("expected %#v\n", want)
				t.Logf("got %#v\n", lines)
				t.Fatal("error: mismatched.")
			}
		}
	}
	test("empty", "", []string(nil))
	test("bad text", "adada", errors.New("header"))
	test("inline eof", `
# hello`, []string{"hello"})
	test("newline eof", `


# hello
`, []string{"hello"})
}

// returns true if expected success and got success
func matchError(t *testing.T, got error, want any) (okay bool) {
	if want, wantedError := want.(error); got == nil {
		if okay = !wantedError; !okay {
			t.Fatalf("expected an error %q", got)
		}
	} else {
		if !wantedError {
			t.Fatalf("expected success, but got %q", got)
		} else {
			gotStr, wantStr := got.Error(), want.Error()
			if !strings.HasPrefix(gotStr, wantStr) {
				t.Fatalf("expected %q but got %q", wantStr, gotStr)
			}
		}
	}
	return
}

// count the number of sections in some sample darta.
func TestSections(t *testing.T) {
	expect := []int{0, 3, 9, 13}
	var got []int
	// break the text into lines
	in := strings.NewReader(testDoc)
	for k := flex.MakeSection(in); k.NextSection(); {
		got = append(got, k.StartingLine)
		for {
			if _, _, e := k.ReadRune(); e != nil {
				if e != io.EOF {
					t.Fatalf("fatal error %v after %#v", e, expect)
				}
				break
			}
		}
	}
	if !r.DeepEqual(expect, got) {
		t.Logf("expected %#v\n", expect)
		t.Logf("got %#v\n", got)
		t.Fatal("mismatch")
	}
}

var testDoc = `
# Header:
---
# First Structured:
Define scene:requires:with:
- "cloak"
- "tapestry"
- Declare: ""
---
# Plain:
The title of the story is "The Cloak of Darkness."
The story has the headline "An example story."
---
# Second Structured:
Define rule:noun:do:
- "instead of traveling"
- "entrance"
- Say: "You've only just arrived and besides the weather outside is terrible."
---
`
