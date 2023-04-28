package weave

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestry(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "q", Ancestor: "j"}, // same domain
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // root domain
		&eph.Kinds{Kind: "j", Ancestor: "m"}, // parent domain
	)
	out := testOut{mdl.Kind}
	if e := writeKinds(&dt, &out); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "t", Ancestor: "p"},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil || !strings.HasPrefix(e.Error(), "circular reference detected") {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// multiple inheritance should fail
func TestAncestryMultipleParents(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "q", Ancestor: "t"},
		&eph.Kinds{Kind: "k", Ancestor: "p"},
		&eph.Kinds{Kind: "k", Ancestor: "q"},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil || e.Error() != `"k" has more than one parent` {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// respecifying hierarchy is okay
func TestAncestryRedundancy(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // root domain
		&eph.Kinds{Kind: "n", Ancestor: "m"}, // more specific
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // duped
	)
	out := testOut{mdl.Kind}
	if e := writeKinds(&dt, &out); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "x"},
	)
	var out testOut
	if e := writeKinds(&dt, &out); e == nil || e.Error() != `unknown dependency "x" for kind "m"` {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
		&eph.Kinds{Kind: "n", Ancestor: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"}, // should be okay.
	)
	dt.makeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "m", Ancestor: "n"}, // should fail.
	)
	var out testOut
	if e := writeKinds(&dt, &out); e == nil || e.Error() != `can't redefine parent as "n" for kind "m"` {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // second in a parallel domain should be fine
	)
	out := testOut{mdl.Kind}
	if e := writeKinds(&dt, &out); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
	)
	dt.makeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // re: TestScopedRivalsOkay, should be okay.
	)
	dt.makeDomain(dd("z", "b", "d")) // trying to include both should be a problem; they are two unique kinds...
	var out testOut
	if e := writeKinds(&dt, &out); e == nil {
		t.Fatal("expected an error", out)
	} else {
		t.Log("ok", e)
	}
}

func writeKinds(dt *domainTest, w *testOut) (err error) {
	cat := NewCatalog(dt.Open("kinds"))
	if e := dt.addToCat(cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(); e != nil {
		err = e
	} else {
		err = cat.WriteKinds(w)
	}
	return
}
