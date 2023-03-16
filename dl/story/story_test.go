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
		k := imp.NewImporter(collectEphemera(&els))
		if e := story.ImportStory(k, t.Name(), &curr); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}

func collectEphemera(sink *[]eph.Ephemera) imp.WriterFun {
	return func(el eph.Ephemera) {
		*sink = append(*sink, el)
	}
}
