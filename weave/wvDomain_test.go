package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave"

	"github.com/kr/pretty"
)

func TestDomainSimplest(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a", "b"))
	dt.MakeDomain(dd("b"))
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if as, e := dt.ReadDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(as, []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected ancestors:", pretty.Sprint(as))
		t.Fatal(diff)
	}
}

func TestDomainSimpleTest(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a", "b", "d"))
	dt.MakeDomain(dd("b", "c", "d"))
	dt.MakeDomain(dd("c", "d", "e"))
	dt.MakeDomain(dd("e", "d"))
	dt.MakeDomain(dd("d"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if as, e := dt.ReadDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(as, []string{"d", "e", "c", "b"}); len(diff) > 0 {
		// note: c requires d and e; but e requires d; so d is closest to the root, and g is root of all.
		t.Log("a has unexpected ancestors:", pretty.Sprint(as))
		t.Fatal(diff)
	} else if names, e := dt.ReadDomains(); e != nil {
		t.Fatal(e) // test getting the list of domains sorted from least to most dependent
	} else if diff := pretty.Diff(names, []string{"d", "e", "c", "b", "a"}); len(diff) > 0 {
		// d:1, e:2, c:3, b:4, a:5
		t.Log("domain names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

func TestDomainCatchCycles(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a", "b", "d"))
	dt.MakeDomain(dd("b", "c", "d"))
	dt.MakeDomain(dd("c", "d", "e"))
	dt.MakeDomain(dd("d", "a"))
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `circular reference`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

func TestDomainWhenUndeclared(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	// while we say "b" is a dependency of "a",
	// we never explicitly declare "b" --
	// and this should result in an error.
	dt.MakeDomain(dd("a", "b"))
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `circular or unknown domain`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

// various white spacing and casing should become more friendly underscore case
func TestDomainCase(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("alpha   domain", "beta domain"))
	dt.MakeDomain(dd("Beta Domain"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if ds, e := dt.ReadDomain("alpha domain"); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(ds, []string{"beta domain"}); len(diff) > 0 {
			t.Fatal(ds)
			t.Fatal(ds, diff)
		}
	}
}

func TestRivalStandalone(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"), rivalFact("secret"))
	dt.MakeDomain(dd("b"), rivalFact("mongoose"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	}
}

func TestRivalConflict(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"), rivalFact("secret"))
	dt.MakeDomain(dd("b"), rivalFact("mongoose"))
	dt.MakeDomain(dd("c", "a", "b"))
	//
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal("unexpected error:", e)
	}
}

// ephemera for testing which enters a definition
type rivalFact string

func (el rivalFact) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) error {
		return w.Pin().AddFact("rivalFact", string(el))
	})
}
