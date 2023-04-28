package weave

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// ensure fields which reference aspects use the necessary formatting
func AspectParam(aspectName string) eph.Params {
	return eph.Params{Name: aspectName, Affinity: Affinity{eph.Affinity_Text}, Class: aspectName}
}

// generate a kind of aspect containing a few traits.
func TestAspectFormation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("d"),
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("a")}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("a")}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "one", Affinity: Affinity{eph.Affinity_Text}}}},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestAspectConflictingTraits(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{ // add some traits
			"one", "several", "oh so many", //
		}},
		&eph.Kinds{Kind: "b", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "b", Traits: []string{ // add some traits
			"one", "two", "blue", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("a")}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("b")}},
	)
	var conflict *Conflict
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}
