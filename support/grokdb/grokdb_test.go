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
		x := dbSource{domain: "a", db: tables.NewCache(db)}
		groktest.Phrases(t, &x)
	}
}

func TestTraits(t *testing.T) {
	if db, e := setupDB(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		x := dbSource{domain: "a", db: tables.NewCache(db)}
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
		patterns
		macros
		kinds
		where
		carrying
		suspicion
		traits
		things
		containers
		supporters
		domain = "a"
	)
	db := testdb.Create(name)
	if e := tables.CreateModel(db); e != nil {
		err = e
	} else if e := tables.CreateRun(db); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_domain", "domain"},
		domain); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_kind", "ROWID", "domain", "kind", "singular", "path"},
		aspects, domain, kindsOf.Aspect.String(), "", idPath(),
		patterns, domain, kindsOf.Pattern.String(), "", idPath(),
		macros, domain, kindsOf.Macro.String(), "", idPath(patterns),
		kinds, domain, kindsOf.Kind.String(), "", idPath(),
		where, domain, "whereabouts", "", idPath(macros),
		carrying, domain, "carrying", "", idPath(macros),

		suspicion, domain, "suspicion", "", idPath(macros),
		traits, domain, "traits", "", idPath(aspects),
		things, domain, "things", "thing", idPath(kinds),
		containers, domain, "containers", "container", idPath(things, kinds),
		supporters, domain, "supporters", "supporter", idPath(things, kinds),
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field", "domain", "kind", "field", "affinity"},
		domain, traits, "closed", "bool",
		domain, traits, "open", "bool",
		domain, traits, "openable", "bool",
		domain, traits, "transparent", "bool",
		domain, traits, "fixed_in_place", "bool",
		// carrying: one-to-many
		domain, carrying, "actor", "text",
		domain, carrying, "carries", "text_list",
		domain, carrying, "error", "text",
		// whereabouts: one-to-many
		domain, where, "kind", "text",
		domain, where, "other_kinds", "text_list",
		domain, where, "error", "text",
		// suspicion: many-to-many
		domain, suspicion, "kinds", "text_list",
		domain, suspicion, "other_kinds", "text_list",
		domain, suspicion, "error", "text",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_phrase", "domain", "kind", "phrase", "reversed"},
		domain, kinds, "kind of", true, // ex. "a closed kind of container"
		domain, kinds, "kinds of", true, // ex. "are closed containers"
		domain, kinds, "a kind of", true, // ex. "a kind of container"
		//
		domain, carrying, "carrying", false,
		domain, where, "on", false, // on the x are the w,y,z
		domain, where, "in", false,
		domain, suspicion, "suspicious of", false,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
