package eph

import (
	"errors"
	"testing"

	"github.com/kr/pretty"
)

func TestDomainSimpleTest(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "d", "e"))
	dt.makeDomain(dd("e", "d"))
	dt.makeDomain(dd("d"))
	var cat Catalog // the catalog processing requires a global (root) domain.
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if got, e := cat.GetDependentDomains("a"); e != nil {
		t.Fatal(e) // test getting just the domains related to "a"
	} else if diff := pretty.Diff(got, []string{"b", "c", "d", "e", "g"}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(got))
		t.Fatal(diff)
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e) // test getting the list of domains sorted from least to most dependent
	} else {
		got := domainNames(ds)
		if diff := pretty.Diff(got, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			// g:0, d:0, e:2, c:3, b:4, a:5
			t.Log("got:", pretty.Sprint(got))
			t.Fatal(diff)
		}
	}
}

func domainNames(ds ResolvedDomains) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.name
	}
	return out
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
	if got, e := writeDomainTable(true); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []outEl{
		{"g", ""},
		{"d", "g"},
		{"e", "g"},
		{"c", "e,g"},
		{"b", "c,d,e,g"},
		{"a", "b,c,d,e,g"},
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(got))
		t.Fatal(diff)
	}
}

// same as table test, but the exclusive parents of each domain
func TestDomainParents(t *testing.T) {
	if got, e := writeDomainTable(false); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []outEl{
		{"g", ""},
		{"d", "g"},
		{"e", "g"},
		{"c", "e"},
		{"b", "c,d"},
		{"a", "b"},
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(got))
		t.Fatal(diff)
	}
}

// build a set of domains, and write them as if to a sql table
func writeDomainTable(fullTree bool) (ret []outEl, err error) {
	var dt domainTest
	dt.makeDomain(dd("a", "b", "d"))
	dt.makeDomain(dd("b", "c", "d"))
	dt.makeDomain(dd("c", "e"))
	dt.makeDomain(dd("d"))
	dt.makeDomain(dd("e"))
	var out testOut
	cat := Catalog{Writer: &out}
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if ds, e := cat.ResolveDomains(); e != nil {
		err = e
	} else if e := cat.WriteDomains(ds, fullTree); e != nil {
		err = e
	} else {
		// domain name and the table
		ret = out[mdl_domain]
	}
	return
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
	} else if got, e := cat.GetDependentDomains("alpha_domain"); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []string{"beta_domain", "g"}); len(diff) > 0 {
		t.Fatal(got, diff)
	}
}

func TestRivalStandalone(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"), rivalFact("secret"))
	dt.makeDomain(dd("b"), rivalFact("mongoose"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e)
	} else if e := ds.ProcessDomains(&cat); e != nil {
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
	} else if ds, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e)
	} else {
		var conflict *Conflict
		if e := ds.ProcessDomains(&cat); !errors.As(e, &conflict) {
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

func (el rivalFact) Phase() Phase { return Tests }

func (el rivalFact) Catalog(c *Catalog, d *Domain, at string) (err error) {
	key, value := "rivalFact", string(el)
	return c.CheckConflict(d.name, "rivalFacts", at, key, value)
}
