package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/kr/pretty"
)

// add some fields to a kind
func TestFieldAssembly(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text, Class: "k"}}},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:n:number:", // field output gets sorted by name
		"a:k:t:text:k",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can define a kind in one domain, and its fields in another
func TestFieldsCrossDomain(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "b", Affinity: affine.Bool}}},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:b:bool:", // field output gets sorted by name
		"a:k:n:number:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can redefine fields in the same domain, and in another
func TestFieldsRedefine(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	//
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := warnings.Expect(`Duplicate field "n" for kind "k"`); e != nil {
		t.Fatal(e)
	} else if e := warnings.Expect(`Duplicate field "n" for kind "k"`); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadFields(); e != nil {
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
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Number}}},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "n", Affinity: affine.Text}}},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict field "n" for kind "k"`); !ok {
		t.Fatal("unexpected error:", e)
	}
}

// rival fields are fine so long as they match
// ( really the fields exist all at the same time )
func TestFieldsMatchingRivals(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()

	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.MakeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.MakeDomain(dd("z", "c", "d"))
	// fix: is this supposed to be an error?
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := warnings.Expect(`Duplicate field "t" for kind "k"`); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:k:t:text:",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fields in a given kind exist all at once; there's really not "rival" fields.
// this is really a name conflict.
func TestFieldsMismatchingRivals(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()

	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
	)
	dt.MakeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)
	dt.MakeDomain(dd("d", "a"),
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "t", Affinity: affine.Bool}}},
	)
	_, err := dt.Assemble()
	if ok, e := testweave.OkayError(t, err, `Conflict field "t" for kind "k"`); !ok {
		t.Fatal("unexpected error:", e)
	}
}

// classes cant refer to kinds that dont exist.
func TestFieldsUnknownClass(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			// m doesnt exist in this domain, so the test will fail.
			{Name: "t", Affinity: affine.Text, Class: "m"}}},
	)
	dt.MakeDomain(dd("c", "a"),
		&eph.Kinds{Kind: "m"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Missing kind "m" in domain "a" trying to find field "t"`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

// note: the original code would push shared fields upwards; the new code doesnt
func TestFieldLca(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "q", Ancestor: "t"},
		//
		&eph.Kinds{Kind: "p", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
		&eph.Kinds{Kind: "q", Contain: []eph.Params{{Name: "t", Affinity: affine.Text}}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:p:t:text:",
		"a:q:t:text:",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
