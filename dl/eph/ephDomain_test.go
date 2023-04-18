package eph

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestDomainSimplest(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b"))
	dt.makeDomain(dd("b"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if a, e := cat.resolveDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(namesOf(a.Ancestors()), []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected ancestors:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if diff := pretty.Diff(namesOf(a.Parents()), []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected parents:", pretty.Sprint(a))
		t.Fatal(diff)
	}
}

func TestDomainSimpleTest(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "d", "e"))
	dt.makeDomain(dd("e", "d"))
	dt.makeDomain(dd("d"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if a, e := cat.resolveDomain("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(namesOf(a.Ancestors()), []string{"d", "e", "c", "b"}); len(diff) > 0 {
		// note: c requires d and e; but e requires d; so d is closest to the root, and g is root of all.
		t.Log("a has unexpected ancestors:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if diff := pretty.Diff(namesOf(a.Parents()), []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected parents:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e) // test getting the list of domains sorted from least to most dependent
	} else {
		names := ds.Names()
		if diff := pretty.Diff(names, []string{"d", "e", "c", "b", "a"}); len(diff) > 0 {
			// d:1, e:2, c:3, b:4, a:5
			t.Log("domain names:", pretty.Sprint(names))
			t.Fatal(diff)
		}
	}
}

func TestDomainCatchCycles(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "d", "e"))
	dt.makeDomain(dd("d", "a"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if got, e := cat.resolveDomain(
		"a"); e == nil {
		t.Fatal(got) // we expected failure
	} else {
		t.Log("ok:", e)
	}
}

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func TestDomainTable(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "e"))
	dt.makeDomain(dd("d"))
	dt.makeDomain(dd("e"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{""} // write just the parents
		if e := ds.WriteTable(&out, "", false); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"d::x",
			"e::x",
			"c:e:x",
			"b:d,c:x", // fix? why does d wind up being listed before c? ( and in ancestors too )
			"a:b:x",
		}); len(diff) > 0 {
			t.Log("parents:", pretty.Sprint(out))
			t.Fatal(diff)
		} else {
			out := testOut{mdl.Domain}
			if e := cat.WriteDomains(&out); e != nil {
				t.Fatal(e)
			} else if diff := pretty.Diff(out[1:], testOut{
				"d::x",
				"e::x",
				"c:e:x",
				"b:d,c,e:x",
				"a:b,d,c,e:x",
			}); len(diff) > 0 {
				t.Log("ancestors:", pretty.Sprint(out))
				t.Fatal(diff)
			}
		}
	}
}
func TestDomainWhenUndeclared(t *testing.T) {
	var dt domainTest
	// while we say "b" is a dependency of "a",
	// we never explicitly declare "b" --
	// and this should result in an error.
	dt.makeDomain(dd("a", "b"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if ds, e := cat.ResolveDomains(); e == nil {
		t.Fatal("expected failure", ds)
	} else {
		t.Log("okay:", e)
	}
}

// various white spacing and casing should become more friendly underscore case
func TestDomainCase(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("alpha   domain", "beta domain"))
	dt.makeDomain(dd("BetaDomain"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if ds, e := cat.resolveDomain("alpha_domain"); e != nil {
		t.Fatal(e)
	} else {
		got := namesOf(ds.Ancestors())
		if diff := pretty.Diff(got, []string{"beta_domain"}); len(diff) > 0 {
			t.Fatal(got)
			t.Fatal(got, diff)
		}
	}
}

func TestRivalStandalone(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"), rivalFact("secret"))
	dt.makeDomain(dd("b"), rivalFact("mongoose"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(nil, nil); e != nil {
		t.Fatal(e)
	}
}

func TestRivalConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"), rivalFact("secret"))
	dt.makeDomain(dd("b"), rivalFact("mongoose"))
	dt.makeDomain(dd("c", "a", "b"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		var conflict *Conflict
		if e := cat.AssembleCatalog(nil, nil); !errors.As(e, &conflict) {
			t.Fatal("expected a conflict", e)
		} else if conflict.Reason != Redefined {
			t.Fatal("expected a redefinition error", e)
		} else {
			t.Log("ok", e)
		}
	}
}

func namesOf(ds []Dependency) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.Name()
	}
	return out
}

// ephemera for testing which enters a "
type rivalFact string

func (el rivalFact) Phase() assert.Phase { return assert.ValuePhase }

func (el rivalFact) Assemble(c *Catalog, d *Domain, at string) (err error) {
	return d.AddDefinition(MakeKey("rivalFact"), at, string(el))
}

func (c *Catalog) resolveDomain(n string) (ret Dependencies, err error) {
	if d, ok := c.GetDomain(n); !ok {
		err = errutil.New("unknown domain", n)
	} else {
		ret, err = d.Resolve()
	}
	return
}
