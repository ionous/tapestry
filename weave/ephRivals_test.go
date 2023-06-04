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
	} else if _, e := db.Exec(`insert into mdl_domain(domain, requires) values
	('p1', ''),
	('p2', '')`); e != nil {
		t.Fatal(e)
	} else if _, e := db.Exec(`insert into mdl_plural(domain, many, one) values
	('p1', 'people', 'human'),
	('p1', 'bats', 'bat'),
	('p2', 'people', 'person'),
	('p2', 'rabbits', 'rabbit')`); e != nil {
		t.Fatal(e)
	} else if _, e := db.Exec(`insert into run_domain(domain, active) values
	('p1', 1),
	('p2', 1)`); e != nil {
		t.Fatal(e)
	} else if conflicts, e := findRivals(db); e != nil {
		t.Fatal(e)
	} else {
		expect := []conflict{
			{Category: "plural", Domain: "p1", Key: "people", Value: "human"},
			{Category: "plural", Domain: "p2", Key: "people", Value: "person"},
		}
		if diff := pretty.Diff(expect, conflicts); len(diff) > 0 {
			t.Log("got", pretty.Sprint(conflicts))
			t.Fatal("unexpected conflicts", diff)
		}
	}
}