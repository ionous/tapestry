package flex_test

import (
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/flex"
	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/kr/pretty"
)

// tokens tests the individual elements
// this test the output commands
// fails because of positions in match output.
func xTestPlainText(t *testing.T) {
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
			Matches: story.JessMatches(panicTokens(
				// fix: currently can't handle "It's"
				"It is a good day for plain text.",
			)),
		},
		&story.Comment{
			Lines: []string{"And comments.", "I guess."},
		},
		&story.DeclareStatement{
			Text: &literal.TextValue{
				Value: "No. Really.",
			},
			Matches: story.JessMatches(panicTokens(
				"No. Really",
			)),
		},
	)
}

func panicTokens(str string) [][]match.TokenValue {
	c := match.Collector{BreakLines: true, KeepComments: true}
	if e := c.Collect(str); e != nil {
		panic(e)
	} else {
		return c.Lines
	}
}

func testText(t *testing.T, in string, expect ...story.StoryStatement) {
	runes := strings.NewReader(in)
	if out, e := flex.ReadText(runes); e != nil {
		t.Fatal(e)
	} else if !reflect.DeepEqual(out, expect) {
		t.Log(pretty.Sprint(out))
		t.Fatal("mismatch")
	}
}
