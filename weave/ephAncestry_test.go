package weave

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestry(t *testing.T) {
	dt := newTest(t.Name())
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
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Kind}
		if e := cat.WriteKinds(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:k::x", "b:m:k:x", "c:j:m,k:x", "c:q:j,m,k:x", "c:n:k:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// cycles should fail
func TestAncestryCycle(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "t", Ancestor: "p"},
	)
	_, e := dt.Assemble()
	if e == nil || !strings.HasPrefix(e.Error(), "circular reference detected") {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// multiple inheritance should fail
func TestAncestryMultipleParents(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "q", Ancestor: "t"},
		&eph.Kinds{Kind: "k", Ancestor: "p"},
		&eph.Kinds{Kind: "k", Ancestor: "q"},
	)
	_, e := dt.Assemble()
	if e == nil || e.Error() != `"k" has more than one parent` {
		t.Fatal("expected failure", e)
	} else {
		t.Log("ok", e)
	}
}

// respecifying hierarchy is okay
func TestAncestryRedundancy(t *testing.T) {
	dt := newTest(t.Name())
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
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Kind}
		if e := cat.WriteKinds(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:k::x", "b:m:k:x", "c:n:m,k:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

func TestAncestryMissing(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "x"},
	)
	_, e := dt.Assemble()
	if e == nil || e.Error() != `unknown dependency "x" for kind "m"` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok", e)
	}
}

func TestAncestryRedefined(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	//
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
		&eph.Kinds{Kind: "n", Ancestor: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"}, // should be okay. << --- failing here
	)
	dt.makeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "m", Ancestor: "n"}, // should fail.
	)
	if _, e := dt.Assemble(); e == nil || e.Error() != `can't redefine parent as "n" for kind "m"` {
		t.Fatal("expected error; got:", e)
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
	dt := newTest(t.Name())
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

	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Kind}
		if e := cat.WriteKinds(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:k::x", "b:m:k:x", "d:m:k:x",
		}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
}

// two different kinds named in two different parent trees should fail
func TestAncestryRivalConflict(t *testing.T) {
	dt := newTest(t.Name())
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
	//
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected an error")
	}
}

func newTest(name string) *domainTest {
	db := testdb.Open(name, testdb.Memory, "")
	return &domainTest{
		name: name,
		db:   db,
		cat:  NewCatalog(db),
	}
}

func newTestShuffle(name string, shuffle bool) *domainTest {
	db := testdb.Open(name, testdb.Memory, "")
	return &domainTest{
		name:      name,
		db:        db,
		cat:       NewCatalog(db),
		noShuffle: !shuffle,
	}
}
