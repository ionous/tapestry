package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// ensure fields which reference aspects use the necessary formatting
var AspectParam = eph.AspectParam

// generate a kind of aspect containing a few traits.
func TestAspectFormation(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("d"),
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{
			"one", "several", "oh so many", //
		}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		// names are sorted
		"d:a:oh_so_many:bool:",
		"d:a:one:bool:",
		"d:a:several:bool:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// generate a some other kind that has a field of using that aspect.
// ( not a very exciting test )
func TestAspectUsage(t *testing.T) {
	dt := newTest(t.Name())
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

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:a:oh_so_many:bool:",
		"a:a:one:bool:",
		"a:a:several:bool:",
		"b:k:a:text:a",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fail to generate a kind that has a field named the same as a trait.
func TestAspectConflictingFields(t *testing.T) {
	dt := newTest(t.Name())
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
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "one", Affinity: affine.Text}}},
	)
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected error")
	} else if ok, e := okError(t, e, "conflict:"); !ok {
		t.Fatal(e)
	}
}

// try to add two aspects with conflicting traits to a single kind.
func TestAspectConflictingTraits(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Aspect.String()},                // say that aspects exist
		&eph.Kinds{Kind: "a", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "a", Traits: []string{ // add some traits
			"one", "several", "oh so many", //
		}},
		&eph.Kinds{Kind: "b", Ancestor: kindsOf.Aspect.String()}, // make an aspect
		&eph.Aspects{Aspects: "b", Traits: []string{ // add some different traits
			"one", "two", "blue", //
		}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("a")}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{AspectParam("b")}},
	)
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected error")
	} else if ok, e := okError(t, e, "conflict:"); !ok {
		t.Fatal(e)
	}
}