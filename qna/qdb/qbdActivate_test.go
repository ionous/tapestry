package qdb_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// write low level activation of domains
// ( not deletion of nouns, or setup of relative pairs --
// just hierarchical selection and detection of changes )
func TestActivate(t *testing.T) {
	// db := testdb.Open(t.Name(), "", "")
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	const at = ""
	if e := createTable(db,
		func(w mdl.Writer) (err error) {
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
	} else if q, e := qdb.NewQueryTest(db); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "main", "sub", "boop", "beep"); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "boop", 1,
		"main:1", "sub:1", "boop:1",
	); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, true, "main", "sub", "boop"); e != nil {
		t.Fatal(e)
	} else if e := isActive(q, false, "beep"); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "beep", 2,
		"beep:2",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "main", 3,
		"main:3",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "sub", 4,
		"main:3",
		"sub:4",
	); e != nil {
		t.Fatal(e)
	} else if e := activate(q, "sub", 5,
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

func isActive(q *qdb.QueryTest, want bool, names ...string) (err error) {
	for _, n := range names {
		if ok, e := q.IsDomainActive(n); e != nil || ok != want {
			err = errutil.New("expected", n, "active", want, e)
			break
		}
	}
	return
}

func activate(q *qdb.QueryTest, name string, act int, expect ...string) (err error) {
	if els, e := q.InnerActivate(name, act); e != nil {
		err = errutil.New("couldnt activate", name, e)
	} else if diff := pretty.Diff(els, expect); len(diff) > 0 {
		err = errutil.New("diff", name, els, diff)
	}
	return
}
