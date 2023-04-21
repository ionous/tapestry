package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/test/debug"
)

func TestImportStory(t *testing.T) {
	var curr story.StoryFile
	if e := story.Decode(&curr, debug.Blob, tapestry.AllSignatures); e != nil {
		t.Fatal(e)
	} else {
		var els []eph.Ephemera
		k := imp.NewImporter(eph.NewCommandQueue(&els))
		if e := story.ImportStory(k, t.Name(), &curr); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}
