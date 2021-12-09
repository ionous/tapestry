package eph

import (
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
	if cat, e := buildFields(dt); e != nil {
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
		&EphAspects{Aspects: "a", Traits: []string{ // fix? should "Aspects" be singular?
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "a", Affinity: Affinity{Affinity_Text}, Class: "a"},
	)
	out := testOut{mdl_field}
	if cat, e := buildFields(dt); e != nil {
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

// fail to generate a kind that has a normal field named the same as the aspect.
func TestAspectConflictingFields(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfAspect},            // say that aspects exist
		&EphKinds{Kinds: "a", From: KindsOfAspect}, // make an aspect
		&EphAspects{Aspects: "a", Traits: []string{ // fix? should "Aspects" be singular?
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "a", Affinity: Affinity{Affinity_Text}, Class: "a"},
		&EphFields{Kinds: "k", Name: "one", Affinity: Affinity{Affinity_Text}},
	)
	if _, e := buildFields(dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

// 4. fail to generate a kind that has a field named the same as a trait.
// x. how do you write the literal back into the db -- its not a "value",
//    do you need an StoredValue() interface or something :/
