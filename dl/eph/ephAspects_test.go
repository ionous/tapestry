package eph

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// generate a kind of aspect containing a few traits.
func TestAspectFormation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("d"),
		&EphKinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&EphKinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	out := testOut{mdl.Field}
	if cat, e := buildAncestors(&dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"d:a:one:bool::x",
		"d:a:several:bool::x",
		"d:a:oh_so_many:bool::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// generate a some other kind that has a field of using that aspect.
// ( not a very exciting test )
func TestAspectUsage(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&EphKinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{AspectParam("a")}},
	)
	out := testOut{mdl.Field}
	if cat, e := buildAncestors(&dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:a:one:bool::x",
		"a:a:several:bool::x",
		"a:a:oh_so_many:bool::x",
		"b:k:a:text:a:x", // MISSING THIS.... fields for the aspect
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fail to generate a kind that has a field named the same as a trait.
func TestAspectConflictingFields(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&EphKinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{AspectParam("a")}},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "one", Affinity: Affinity{Affinity_Text}}}},
	)
	if _, e := buildAncestors(&dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestAspectConflictingTraits(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&EphKinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{ // add some traits
			"one", "several", "oh so many", //
		}},
		&EphKinds{Kind: "b", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&EphAspects{Aspects: "b", Traits: []string{ // add some traits
			"one", "two", "blue", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{AspectParam("a")}},
		&EphKinds{Kind: "k", Contain: []EphParams{AspectParam("b")}},
	)
	var conflict *Conflict
	if _, e := buildAncestors(&dt); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}
