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
		// built in kinds
		aspects = iota + 1
		patterns
		macros
		kinds
		traits
		// kinds of nouns
		things
		containers
		supporters
		// macros
		carry
		contain
		inherit
		support
		suspect
		// domain string
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
	} else if e := testdb.Ins(db, []string{"mdl_kind",
		"ROWID", "domain", "kind", "singular", "path"},
		aspects, domain, kindsOf.Aspect.String(), "", idPath(),
		patterns, domain, kindsOf.Pattern.String(), "", idPath(),
		macros, domain, kindsOf.Macro.String(), "", idPath(patterns),
		kinds, domain, kindsOf.Kind.String(), "", idPath(),
		traits, domain, "traits", "", idPath(aspects),
		things, domain, "things", "thing", idPath(kinds),
		containers, domain, "containers", "container", idPath(things, kinds),
		supporters, domain, "supporters", "supporter", idPath(things, kinds),
		// macros:
		carry, domain, "carry", "", idPath(macros),
		contain, domain, "contain", "", idPath(macros),
		inherit, domain, "inherit", "", idPath(macros),
		support, domain, "support", "", idPath(macros),
		suspect, domain, "suspect", "", idPath(macros),
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field",
		"domain", "kind", "field", "affinity"},
		//
		domain, traits, "closed", "bool",
		domain, traits, "open", "bool",
		domain, traits, "openable", "bool",
		domain, traits, "transparent", "bool",
		domain, traits, "fixed_in_place", "bool",
		// macros:
		domain, carry, "source", "text",
		domain, carry, "targets", "text_list",
		domain, carry, "error", "text",
		//
		domain, contain, "source", "text",
		domain, contain, "targets", "text_list",
		domain, contain, "error", "text",
		//
		domain, inherit, "sources", "text_list",
		domain, inherit, "error", "text",
		//
		domain, support, "source", "text",
		domain, support, "targets", "text_list",
		domain, support, "error", "text",
		// suspicion: many-to-many
		domain, suspect, "sources", "text_list",
		domain, suspect, "targets", "text_list",
		domain, suspect, "error", "text",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_pat",
		"kind", "result"},
		//
		carry, "error",
		contain, "error",
		inherit, "error",
		support, "error",
		suspect, "error",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_phrase",
		"domain", "macro", "phrase", "reversed"},
		//
		domain, carry, "carried by", true, // ex. source carrying targets
		domain, carry, "carrying", false, // ex. source carrying targets
		domain, contain, "in", true,
		domain, inherit, "kind of", false, // ex. "a closed kind of container"
		domain, inherit, "kinds of", false, // ex. "are closed containers"
		domain, inherit, "a kind of", false, // ex. "a kind of container"
		domain, support, "on", true, // on the x are the w,y,z
		domain, suspect, "suspicious of", false,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
