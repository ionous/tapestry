package eph

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestry(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kind: "q", Ancestor: "j"}, // same domain
		&EphKinds{Kind: "n", Ancestor: "k"}, // root domain
		&EphKinds{Kind: "j", Ancestor: "m"}, // parent domain
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
		&EphKinds{Kind: "t"},
		&EphKinds{Kind: "p", Ancestor: "t"},
		&EphKinds{Kind: "t", Ancestor: "p"},
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
		&EphKinds{Kind: "t"},
		&EphKinds{Kind: "p", Ancestor: "t"},
		&EphKinds{Kind: "q", Ancestor: "t"},
		&EphKinds{Kind: "k", Ancestor: "p"},
		&EphKinds{Kind: "k", Ancestor: "q"},
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
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kind: "n", Ancestor: "k"}, // root domain
		&EphKinds{Kind: "n", Ancestor: "m"}, // more specific
		&EphKinds{Kind: "n", Ancestor: "k"}, // duped
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
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "x"},
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
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "q", Ancestor: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"},
		&EphKinds{Kind: "n", Ancestor: "k"},
		&EphKinds{Kind: "q", Ancestor: "k"}, // should be okay.
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kind: "m", Ancestor: "n"}, // should fail.
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
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"}, // second in a parallel domain should be fine
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
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kind: "m", Ancestor: "k"}, // re: TestScopedRivalsOkay, should be okay.
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
