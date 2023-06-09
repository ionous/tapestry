package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralConflict(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// one singular can have several plurals:
		// exTestAncestryMultipleParents. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict`); !ok {
		t.Fatal("expected error; got:", e)
	} else {
		t.Log("ok:", e)
	}
}

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	// because this test picks out two warnings, one by one...
	// we cant shuffle the statements...
	dt := newTestShuffle(t.Name(), false)
	defer dt.Close()
	// yes, these are collective nouns not plurals... shhh...
	dt.makeDomain(dd("a"),
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&eph.Plurals{Singular: "bat", Plural: "cloud"},
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("b", "a"),
		// add something new:
		&eph.Plurals{Singular: "fish", Plural: "school"},
		// collapse:
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("c", "a"),
		// redefine; this isnt allowed; but everything else should have worked.
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	//
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	} else if ok, e := okError(t, warnings.shift(), `Duplicate plural "cauldron"`); !ok {
		t.Fatal(e)
	} else if out, e := dt.readPlurals(); e != nil {
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
