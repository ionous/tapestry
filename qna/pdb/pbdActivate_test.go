package pdb

import (
	"database/sql"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/tables/mdl"
	"git.sr.ht/~ionous/iffy/test/testdb"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// write test data to the database, and ensure we can query it back.
// this exercises the asm.writer ( xforming from strings to ids )
// and the various runtime queries we need.
func TestActivate(t *testing.T) {
	// db := testdb.Open(t.Name(), "", "")
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	const at = ""
	if e := createTable(db,
		func(w eph.Writer) (err error) {
			return write(w,
				// name, path, at
				// -------------------------
				mdl.Domain,
				"main", "", at,
				"sub", "main", at,
				"beep", "", at,
				"boop", "sub,main", at,
			)
		}); e != nil {
		t.Fatal("failed to create table", e)
	} else if q, e := NewQueries(db); e != nil {
		t.Fatal(e)
	} else if scan, e := db.Prepare(
		`select md.domain || ':' || rd.active
		from run_domain rd
		join mdl_domain md
			on (rd.domain=md.rowid)
		where rd.active > 0
		`); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "main", "sub", "boop", "beep"); e != nil {
		t.Fatal(e)
	} else if e := activate(q, scan, "boop", 1,
		"main:1", "sub:1", "boop:1",
	); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, true, "main", "sub", "boop"); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "beep"); e != nil {
		t.Fatal(e)
	} else if e := activate(q, scan, "beep", 2,
		"beep:2",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, scan, "main", 3,
		"main:3",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, scan, "sub", 4,
		"main:3",
		"sub:4",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, scan, "sub", 5,
		"main:3",
		"sub:4",
	); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, true, "main", "sub"); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "boop", "beep"); e != nil {
		t.Fatal(e)
	}
}

func isActive(q *Query, want bool, names ...string) (err error) {
	for _, n := range names {
		if ok, e := q.IsDomainActive(n); e != nil || ok != want {
			err = errutil.New("expected", n, "active", want, e)
			break
		}
	}
	return
}

func activate(q *Query, scan *sql.Stmt, name string, act int, expect ...string) (err error) {
	if _, e := q.domainActivation.Exec(name, act); e != nil {
		err = errutil.New("couldnt activate", name, e)
	} else if els, e := scanStrings(scan); e != nil {
		err = errutil.New("couldnt scan", name, e)
	} else if diff := pretty.Diff(els, expect); len(diff) > 0 {
		err = errutil.New("diff", name, els, diff)
	}
	return
}
