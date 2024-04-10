package flex_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/flex"
	"github.com/kr/pretty"
)

// FIX
func xTestDoc(t *testing.T) {
	var out story.StoryFile
	if e := flex.ReadStory("beep", strings.NewReader(testDoc), &out); e != nil {
		t.Fatal(e)
	} else {
		pretty.Println(out)
	}
}
