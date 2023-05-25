package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// add some fields to a kind
func TestFields(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text, Class: "k"}}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:t:text:k",
		"a:k:n:number:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can define a kind in one domain, and its fields in another
func TestFieldsCrossDomain(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "b", Affinity: affine.Bool}}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:n:number:",
		"a:k:b:bool:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can redefine fields in the same domain, and in another
func TestFieldsRedefine(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	//
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:n:number:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fields conflict in sub-domains
// we can redefine fields in the same domain, and in another
func TestFieldsConflict(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Text}}},
	)
	_, e := dt.Assemble()
	if e := okDomainConflict("a", Redefined, e); e != nil {
		t.Fatal(e)
	} else {
		t.Log("ok:", e)
	}
}

// rival fields are fine so long as they match
// ( really the fields exist all at the same time )
//
// fix: this is failing -- checkRivals is allowDupes false for ancestry phase
// and its conflicting on kinds, k, parent name: "" -- from resolveKinds()
// not 100% why that worked before
func xxxTestFieldsMatchingRivals(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.makeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.makeDomain(dd("z", "c", "d"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:t:text:",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fields in kinds exist all at once --
// there's not really "rival" fields
func TestFieldsMismatchingRivals(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.makeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Bool}}},
	)
	// dt.makeDomain(dd("z", "c", "d")) <-- fails even without this.
	_, e := dt.Assemble()
	if e := okDomainConflict("a", Redefined, e); e != nil {
		t.Fatal(e)
	} else {
		t.Log("ok:", e)
	}
}

// classes cant refer to kinds that dont exist.
func TestFieldsUnknownClass(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			// m doesnt exist in this domain, so the test will fail.
			{Name: "t", Affinity: affine.Text, Class: "m"}}},
	)
	dt.makeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "m"},
	)
	_, e := dt.Assemble()
	if e == nil || e.Error() != `no such kind "m" in domain "a" trying to write field "t"` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok:", e)
	}
}

// note: the original code would push shared fields upwards; the new code doesnt
func TestFieldLca(t *testing.T) {
	dt := newTestShuffle(t.Name(), false) // fields arent sorted
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "q", Ancestor: "t"},
		//
		&eph.Kinds{Kind: "p", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
		&eph.Kinds{Kind: "q", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:p:t:text:",
		"a:q:t:text:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
