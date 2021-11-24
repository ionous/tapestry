package eph

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/kr/pretty"
)

// sorting resolved dependencies should work as expected.
func TestDomainSort(t *testing.T) {
	{
		og, tx := sorter("a", "b", "abc")
		if diff := pretty.Diff(tx, []string{"a", "b", "c"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("g", "gec", "gd", "ge", "gecdb", "gecdba")
		if diff := pretty.Diff(tx, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("g", "gdec", "gd", "gde", "gdecb", "gdecba")
		if diff := pretty.Diff(tx, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("t", "tu", "tuv", "a", "ab", "abc")
		if diff := pretty.Diff(tx, []string{"a", "b", "c", "t", "u", "v"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("a", "ab", "ac")
		if diff := pretty.Diff(tx, []string{"a", "b", "c"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
}

// single letter domains, root on the left, name on the right.
// ie. "agc" means "a" domain with ancestors "g" and "c"
func sorter(strs ...string) ([]string, []string) {
	ds := make(DependencyTable, len(strs))
	for i, str := range strs {
		row := make([]string, len(str))
		for j, el := range str {
			row[j] = string(el)
		}
		ds[i] = Dependents{ancestors: row}
	}
	rand.Shuffle(len(ds), func(i, j int) {
		ds[i], ds[j] = ds[j], ds[i]
	})
	was := ds.Names()
	ds.SortTable()
	return was, ds.Names()
}

func TestDomainSimplest(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b"))
	dt.makeDomain(dd("b"))
	var cat Catalog // the catalog adds a global "g" domain.
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if a, e := cat.GetDependentDomains("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(a.Ancestors(true), []string{"g", "b"}); len(diff) > 0 {
		t.Log("a has unexpected ancestors:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if diff := pretty.Diff(a.Ancestors(false), []string{"b"}); len(diff) > 0 {
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
	var cat Catalog // the catalog adds a global "g" domain.
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if a, e := cat.GetDependentDomains("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(a.Ancestors(true), []string{"g", "d", "e", "c", "b"}); len(diff) > 0 {
		// note: c requires d and e; but e requires d; so d is closest to the root, and g is root of all.
		t.Log("a has unexpected ancestors:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if diff := pretty.Diff(a.Ancestors(false), []string{"b"}); len(diff) > 0 {
		t.Log("a has unexpected parents:", pretty.Sprint(a))
		t.Fatal(diff)
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e) // test getting the list of domains sorted from least to most dependent
	} else {
		names := ds.Names()
		if diff := pretty.Diff(names, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			// g:0, d:1, e:2, c:3, b:4, a:5
			// ----- {"c", "d", "e", "b", "a", "g"}
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
	} else if got, e := cat.GetDependentDomains("a"); e == nil {
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
		var ancestors testOut
		if e := ds.WriteTable(&ancestors, mdl_domain, true); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(ancestors[mdl_domain], []string{
			"g:",
			"d:g",
			"e:g",
			"c:e,g",
			"b:d,c,e,g",
			"a:b,d,c,e,g",
		}); len(diff) > 0 {
			t.Log("ancestors:", pretty.Sprint(ancestors))
			t.Fatal(diff)
		} else {
			var parents testOut
			if e := ds.WriteTable(&parents, mdl_domain, false); e != nil {
				t.Fatal(e)
			} else if diff := pretty.Diff(parents[mdl_domain], []string{
				"g:",
				"d:g",
				"e:g",
				"c:e",
				"b:d,c", // fix? why does d wind up being listed before c? ( and in ancestors too )
				"a:b",
			}); len(diff) > 0 {
				t.Log("parents:", pretty.Sprint(parents))
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
	} else if ds, e := cat.GetDependentDomains("alpha_domain"); e != nil {
		t.Fatal(e)
	} else {
		got := ds.Ancestors(true)
		if diff := pretty.Diff(got, []string{"g", "beta_domain"}); len(diff) > 0 {
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
	} else if e := cat.ProcessDomains(nil); e != nil {
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
		if e := cat.ProcessDomains(nil); !errors.As(e, &conflict) {
			t.Fatal("expected a conflict", e)
		} else if conflict.Reason != Redefined {
			t.Fatal("expected a redefinition error", e)
		} else {
			t.Log("ok", e)
		}
	}
}

// ephemera for testing which enters a "
type rivalFact string

func (el rivalFact) Phase() Phase { return TestPhase }

func (el rivalFact) Catalog(c *Catalog, d *Domain, at string) (err error) {
	key, value := "rivalFact", string(el)
	return c.CheckConflict(d.name, "rivalFacts", at, key, value)
}
