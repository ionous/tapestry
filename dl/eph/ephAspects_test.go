package eph

import (
	"errors"
	"testing"

	"github.com/kr/pretty"
)

// generate a kind of aspect containing a few traits.
func TestAspectFormation(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfAspect},            // say that aspects exist
		&EphKinds{Kinds: "a", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{ // fix? should "Aspects" be singular?
			"one", "several", "oh so many", //
		}},
	)
	out := testOut{mdl_aspect}
	if cat, e := buildAncestors(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteAspects(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:one:0",
		"a:several:1",
		"a:oh_so_many:2",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// generate a some other kind that has a field of using that aspect.
// ( not a very exciting test )
func TestAspectUsage(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfAspect},            // say that aspects exist
		&EphKinds{Kinds: "a", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "k", Contain: []EphParams{AspectParam("a")}},
	)
	out := testOut{mdl_field}
	if cat, e := buildAncestors(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"b:k:a:text:a:x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fail to generate a kind that has a field named the same as a trait.
func TestAspectConflictingFields(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfAspect},            // say that aspects exist
		&EphKinds{Kinds: "a", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "k", Contain: []EphParams{AspectParam("a")}},
		&EphKinds{Kinds: "k", Contain: []EphParams{{Name: "one", Affinity: Affinity{Affinity_Text}}}},
	)
	if _, e := buildAncestors(dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestAspectConflictingTraits(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfAspect},            // say that aspects exist
		&EphKinds{Kinds: "a", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{ // add some traits
			"one", "several", "oh so many", //
		}},
		&EphKinds{Kinds: "b", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "b", Traits: []string{ // add some traits
			"one", "two", "blue", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "k", Contain: []EphParams{AspectParam("a")}},
		&EphKinds{Kinds: "k", Contain: []EphParams{AspectParam("b")}},
	)
	var conflict *Conflict
	if _, e := buildAncestors(dt); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}
