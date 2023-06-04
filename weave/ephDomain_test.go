package weave

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/assert"

	"github.com/kr/pretty"
)

func TestDomainSimplest(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a", "b"))
	dt.makeDomain(dd("b"))
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if as, e := dt.readDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(as, []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected ancestors:", pretty.Sprint(as))
		t.Fatal(diff)
	}
}

func TestDomainSimpleTest(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "d", "e"))
	dt.makeDomain(dd("e", "d"))
	dt.makeDomain(dd("d"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if as, e := dt.readDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(as, []string{"d", "e", "c", "b"}); len(diff) > 0 {
		// note: c requires d and e; but e requires d; so d is closest to the root, and g is root of all.
		t.Log("a has unexpected ancestors:", pretty.Sprint(as))
		t.Fatal(diff)
	} else if names, e := dt.readDomains(); e != nil {
		t.Fatal(e) // test getting the list of domains sorted from least to most dependent
	} else if diff := pretty.Diff(names, []string{"d", "e", "c", "b", "a"}); len(diff) > 0 {
		// d:1, e:2, c:3, b:4, a:5
		t.Log("domain names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

func TestDomainCatchCycles(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "d", "e"))
	dt.makeDomain(dd("d", "a"))
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `circular reference`); !ok {
		t.Fatal("expected error; got:", e)
	} else {
		t.Log("ok:", e)
	}
}

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func xxxxxTestDomainTable(t *testing.T) {
	// dt := newTestShuffle(t.Name(), false) // because of ids
	// defer dt.Close()
	// dt.makeDomain(dd("a", "b", "d"))
	// dt.makeDomain(dd("b", "c", "d"))
	// dt.makeDomain(dd("c", "e"))
	// dt.makeDomain(dd("d"))
	// dt.makeDomain(dd("e"))
	// if cat, e := dt.Assemble(); e != nil {
	// 	t.Fatal(e)
	// } else if out, e := dt.readDomains(); e != nil {
	// 	t.Fatal(e)
	// } else if diff := pretty.Diff(out, []string{
	// 		"d:",
	// 		"e:",
	// 		"c:e",
	// 		"b:d,c", // fix? why does d wind up being listed before c? ( and in ancestors too )
	// 		"a:b",
	// 	}); len(diff) > 0 {
	// 		t.Log("parents:", pretty.Sprint(out))
	// 		t.Fatal(diff)
	// 	} else {
	// 		if out, e := tables.QueryStrings(dt.db,
	// 			`select domain  || rowid || ':' || path from mdl_domain`); e != nil {
	// 			t.Fatal(e)
	// 		} else if diff := pretty.Diff(out, []string{
	// 			"a1:",
	// 			"b2:",
	// 			"c3:b2,",
	// 			"d4:a1,c3,b2,",
	// 			"e5:d4,a1,c3,b2,",
	// 		}); len(diff) > 0 {
	// 			t.Log("ancestors:", pretty.Sprint(out))
	// 			t.Fatal(diff)
	// 		}
	// 	}
}

func TestDomainWhenUndeclared(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	// while we say "b" is a dependency of "a",
	// we never explicitly declare "b" --
	// and this should result in an error.
	dt.makeDomain(dd("a", "b"))
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `circular or unknown domain`); !ok {
		t.Fatal("expected error; got:", e)
	} else {
		t.Log("ok:", e)
	}
}

// various white spacing and casing should become more friendly underscore case
func TestDomainCase(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("alpha   domain", "beta domain"))
	dt.makeDomain(dd("BetaDomain"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if ds, e := dt.readDomain("alpha_domain"); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(ds, []string{"beta_domain"}); len(diff) > 0 {
			t.Fatal(ds)
			t.Fatal(ds, diff)
		}
	}
}

func TestRivalStandalone(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), rivalFact("secret"))
	dt.makeDomain(dd("b"), rivalFact("mongoose"))

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	}
}

func TestRivalConflict(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), rivalFact("secret"))
	dt.makeDomain(dd("b"), rivalFact("mongoose"))
	dt.makeDomain(dd("c", "a", "b"))

	var conflict *Conflict
	if _, e := dt.Assemble(); !errors.As(e, &conflict) {
		t.Fatal("expected a conflict", e)
	} else if conflict.Reason != Redefined {
		t.Fatal("expected a redefinition error", e)
	} else {
		t.Log("ok:", e)
	}
}

func namesOf(ds []Dependency) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.Name()
	}
	return out
}

// ephemera for testing which enters a definition
type rivalFact string

func (el rivalFact) Phase() assert.Phase { return assert.ValuePhase }

func (el rivalFact) Assert(cat assert.Assertions) error {
	return cat.AssertDefinition("rivalFact", string(el))
}