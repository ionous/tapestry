package flex_test

import (
	"io"
	r "reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/flex"
	"github.com/kr/pretty"
)

// count the number of sections in some sample data.
func TestSectionCount(t *testing.T) {
	expect := []int{0, 2, 8, 12}
	var got []int
	// break the text into lines
	in := strings.NewReader(testOne)
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

// fix? doesnt compare against a particular result
// because some things have line numbers, and some dont
// compare against tell maps ( instead of parsed stories? )
func TestDoc(t *testing.T) {
	if out, e := flex.ReadStory(strings.NewReader(testTwo)); e != nil {
		t.Fatal(e)
	} else {
		pretty.Println(out)
	}
	if out, e := flex.ReadStory(strings.NewReader(testOne)); e != nil {
		t.Fatal(e)
	} else {
		pretty.Println(out)
	}
}

var testOne = `
# First Plain Text:
---
# First Structured:
Define scene:requires:with:
- "cloak"
- "tapestry"
- Declare: ""
---
# Second Plain Text:
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

var testTwo = `---
# Everything in the Tapestry directory
Define scene: "tapestry"`
