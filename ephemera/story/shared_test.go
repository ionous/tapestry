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
	db := testdb.Open(t.Name(), testdb.Memory, "")
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	}
	iffy.RegisterGobs()
	dec := decode.NewDecoderReporter(func(pos reader.Position, err error) {
		t.Errorf("%s at %s", err, pos)
	})
	k := story.NewImporterDecoder(db, dec).SetSource(t.Name())
	dec.AddDefaultCallbacks(core.Slats)
	k.AddModel(story.Model)
	ret, retDec, retDB = k, dec, db
	return
}
