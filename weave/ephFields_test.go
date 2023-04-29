package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// add some fields to a kind
func TestFields(t *testing.T) {
	dt := newTestShuffle(t.Name(), false) // fields arent sorted
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text, Class: "k"}}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	out := testOut{mdl.Field}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:t:text:k:x",
		"a:k:n:number::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can define a kind in one domain, and its fields in another
func TestFieldsCrossDomain(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close() // fields arent sorted, but are probably added in domain order so...
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "b", Affinity: affine.Bool}}},
	)
	out := testOut{mdl.Field}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:n:number::x",
		"a:k:b:bool::x",
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
	out := testOut{mdl.Field}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:n:number::x",
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
func TestFieldsMatchingRivals(t *testing.T) {
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
	out := testOut{mdl.Field}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:t:text::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
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
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text, Class: "m"}}},
	)
	dt.makeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "m"},
	)
	_, e := dt.Assemble()
	if e == nil || e.Error() != `unknown class "m" for field "t" for kind "k"` {
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
	out := testOut{mdl.Field}

	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:p:t:text::x",
		"a:q:t:text::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
