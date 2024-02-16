package grokdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/grok"
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
		groktest.RunPhraseTests(t, func(testPhrase string) (ret grok.Results, err error) {
			if ws, e := grok.MakeSpan(testPhrase); e != nil {
				err = e
			} else {
				ret, err = jess.Match(&x, ws)
			}
			return
		})
	}
}

func TestTraits(t *testing.T) {
	if db, e := setupDB(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		x := dbSource{domain: "a", db: tables.NewCache(db)}
		groktest.RunTraitTests(t, func(testPhrase string) (ret grok.TraitSet, err error) {
			if ws, e := grok.MakeSpan(testPhrase); e != nil {
				err = e
			} else {
				var t jess.TraitsKind //
				input := jess.InputState(ws)
				if !t.Match(jess.MakeQuery(&x), &input) {
					err = errors.New("failed to match traits")
				} else if cnt := len(input); cnt != 0 {
					err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
				} else {
					ret = t.GetTraitSet()
				}
			}
			return
		})
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
		implies
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
		aspects, domain, "aspects", "aspect", idPath(),
		patterns, domain, "patterns", "pattern", idPath(),
		macros, domain, "macros", "macro", idPath(patterns),
		kinds, domain, "kinds", "kind", idPath(),
		traits, domain, "traits", "", idPath(aspects),
		things, domain, "things", "thing", idPath(kinds),
		containers, domain, "containers", "container", idPath(things, kinds),
		supporters, domain, "supporters", "supporter", idPath(things, kinds),
		// macros:
		carry, domain, "carry", "", idPath(macros),
		contain, domain, "contain", "", idPath(macros),
		inherit, domain, "inherit", "", idPath(macros),
		implies, domain, "implies", "", idPath(macros),
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
		domain, traits, "fixed in place", "bool",
		// macros:
		domain, carry, "primary", "text",
		domain, carry, "secondary", "text_list",
		domain, carry, "error", "text",
		//
		domain, contain, "primary", "text",
		domain, contain, "secondary", "text_list",
		domain, contain, "error", "text",
		//
		domain, inherit, "primary", "text_list",
		domain, inherit, "error", "text",

		domain, implies, "primary", "text_list",
		domain, implies, "error", "text",
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
		inherit, "error",
		implies, "error",
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
		domain, inherit, "kinds of", false, // for "are kinds of containers"
		domain, inherit, "a kind of", false, // ex. "a kind of container"
		domain, implies, "usually", false, // ex. "usually closed"
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
