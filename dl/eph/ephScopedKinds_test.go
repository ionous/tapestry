package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestScopedKinds(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kinds: "q", From: "j"}, // same domain
		&EphKinds{Kinds: "n", From: "k"}, // root domain
		&EphKinds{Kinds: "j", From: "m"}, // parent domain
	)
	var out testOut
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"k:", "m:k", "j:m", "q:j", "n:k",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func TestScopedRedundant(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kinds: "n", From: "k"}, // root domain
		&EphKinds{Kinds: "n", From: "m"}, // more specific
		&EphKinds{Kinds: "n", From: "k"}, // duped
	)
	var out testOut
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"k:", "m:k", "n:m",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func TestScopedKindMissing(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "x"},
	)
	var out testOut
	if e := writeKinds(dt, &out); e == nil || e.Error() != `unknown dependency "x" for kind "m"` {
		t.Fatal("expected error", e, out)
	} else {
		t.Log("ok", e)
	}
}

func TestScopedKindConflict(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "q", From: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"},
		&EphKinds{Kinds: "n", From: "k"},
		&EphKinds{Kinds: "q", From: "k"}, // should be okay.
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kinds: "m", From: "n"},
	)
	var out testOut
	if e := writeKinds(dt, &out); e == nil || e.Error() != `can't redefine parent as "n" for kind "m"` {
		t.Fatal("expected error", e, out)
	} else if warned := warnings.all(); len(warned) != 1 {
		t.Fatal("expected one warning", warned)
	} else {
		t.Log("ok", e, warned)
	}
}

func TestScopedRivalsOkay(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kinds: "m", From: "k"}, // second in a parallel domain should be fine
	)
	var out testOut
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	}
}

func TestScopedRivalConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kinds: "m", From: "k"}, // re: TestScopedRivalsOkay, should be okay.
	)
	dt.makeDomain(dd("z", "b", "d")) // trying to include both should be a problem; they are two unique kinds...
	var out testOut
	if e := writeKinds(dt, &out); e == nil {
		t.Fatal("expected an error", out)
	} else {
		t.Log("ok", e)
	}
}

func writeKinds(dt domainTest, pout *testOut) (err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else {
		err = cat.AssembleCatalog(PhaseActions{
			AncestryPhase: PhaseAction{
				PhaseFlags{NoDuplicates: true},
				func(d *Domain) (err error) {
					if ks, e := d.ResolveKinds(); e != nil {
						err = e
					} else if e := ks.WriteTable(pout, "", false); e != nil {
						err = e
					}
					return
				},
			},
		})
	}
	return
}
