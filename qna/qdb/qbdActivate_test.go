package qdb_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// test hierarchical selection of domains and detection of changes
func TestActivate(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	const at = ""
	if e := createTable(db,
		func(w mdl.Modeler) (err error) {
			return write(w,
				// name, dependency, at
				// -------------------------
				mdl.Domain,
				"main", "", at,
				"sub", "main", at,
				"beep", "", at,
				"boop", "sub", at,
			)
		}); e != nil {
		t.Fatal("failed to create table", e)
	} else if q, e := qdb.NewQueryTest(db); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "main", "sub", "boop", "beep"); e != nil {
		t.Fatal(e) // ^ verify nothing is active by default.
	} else if e := activate(q, "boop",
		"main:1", // active domains and the value of their activation count
		"sub:1",  // ( in order of most root to leaf )
		"boop:1", // ( ie. boop depends on sub, which depends on main )
	); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, true, "main", "sub", "boop"); e != nil {
		t.Fatal(e) // ^ verify using the actual IsDomainActive call
	} else if e := isActive(q, false, "beep"); e != nil {
		t.Fatal(e) // ^ verify using the actual IsDomainActive call
	} else if e := activate(q, "beep",
		"beep:2", // only beep is active; it's newly active and has the new activation count
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "main",
		"main:3", // only main is active; it's newly active and has the new activation count
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "sub",
		"main:3", // main is still active, and so its counter is unchanged
		"sub:4",  // sub is newly active, and has the current count
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "sub",
		"main:3", // nothing has changed
		"sub:4",
	); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, true, "main", "sub"); e != nil {
		t.Fatal(e) // ^ verify using the actual IsDomainActive call
	} else if e := isActive(q, false, "boop", "beep"); e != nil {
		t.Fatal(e) // ^ verify using the actual IsDomainActive call
	}
}

func isActive(q *qdb.QueryTest, want bool, names ...string) (err error) {
	for _, n := range names {
		if ok, e := q.IsDomainActive(n); e != nil || ok != want {
			err = errutil.New("expected", n, "active", want, e)
			break
		}
	}
	return
}

func activate(q *qdb.QueryTest, name string, expect ...string) (err error) {
	if els, e := q.InnerActivate(name); e != nil {
		err = errutil.New("couldnt activate", name, e)
	} else if diff := pretty.Diff(els, expect); len(diff) > 0 {
		err = errutil.New("diff", name, els, diff)
	}
	return
}
