package flex_test

import (
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/flex"
	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
	"github.com/kr/pretty"
)

// tokens tests the individual elements
// this test the output commands
func TestPlainText(t *testing.T) {
	testText(t, `
It is a good day for plain text.
# And comments.
# I guess.
No. Really.
`,
		&story.DeclareStatement{
			Text: &literal.TextValue{
				Value: "It is a good day for plain text.",
			},
			Matches: story.JessMatches(match.PanicSpans(
				// FIX: currently can't handle It's
				"It is a good day for plain text", // no terminal
			)),
		},
		&story.Comment{
			Lines: []string{"And comments.", "I guess."},
		},
		&story.DeclareStatement{
			Text: &literal.TextValue{
				Value: "No. Really.",
			},
			Matches: story.JessMatches(match.PanicSpans(
				"No",
				"Really",
			)),
		},
	)
}

func testText(t *testing.T, in string, expect ...story.StoryStatement) {
	var pt flex.PlainText
	runes := strings.NewReader(in)
	run := flex.NewTokenizer(&pt)
	if e := charm.Read(runes, run); e != nil {
		t.Fatal(e)
	} else {
		out := pt.Finalize()
		if !reflect.DeepEqual(out, expect) {
			t.Log(pretty.Sprint(out))
			t.Fatal("mismatch")
		}

	}
	return
}
