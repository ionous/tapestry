package eph

import (
	"errors"
	"testing"

	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var warnings []error
	unwarn := catchWarnings(&warnings)
	defer unwarn()
	//
	var dt domainTest
	// yes, these are collective nouns not plurals... shhh...
	dt.makeDomain(dd("a"),
		&EphPlural{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&EphPlural{Singular: "bat", Plural: "cloud"},
		&EphPlural{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("b", "a"),
		// add something new:
		&EphPlural{Singular: "fish", Plural: "school"},
	)
	dt.makeDomain(dd("c", "a"),
		// redefine:
		&EphPlural{Singular: "witch", Plural: "unkindness"},
		// collapse:
		&EphPlural{Singular: "bat", Plural: "cauldron"},
	)
	var out testOut
	var cat Catalog
	// addToCat will define all of the above domains
	// and then it will.... queue all of the ephemera
	// FIX? should end domain write the domain? hmm... maybe?
	// ( and for example -- run any queued commands? )
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.ProcessDomains(nil); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Redefined, warnings[0]); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings[1]); e != nil {
		t.Fatal(e)
	} else {
		if e := cat.plurals.WritePlurals(&out); e != nil {
			// try seeing what we made
			if diff := pretty.Diff(out, testOut{
				"a:unkindness:raven",
				"a:cloud:bat",
				"a:cauldron:bat",
				"b:school:fish",
				// we dont expect to see our duplicated definition of cauldron of bat(s)
				// we do expect that it's okay to redefine the collective "witch" as "unkindness"
				// ( wicca good and love the earth, and i'll be over here. )
				"c:unkindness:witch",
			}); len(diff) > 0 {
				t.Log(pretty.Sprint(out))
				t.Fatal(diff)
			}
		}
	}
}

func okDomainConflict(d string, y ReasonForConflict, e error) (err error) {
	var de DomainError
	var conflict *Conflict
	if !errors.As(e, &de) || de.Domain != d ||
		!errors.As(de.Err, &conflict) || conflict.Reason != y {
		err = errutil.New("unexpected warning", e)
	}
	return
}

func TestPluralDomainConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&EphPlural{Singular: "raven", Plural: "unkindness"},
		&EphPlural{Singular: "witch", Plural: "unkindness"},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.ProcessDomains(nil); e == nil {
		t.Fatal("expected an error")
	} else {
		t.Log("ok:", e)
	}
}
