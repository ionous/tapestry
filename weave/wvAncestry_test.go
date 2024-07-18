package weave_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestAncestryFormation(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.MakeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "q", Ancestor: "j"}, // same domain
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // root domain
		&eph.Kinds{Kind: "j", Ancestor: "m"}, // parent domain
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal("failed assembly", e)
	} else if out, e := dt.ReadKinds(); e != nil {
		t.Fatal(e)
	} else if !reflect.DeepEqual(out, []string{
		// domain, kind
		"a:k:",
		"b:m:k",
		"c:n:k",
		"c:j:m,k",
		"c:q:j,m,k",
	}) {
		t.Fatalf("%#v", out)
	}
}

// cycles should fail
func TestAncestryCycle(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "t", Ancestor: "p"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict circular reference detected`); !ok {
		t.Fatal("unexpected error:", e)
	}
}

// multiple inheritance should fail
func TestAncestryMultipleParents(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "q", Ancestor: "t"},
		&eph.Kinds{Kind: "k", Ancestor: "p"},
		&eph.Kinds{Kind: "k", Ancestor: "q"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Missing a definition in domain "a" that would allow "k"`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

// respecifying hierarchy is okay
func TestAncestryRedundancy(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // parent/root domain
	)
	dt.MakeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // root domain
		&eph.Kinds{Kind: "n", Ancestor: "m"}, // more specific
		&eph.Kinds{Kind: "n", Ancestor: "k"}, // duped
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadKinds(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:", "b:m:k", "c:n:m,k",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}

}

func TestAncestryMissing(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	//
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "x"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Missing kind "x" in domain "b"`); !ok {
		t.Fatal("unexpected error:", e)
	}
}

func TestAncestryRedefined(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	//
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
		&eph.Kinds{Kind: "n", Ancestor: "k"},
		&eph.Kinds{Kind: "q", Ancestor: "k"}, // should be okay; even though the domains differ, they're still compatible.
	)
	dt.MakeDomain(dd("c", "b"),
		&eph.Kinds{Kind: "m", Ancestor: "n"}, // should fail.
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict can't redefine the ancestor of "m" as "n"`); !ok {
		t.Fatal("unexpected error:", e)
	} else if e := warnings.Expect(`Duplicate "k" already declared as an ancestor of "q"`); e != nil {
		t.Fatal("unexpected warning:", e)
	}
}

// a similar named kind in another domain should be fine
func TestAncestryRivalsOkay(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
	)
	dt.MakeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // second in a parallel domain should be fine
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadKinds(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:", "b:m:k", "d:m:k",
	}); len(diff) > 0 {
		t.Fatal(diff)
	}
}

// two different kinds named in two different parent trees should fail
func TestAncestryRivalConflict(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"},
	)
	dt.MakeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "m", Ancestor: "k"}, // re: TestScopedRivalsOkay, should be okay.
	)
	dt.MakeDomain(dd("z", "b", "d")) // trying to include both should be a problem; they are two unique kinds...
	//
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict in domain`); !ok {
		t.Fatal(e)
	}
}
