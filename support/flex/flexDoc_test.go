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
	if e := flex.ReadStory("beep", strings.NewReader(testDoc), &out); e != nil {
		t.Fatal(e)
	} else {
		// fix? doesnt compare against a particular result
		// because some things have line numbers, and some dont
		// probably they arent event correct yet
		pretty.Println(out)
	}
}
