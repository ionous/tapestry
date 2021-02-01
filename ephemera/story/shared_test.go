package story_test

import (
	"database/sql"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

// for tests where we need a default decoder to read json
func newImporter(t *testing.T, where string) (ret *story.Importer, retDec *decode.Decoder, retDB *sql.DB) {
	db := newImportDB(t, where)
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	} else {
		iffy.RegisterGobs()
		dec := decode.NewDecoderReporter(t.Name(), func(pos reader.Position, err error) {
			t.Errorf("%s at %s", err, pos)
		})
		k := story.NewImporterDecoder(t.Name(), db, dec)
		dec.AddDefaultCallbacks(core.Slats)
		k.AddModel(story.Model)
		ret, retDec, retDB = k, dec, db
	}
	return
}

// if path is nil, it will use a file db.
func newImportDB(t *testing.T, where string) (ret *sql.DB) {
	var source string
	if len(where) > 0 {
		source = where
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		source = p
	}
	//
	if db, e := sql.Open(tables.DefaultDriver, source); e != nil {
		t.Fatal(e)
	} else {
		t.Log("opened db", source)
		ret = db
	}
	return
}
