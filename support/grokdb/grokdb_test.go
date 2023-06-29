package grokdb

import (
	"database/sql"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/groktest"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
)

func TestPhrases(t *testing.T) {
	if db, e := setupDB(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		x := dbg{domain: "a", db: db}
		groktest.Phrases(t, &x)
	}
}

func TestTraits(t *testing.T) {
	if db, e := setupDB(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		x := dbg{domain: "a", db: db}
		groktest.Traits(t, &x)
	}
}

func setupDB(name string) (ret *sql.DB, err error) {
	db := testdb.Create(name)
	if e := tables.CreateModel(db); e != nil {
		err = e
	} else if e := tables.CreateRun(db); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_domain", "domain"},
		"a"); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_kind", "ROWID", "domain", "kind", "singular", "path"},
		1, "a", kindsOf.Aspect.String(), "", "",
		2, "a", kindsOf.Relation.String(), "", "",
		5, "a", kindsOf.Kind.String(), "", "",
		10, "a", "whereabouts", "", "2,",
		11, "a", "suspicion", "", "2,",
		20, "a", "traits", "", "1,",
		88, "a", "things", "thing", "5,",
		89, "a", "containers", "container", "88,5,",
		90, "a", "supporters", "supporter", "88,5,",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field", "domain", "kind", "field", "affinity"},
		"a", 20, "closed", "bool",
		"a", 20, "open", "bool",
		"a", 20, "openable", "bool",
		"a", 20, "transparent", "bool",
		"a", 20, "fixed_in_place", "bool",
		// whereabouts: one-to-many
		"a", 10, "kind", "text",
		"a", 10, "other_kinds", "text_list",
		// suspicion: many-to-many
		"a", 11, "kinds", "text_list",
		"a", 11, "other_kinds", "text_list",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_grok", "domain", "kind", "phrase", "reversed"},
		"a", 10, "kind of", true, // for "a closed kind of container"
		"a", 10, "kinds of", true, // for "are closed containers"
		"a", 10, "a kind of", true, // for "a kind of container"
		"a", 10, "on", false, // on the x are the w,y,z
		"a", 10, "in", false,
		"a", 11, "suspicious of", false,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
