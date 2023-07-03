package grokdb

import (
	"database/sql"
	"strconv"
	"strings"
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

func idPath(ids ...int) string {
	var b strings.Builder
	for _, el := range ids {
		b.WriteString(strconv.Itoa(el))
		b.WriteRune(',')
	}
	return b.String()
}

func setupDB(name string) (ret *sql.DB, err error) {
	const (
		aspects = iota + 1
		rels
		kinds
		where
		suspicion
		traits
		things
		containers
		supporters
	)
	db := testdb.Create(name)
	if e := tables.CreateModel(db); e != nil {
		err = e
	} else if e := tables.CreateRun(db); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_domain", "domain"},
		"a"); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_kind", "ROWID", "domain", "kind", "singular", "path"},
		aspects, "a", kindsOf.Aspect.String(), "", idPath(),
		rels, "a", kindsOf.Relation.String(), "", idPath(),
		kinds, "a", kindsOf.Kind.String(), "", idPath(),
		where, "a", "whereabouts", "", idPath(rels),
		suspicion, "a", "suspicion", "", idPath(rels),
		traits, "a", "traits", "", idPath(aspects),
		things, "a", "things", "thing", idPath(kinds),
		containers, "a", "containers", "container", idPath(things, kinds),
		supporters, "a", "supporters", "supporter", idPath(things, kinds),
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field", "domain", "kind", "field", "affinity"},
		"a", traits, "closed", "bool",
		"a", traits, "open", "bool",
		"a", traits, "openable", "bool",
		"a", traits, "transparent", "bool",
		"a", traits, "fixed_in_place", "bool",
		// whereabouts: one-to-many
		"a", where, "kind", "text",
		"a", where, "other_kinds", "text_list",
		// suspicion: many-to-many
		"a", suspicion, "kinds", "text_list",
		"a", suspicion, "other_kinds", "text_list",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_phrase", "domain", "kind", "phrase", "reversed"},
		"a", kinds, "kind of", true, // ex. "a closed kind of container"
		"a", kinds, "kinds of", true, // ex. "are closed containers"
		"a", kinds, "a kind of", true, // ex. "a kind of container"
		//
		"a", where, "on", false, // on the x are the w,y,z
		"a", where, "in", false,
		"a", suspicion, "suspicious of", false,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
