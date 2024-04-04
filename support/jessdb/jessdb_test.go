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
			const at = -1
			for i, p := range jesstest.Phrases {
				if i != at && at >= 0 {
					continue
				}
				if str, ok := p.Test(); !ok {
					continue // skip unused tests
				} else {
					// reset the dynamic noun pool every test
					dynamicNouns := make(map[string]string)
					// block logging of known nouns to match jess_test
					for _, name := range []string{"story"} {
						dynamicNouns[name] = name
						dynamicNouns["$"+name] = "things"
					}
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

func (ta testAdapter) FindNoun(ws match.Span, pkind *string) (ret string, width int) {
	if n, w := ta.Query.FindNoun(ws, pkind); w > 0 {
		ret, width = n, w
	} else {
		str := ws.String()
		if noun, ok := ta.dynamicNouns[str]; ok {
			ret, width = noun, len(ws)
			if pkind != nil {
				*pkind = ta.dynamicNouns["$"+str]
			}
		}
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
		// domain string
		domain = "a"
		// built in kinds
		kinds = iota + 1
		patterns
		relations
		verbs
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
		relSuspicion
		relWhereabouts
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
		// used when matching the names of verbs
		verbCarrying
		verbCarriedBy
		verbIn
		verbOn
		verbSuspiciousOf
		//
		fieldOpen
		fieldClosed
		fieldOpenable
		fieldTransparent
		fieldFixed
		fixedDark
		fieldDesc
		fieldTitle
		fieldAge
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
		verbs, domain, "verbs", "verb", idPath(),
		aspects, domain, "aspects", "aspect", idPath(),
		traits, domain, "traits", "", idPath(aspects),
		colors, domain, "color", nil, idPath(aspects),
		//
		relations, domain, "relations", "relation", idPath(),
		patterns, domain, "patterns", "pattern", idPath(),
		actions, domain, "actions", "action", idPath(patterns),
		storing, domain, "storing", nil, idPath(actions),
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
		// relations:
		relWhereabouts, domain, "whereabouts", nil, idPath(relations),
		relSuspicion, domain, "suspects", nil, idPath(relations),
		//,
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_field",
		"domain", "kind", "ROWID", "field", "affinity"},
		// traits: for matching traits to set
		domain, traits, fieldClosed, "closed", "bool",
		domain, traits, fieldOpen, "open", "bool",
		domain, traits, fieldOpenable, "openable", "bool",
		domain, traits, fieldTransparent, "transparent", "bool",
		domain, traits, fieldFixed, "fixed in place", "bool",
		domain, traits, fixedDark, "dark", "bool",
		// fields: for matching properties to set
		domain, things, fieldDesc, "description", "text",
		domain, things, fieldTitle, "title", "text",
		domain, things, fieldAge, "age", "text",
	); e != nil {
		err = e
	} else if e := testdb.Ins(db, []string{"mdl_noun",
		"domain", "ROWID", "noun", "kind"},
		// things
		domain, story, "story", things,
		domain, message, "message", things,
		domain, missive, "missive", things,
		// directions
		domain, north, "north", directions,
		domain, west, "west", directions,
		domain, east, "east", directions,
		domain, south, "south", directions,
		// links
		domain, river, "river", doors,
		domain, ocean, "ocean", rooms,
		// verbs
		domain, verbCarrying, "carrying", verbs,
		domain, verbCarriedBy, "carried by", verbs,
		domain, verbIn, "in", verbs,
		domain, verbOn, "on", verbs,
		domain, verbSuspiciousOf, "suspicious of", verbs,
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
		domain, verbCarrying, "carrying", 0,
		domain, verbCarriedBy, "carried by", 0,
		domain, verbIn, "in", 0,
		domain, verbOn, "on", 0,
		domain, verbSuspiciousOf, "suspicious of", 0,
	); e != nil {
		err = e
	}
	if err != nil {
		db.Close()
	}
	return db, err
}
