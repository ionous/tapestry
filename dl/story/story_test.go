package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave"
)

// decode and import a story;
// only checks that the process finishes; doesnt check the results.
func TestImportStory(t *testing.T) {
	if curr, e := story.CompactDecode(debug.Blob); e != nil {
		t.Fatal(e)
	} else {
		db := testdb.Create(t.Name())
		defer db.Close()
		k := weave.NewCatalog(db)
		if e := k.DomainStart("tapestry", nil); e != nil {
			t.Fatal("import", e)
		} else if e := story.ImportStory(k, t.Name(), &curr); e != nil {
			t.Fatal("import", e)
		} else if e := k.DomainEnd(); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}
