package eph

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	// because this test picks out two warnings, one by one...
	// we cant shuffle the statements...
	dt := domainTest{noShuffle: true}
	// yes, these are collective nouns not plurals... shhh...
	dt.makeDomain(dd("a"),
		&EphPlurals{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&EphPlurals{Singular: "bat", Plural: "cloud"},
		&EphPlurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("b", "a"),
		// add something new:
		&EphPlurals{Singular: "fish", Plural: "school"},
	)
	dt.makeDomain(dd("c", "a"),
		// redefine:
		&EphPlurals{Singular: "witch", Plural: "unkindness"},
		// collapse:
		&EphPlurals{Singular: "bat", Plural: "cauldron"},
	)
	var cat Catalog
	// addToCat will define all of the above domains
	// and then it will.... queue all of the ephemera
	// FIX? should end domain write the domain? hmm... maybe?
	// ( and for example -- run any queued commands? )
	if e := dt.addToCat(cat.Weaver()); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(nil, nil); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Redefined, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Plural}
		if e := cat.WritePlurals(&out); e != nil {
			t.Fatal(e)
		} else {
			if diff := pretty.Diff(out[1:], testOut{
				"a:unkindness:raven:x",
				"a:cloud:bat:x",
				"a:cauldron:bat:x",
				"b:school:fish:x",
				// we dont expect to see our duplicated definition of cauldron of bat(s)
				// we do expect that it's okay to redefine the collective "witch" as "unkindness"
				// ( wicca good and love the earth, and i'll be over here. )
				"c:unkindness:witch:x",
			}); len(diff) > 0 {
				t.Log(pretty.Sprint(out))
				t.Fatal(diff)
			}
		}
	}
}

func okDomainConflict(d string, y ReasonForConflict, e error) (err error) {
	var de domainError
	var conflict *Conflict
	if !errors.As(e, &de) || de.Domain != d ||
		!errors.As(de.Err, &conflict) || conflict.Reason != y {
		err = errutil.New("unexpected warning:", e)
	}
	return
}

func TestPluralDomainConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&EphPlurals{Singular: "raven", Plural: "unkindness"},
		&EphPlurals{Singular: "witch", Plural: "unkindness"},
	)
	var cat Catalog
	if e := dt.addToCat(cat.Weaver()); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(nil, nil); e == nil {
		t.Fatal("expected an error")
	} else {
		t.Log("ok:", e)
	}
}
