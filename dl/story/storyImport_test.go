package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"git.sr.ht/~ionous/tapestry/weave"
)

// decode and import a story;
// only checks that the process finishes; doesnt check the results.
func TestImportStory(t *testing.T) {
	var msg map[string]any
	var curr story.StoryFile
	if e := json.Unmarshal(debug.Blob, &msg); e != nil {
		t.Fatal("couldn't unmarshal blob", e)
	} else if e := story.Decode(&curr, msg); e != nil {
		t.Fatal("couldn't decode story", e)
	} else {
		db := tables.CreateTest(t.Name(), true)
		defer db.Close()
		cat := weave.NewCatalog(db)
		//
		d := cat.EnsureScene("tapestry")
		if _, e := cat.SceneBegin(d, compact.Source{}, nil); e != nil {
			t.Fatal(e)
		} else {
			defer cat.SceneEnd()
			if e := curr.Weave(cat); e != nil {
				t.Fatal("failed story import", e)
			} else {
				t.Log("ok")
			}
		}
	}
}
