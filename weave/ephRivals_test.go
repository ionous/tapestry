package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"github.com/kr/pretty"
)

// using the plurals as an example,
// verify the basic algorithm behind rival testing.
func TestRivals(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	if e := tables.CreateAll(db); e != nil {
		t.Fatal(e)
	} else if _, e := db.Exec(`insert into mdl_domain(rowid, domain, path) values
	(1, 'p1', ''),
	(2, 'p2', '')`); e != nil {
		t.Fatal(e)
	} else if _, e := db.Exec(`insert into mdl_plural(domain, many, one) values
	(1, 'people', 'human'),
	(1, 'bats', 'bat'),
	(2, 'people', 'person'),
	(2, 'rabbits', 'rabbit')`); e != nil {
		t.Fatal(e)
	} else if _, e := db.Exec(`insert into run_domain(domain, active) values
	(1, 1),
	(2, 1)`); e != nil {
		t.Fatal(e)
	} else if conflicts, e := findConflicts(db); e != nil {
		t.Fatal(e)
	} else {
		expect := []conflict{
			{Category: "plural", Domain: "p1", Key: "people", Value: "human"},
			{Category: "plural", Domain: "p2", Key: "people", Value: "person"},
		}
		if diff := pretty.Diff(expect, conflicts); len(diff) > 0 {
			t.Fatal("unexpected conflicts", diff)
		}
	}
}
