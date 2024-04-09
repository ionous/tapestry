package flex_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/flex"
	"github.com/kr/pretty"
)

func TestDoc(t *testing.T) {
	var out story.StoryFile
	if e := flex.ReadStory(strings.NewReader(testDoc), &out); e != nil {
		t.Fatal(e)
	} else {
		pretty.Println(out)
	}
}
