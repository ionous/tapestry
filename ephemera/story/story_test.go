package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func TestImportStory(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	//
	var in reader.Map
	if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
		t.Fatal("read json", e)
	} else if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else {
		k := story.NewImporter(db, func(pos reader.Position, err error) {
			t.Errorf("%s at %s", err, pos)
		})
		if _, e := k.ImportStory(t.Name(), in); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}
