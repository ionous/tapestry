package eph

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestry(t *testing.T) {
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
	out := testOut{mdl.Kind}
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k::x", "b:m:k:x", "c:j:m,k:x", "c:q:j,m,k:x", "c:n:k:x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// cycles should fail
func TestAncestryCycle(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "t"},
		&EphKinds{Kinds: "p", From: "t"},
		&EphKinds{Kinds: "t", From: "p"},
	)
	if _, e := buildAncestors(dt); e == nil || !strings.HasPrefix(e.Error(), "circular reference detected") {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// multiple inheritance should fail
func TestAncestryMultipleParents(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "t"},
		&EphKinds{Kinds: "p", From: "t"},
		&EphKinds{Kinds: "q", From: "t"},
		&EphKinds{Kinds: "k", From: "p"},
		&EphKinds{Kinds: "k", From: "q"},
	)
	if _, e := buildAncestors(dt); e == nil || e.Error() != `"k" has more than one parent` {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// respecifying hierarchy is okay
func TestAncestryRedundancy(t *testing.T) {
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
	out := testOut{mdl.Kind}
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k::x", "b:m:k:x", "c:n:m,k:x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func TestAncestryMissing(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "x"},
	)
	out := testOut{mdl.Domain}
	if e := writeKinds(dt, &out); e == nil || e.Error() != `unknown dependency "x" for kind "m"` {
		t.Fatal("expected error", e, out)
	} else {
		t.Log("ok", e)
	}
}

func TestAncestryRedefined(t *testing.T) {
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
		&EphKinds{Kinds: "m", From: "n"}, // should fail.
	)
	out := testOut{mdl.Domain}
	if e := writeKinds(dt, &out); e == nil || e.Error() != `can't redefine parent as "n" for kind "m"` {
		t.Fatal("expected error", e, out)
	} else if warned := warnings.all(); len(warned) != 1 {
		t.Fatal("expected one warning", warned)
	} else {
		t.Log("ok", e, warned)
	}
}

// a similar named kind in another domain should be fine
func TestAncestryRivalsOkay(t *testing.T) {
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
	out := testOut{mdl.Kind}
	if e := writeKinds(dt, &out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k::x", "b:m:k:x", "d:m:k:x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// two different kinds named in two different parent trees should fail
func TestAncestryRivalConflict(t *testing.T) {
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
	out := testOut{mdl.Domain}
	if e := writeKinds(dt, &out); e == nil {
		t.Fatal("expected an error", out)
	} else {
		t.Log("ok", e)
	}
}

func writeKinds(dt domainTest, w *testOut) (err error) {
	if cat, e := buildAncestors(dt); e != nil {
		err = e
	} else {
		err = cat.WriteKinds(w)
	}
	return
}
