package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var dt domainTest
	// yes, these are collective nouns not plurals... shhh...
	dt.makeDomain(ds("a"),
		&EphPlural{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&EphPlural{Singular: "bat", Plural: "cloud"},
		&EphPlural{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(ds("b", "a"),
		// add something new:
		&EphPlural{Singular: "fish", Plural: "school"},
	)
	dt.makeDomain(ds("c", "a"),
		// redefine:
		&EphPlural{Singular: "witch", Plural: "unkindness"},
		// collapse:
		&EphPlural{Singular: "bat", Plural: "cauldron"},
	)
	var out testOut
	cat := Catalog{Writer: &out}
	// addToCat will define all of the above domains
	// and then it will.... queue all of the ephemera
	// FIX? should end domain write the domain? hmm... maybe?
	// ( and for example -- run any queued commands? )
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		// walk the domains and run the commands remaining in their queues
		var err error
		for _, n := range ds("a", "b", "c") {
			d := cat.GetDomain(n)
			for _, el := range d.eph.All {
				if e := el.Eph.Catalog(&cat, d, el.At); e != nil {
					err = e
					break
				}
			}
		}
		if err != nil {
			t.Fatal(err)
		} else {
			// try seeing what we made
			got := out[mdl_plural]
			if diff := pretty.Diff(got, []outEl{{
				"a", "unkindness", "raven",
			}, {
				"a", "cloud", "bat",
			}, {
				"a", "cauldron", "bat",
			}, {
				"b", "school", "fish",
			}, {
				// we dont expect to see our duplicated definition of cauldron of bat(s)
				// we do expect that its okay to redefine the collective "witch" as "unkindness"
				// ( but wicca good and love the earth, and i'll be over here. )
				"c", "unkindness", "witch",
			}}); len(diff) > 0 {
				t.Log(pretty.Sprint(got))
				t.Fatal(diff)
			}
		}
	}
}

func TestPluralDomainConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("a"),
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&EphPlural{Singular: "raven", Plural: "unkindness"},
		&EphPlural{Singular: "witch", Plural: "unkindness"},
	)
	var out testOut
	cat := Catalog{Writer: &out}
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		// walk the domains and run the commands remaining in their queues
		var err error
		for _, n := range ds("a", "b", "c") {
			d := cat.GetDomain(n)
			for _, el := range d.eph.All {
				if e := el.Eph.Catalog(&cat, d, el.At); e != nil {
					err = e
					break
				}
			}
		}
		if err == nil {
			t.Fatal("expected an error")
		} else {
			t.Log("ok:", err)
		}
	}
}
