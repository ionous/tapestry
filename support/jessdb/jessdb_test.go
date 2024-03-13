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

			for i, p := range jesstest.Phrases {
				if str, ok := p.Test(); !ok {
					continue // skip unused tests
				} else {
					// reset the dynamic noun pool every test
					dynamicNouns := make(map[string]string)
					q := testAdapter{jessdb.MakeQuery(m, "a"), dynamicNouns}
					// create the test helper
					m := jesstest.MakeMock(q, dynamicNouns)
					// run the test:
					t.Logf("testing: %d %s", i, str)
					if !p.Verify(m.Generate(str)) {
						t.Logf("failed %d", i)
						t.Fail()
					}
				}
			}
		}
	}
}

type testAdapter struct {
	jess.Query
	dynamicNouns map[string]string
}

func (ta testAdapter) FindNoun(ws match.Span, kind string) (ret string, width int) {
	if n, w := ta.Query.FindNoun(ws, kind); w > 0 {
		ret, width = n, w
	} else if noun, ok := ta.dynamicNouns[ws.String()]; ok {
		ret, width = noun, len(ws)
	}
	return
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
		kinds = iota + 1
		patterns
		macros
		aspects
		traits
		records
		actions
		// kinds of nouns
		objects // base for concrete and abstract nouns
		things
		containers
		supporters
		colors     // aspects
		groups     // records
		directions // umm... directions
		doors
		rooms
		// actions
		storing
		// macros
		carry
		contain
		support
		suspect
		// nouns
		story
		message // thing
		missive
		north // directions
		south
		east
		west
		river // predefined door
		ocean // predefined room
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
		//
		aspects, domain, "aspects", "aspect", idPath(),
		traits, domain, "traits", "", idPath(aspects),
		colors, domain, "color", nil, idPath(aspects),
		//
		patterns, domain, "patterns", "pattern", idPath(),
		macros, domain, "macros", "macro", idPath(patterns),
		actions, domain, "actions", "action", idPath(patterns),
		//
		kinds, domain, "kinds", "kind", idPath(),
		objects, domain, "objects", "object", idPath(kinds),
		directions, domain, "directions", "direction", idPath(objects, kinds),
		rooms, domain, "rooms", "room", idPath(objects, kinds),
		things, domain, "things", "thing", idPath(objects, kinds),
		doors, domain, "doors", "door", idPath(objects, things, kinds),
		containers, domain, "containers", "container", idPath(objects, things, kinds),
		supporters, domain, "supporters", "supporter", idPath(objects, things, kinds),
		//
		records, domain, "records", "record", idPath(),
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
		domain, traits, "dark", "bool",
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
		"domain", "ROWID", "noun", "kind"},
		//
		domain, story, "story", things,
		domain, message, "message", things,
		domain, missive, "missive", things,
		domain, north, "north", directions,
		domain, west, "west", directions,
		domain, east, "east", directions,
		domain, south, "south", directions,
		//
		domain, river, "river", doors,
		domain, ocean, "ocean", rooms,
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_name",
		"domain", "noun", "name", "rank"},
		//
		domain, story, "story", 0,
		domain, message, "message", 0,
		domain, missive, "missive", 0,
		domain, north, "north", 0,
		domain, west, "west", 0,
		domain, east, "east", 0,
		domain, south, "south", 0,
		domain, river, "river", 0,
		domain, ocean, "ocean", 0,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
