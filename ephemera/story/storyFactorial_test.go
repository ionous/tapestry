package story_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/story"

	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

// read the factorial test story
func TestFactorialStory(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else {
		k := story.NewImporter(db)
		if e := k.ImportStory(t.Name(), debug.FactorialStory); e != nil {
			t.Fatal(e)
		} else {
			var buf strings.Builder
			tables.WriteCsv(db, &buf, "select count() from eph_check", 1)
			tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
			tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
			tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
			tables.WriteCsv(db, &buf, "select distinct name, category from eph_named where category != 'scene' order by name, category", 2)
			if have, want := buf.String(), lines(
				"1", // eph_check -- 1 unit test
				"2", // eph_rule -- 2 rules
				"3", // eph_prog -- 1 test program, 2 rules
				"4", // eph_pattern specifies types - (1 pattern, 1 parameter) * (1 decl, 1 call)
				// eph_named
				"factorial,pattern", // name of the pattern
				"factorial,test",    // name of the test
				"num,parameter",     // we declared the param
				"num,return",        // we referenced the return
				"number_eval,type",  // we evaluated the var
				"patterns,type",     // we evaluated the return
			); have != want {
				t.Fatal(have)
			}
		}
	}
}
