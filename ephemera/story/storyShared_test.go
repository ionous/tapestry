package story_test

import (
	"database/sql"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

// for tests where we need a default decoder to read json
func newImporter(t *testing.T, where string) (*story.Importer, *sql.DB) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	}
	k := story.NewImporter(db, cout.Marshal)
	k.SetSource(t.Name())
	return k, db
}
