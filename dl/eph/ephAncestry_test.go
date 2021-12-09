package eph

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestry(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{"k", ""},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{"q", "j"}, // same domain
		&EphKinds{"n", "k"}, // root domain
		&EphKinds{"j", "m"}, // parent domain
	)
	out := testOut{mdl_kind}
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
		&EphKinds{"t", ""},
		&EphKinds{"p", "t"},
		&EphKinds{"t", "p"},
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
		&EphKinds{"t", ""},
		&EphKinds{"p", "t"},
		&EphKinds{"q", "t"},
		&EphKinds{"k", "p"},
		&EphKinds{"k", "q"},
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
		&EphKinds{"k", ""},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{"n", "k"}, // root domain
		&EphKinds{"n", "m"}, // more specific
		&EphKinds{"n", "k"}, // duped
	)
	out := testOut{mdl_kind}
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
		&EphKinds{"k", ""},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "x"},
	)
	out := testOut{mdl_domain}
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
		&EphKinds{"k", ""},
		&EphKinds{"q", "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "k"},
		&EphKinds{"n", "k"},
		&EphKinds{"q", "k"}, // should be okay.
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{"m", "n"},
	)
	out := testOut{mdl_domain}
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
		&EphKinds{"k", ""},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{"m", "k"}, // second in a parallel domain should be fine
	)
	out := testOut{mdl_kind}
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
		&EphKinds{"k", ""},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{"m", "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{"m", "k"}, // re: TestScopedRivalsOkay, should be okay.
	)
	dt.makeDomain(dd("z", "b", "d")) // trying to include both should be a problem; they are two unique kinds...
	out := testOut{mdl_domain}
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
