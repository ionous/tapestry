package story

import (
	"encoding/json"
	"testing"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func TestImportStory(t *testing.T) {
	db := newTestDB(t, memory)
	defer db.Close()
	//
	var in reader.Map
	if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
		t.Fatal("read json", e)
	} else if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else if e := ImportStory(t.Name(), db, in); e != nil {
		t.Fatal("import", e)
	} else {
		t.Log("ok")
	}
}