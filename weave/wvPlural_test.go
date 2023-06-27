package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralConflict(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var warnings testweave.Warnings
	unwarn := warnings.Catch(t)
	defer unwarn()
	// cant shuffle because the test picks out warnings one by one
	// tbd: add some warnings.Includes()?
	dt := testweave.NewWeaverOptions(t.Name(), nil, false)
	defer dt.Close()
	// yes, these are collective nouns not plurals... shhh...
	dt.MakeDomain(dd("a"),
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&eph.Plurals{Singular: "bat", Plural: "cloud"},
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.MakeDomain(dd("b", "a"),
		// add something new:
		&eph.Plurals{Singular: "fish", Plural: "school"},
		// collapse:
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.MakeDomain(dd("c", "a"),
		// redefine; this isnt allowed; but everything else should have worked.
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	//
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	} else if ok, e := testweave.OkayError(t, warnings.Shift(), `Duplicate plural "cauldron"`); !ok {
		t.Fatal(e)
	} else if out, e := dt.ReadPlurals(); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(out, []string{
			"a:unkindness:raven",
			"a:cloud:bat",
			"a:cauldron:bat",
			"b:school:fish",
			// plural redefinition is (no longer) allowed.
			// ( wicca good and love the earth: and i'll be over here. )
			// "c:unkindness:witch",
			// we dont expect to see our duplicated definition of cauldron of bat(s)
			// c is dependent on a: so the definition would be redundant.
			// "c:cauldron:bat",
		}); len(diff) > 0 {
			t.Log("got", len(out), out)
			t.Fatal(diff)
		}
	}
}
