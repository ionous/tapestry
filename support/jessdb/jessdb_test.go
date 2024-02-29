package jessdb_test

import (
	"database/sql"
	"strconv"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/support/jesstest"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func TestPhrases(t *testing.T) {
	if db, e := setupDB(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if m, e := mdl.NewModeler(db); e != nil {
			t.Fatal(e)
		} else {
			m.PrecachePaths()
			x := jessdb.NewSource(m)
			jesstest.RunPhraseTests(t, func(testPhrase string) (ret jess.Generator, err error) {
				if ws, e := match.MakeSpan(testPhrase); e != nil {
					err = e
				} else {
					ret, err = x.MatchSpan("a", ws)
				}
				return
			})
		}
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
		records
		actions
		// kinds of nouns
		things
		containers
		supporters
		colors
		groups
		// actions
		storing
		// macros
		carry
		contain
		support
		suspect
		// nouns
		message
		missive
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
		aspects, domain, "aspects", "aspect", idPath(),
		patterns, domain, "patterns", "pattern", idPath(),
		records, domain, "records", "record", idPath(),
		macros, domain, "macros", "macro", idPath(patterns),
		actions, domain, "actions", "action", idPath(patterns),
		kinds, domain, "kinds", "kind", idPath(),
		traits, domain, "traits", "", idPath(aspects),
		things, domain, "things", "thing", idPath(kinds),
		containers, domain, "containers", "container", idPath(things, kinds),
		supporters, domain, "supporters", "supporter", idPath(things, kinds),
		colors, domain, "color", nil, idPath(aspects),
		groups, domain, "groups", "group", idPath(records),
		// macros:
		carry, domain, "carry", nil, idPath(macros),
		contain, domain, "contain", nil, idPath(macros),
		support, domain, "support", nil, idPath(macros),
		suspect, domain, "suspect", nil, idPath(macros),
		// actions
		storing, domain, "storing", nil, idPath(actions),
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field",
		"domain", "kind", "field", "affinity"},
		// traits
		domain, traits, "closed", "bool",
		domain, traits, "open", "bool",
		domain, traits, "openable", "bool",
		domain, traits, "transparent", "bool",
		domain, traits, "fixed in place", "bool",
		// fields
		domain, things, "description", "text",
		domain, things, "title", "text",
		domain, things, "age", "text",
		// macros:
		domain, carry, "primary", "text",
		domain, carry, "secondary", "text_list",
		domain, carry, "error", "text",
		//
		domain, contain, "primary", "text",
		domain, contain, "secondary", "text_list",
		domain, contain, "error", "text",
		//
		domain, support, "primary", "text",
		domain, support, "secondary", "text_list",
		domain, support, "error", "text",
		// suspicion: many-to-many
		domain, suspect, "primary", "text_list",
		domain, suspect, "secondary", "text_list",
		domain, suspect, "error", "text",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_pat",
		"kind", "result"},
		//
		carry, "error",
		contain, "error",
		support, "error",
		suspect, "error",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_phrase",
		"domain", "macro", "phrase", "reversed"},
		//
		domain, carry, "carried by", true, // ex. primary carrying secondary
		domain, carry, "carrying", false, // ex. primary carrying secondary
		domain, contain, "in", true,
		domain, support, "on", true, // on the x are the w,y,z
		domain, suspect, "suspicious of", false,
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_noun",
		"ROWID", "domain", "noun", "kind"},
		//
		message, domain, "message", things,
		missive, domain, "missive", things,
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_name",
		"domain", "noun", "name", "rank"},
		//
		domain, message, "message", 0,
		domain, missive, "missive", 0,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
