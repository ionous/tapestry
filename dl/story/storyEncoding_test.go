package story_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"
)

func TestDecodeStory(t *testing.T) {
	var m map[string]any
	var file story.StoryFile
	if e := json.Unmarshal([]byte(debug.FactorialJs), &m); e != nil {
		t.Fatal(e)
	} else if e := story.Decode(&file, m); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, file); len(diff) != 0 {
		pretty.Print(file)
		t.Fatal(diff)
	}
}
