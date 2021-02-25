package qna

import (
	"testing"

	"git.sr.ht/~ionous/iffy/assembly"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

// no idea where this test should live...
// tests the importation, assembly, and execution of a factorial story.
// doesn't test the *reading* of the story.
func TestFullFactorial(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, assembly.SqlCustomDriver)
	defer db.Close()

	// read factorialStory, assemble and run.
	var ds reader.Dilemmas
	if e := tables.CreateAll(db); e != nil {
		t.Fatal("couldn't create tables", e)
	} else {
		k := story.NewImporterDecoder(db, decode.NewDecoder()).SetSource(t.Name())
		//
		if e := debug.FactorialStory.ImportStory(k); e != nil {
			t.Fatal("couldn't import story", e)
		} else if e := assembly.AssembleStory(db, "kinds", ds.Add); e != nil {
			t.Fatal("couldnt assemble story", e, ds.Err())
		} else if len(ds) > 0 {
			t.Fatal("issues assembling", ds.Err())
		} else if cnt, e := CheckAll(db, ""); e != nil {
			t.Fatal(e)
		} else if cnt != 1 {
			t.Fatal("expected one test", cnt)
		} else {
			t.Log("ok", cnt)
		}
	}
}
