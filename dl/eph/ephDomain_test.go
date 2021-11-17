package eph

import (
	"testing"

	"github.com/kr/pretty"
)

func TestDomainSimpleTest(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("a", "b", "d"))
	dt.makeDomain(ds("b", "c", "d"))
	dt.makeDomain(ds("c", "d", "e"))
	dt.makeDomain(ds("e", "d"))
	var cat Catalog // the catalog processing requires a global (root) domain.
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		a := cat.GetDomain("a")
		if _, e := a.Resolve(); e != nil {
			t.Fatal(e)
		} else {
			got := a.resolved.names()
			if diff := pretty.Diff(got, []string{"b", "c", "d", "e", "g"}); len(diff) > 0 {
				t.Log("got:", pretty.Sprint(got))
				t.Fatal(diff)
			}
		}
	}
}

func TestDomainCatchCycles(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("a", "b", "d"))
	dt.makeDomain(ds("b", "c", "d"))
	dt.makeDomain(ds("c", "d", "e"))
	dt.makeDomain(ds("d", "a"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		a := cat.GetDomain("a")
		if _, e := a.Resolve(); e == nil {
			t.Fatal(a.resolved.names()) // we expected failure
		} else {
			t.Log("ok:", e)
		}
	}
}

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func TestDomainTable(t *testing.T) {
	if got, e := writeDomainTable(true); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []outEl{{
		// we expect that we start generating dependencies in alphabetical domain order
		// so: "a" gets evaluated first, recursively
		"a", "b,c,d,e,g",
	}, {
		"b", "c,d,e,g",
	}, {
		"c", "e,g",
	}, {
		"d", "g",
	}, {
		"e", "g",
	}, {
		"g", "",
	}}); len(diff) > 0 {
		t.Log(pretty.Sprint(got))
		t.Fatal(diff)
	}
}

// same as table test, but the exclusive parents of each domain
func TestDomainParents(t *testing.T) {
	if got, e := writeDomainTable(false); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []outEl{{
		"a", "b",
	}, {
		"b", "c,d",
	}, {
		"c", "e",
	}, {
		"d", "g",
	}, {
		"e", "g",
	}, {
		"g", "",
	}}); len(diff) > 0 {
		t.Log(pretty.Sprint(got))
		t.Fatal(diff)
	}
}

// build a set of domains, and write them as if to a sql table
func writeDomainTable(fullTree bool) (ret []outEl, err error) {
	var dt domainTest
	dt.makeDomain(ds("a", "b", "d"))
	dt.makeDomain(ds("b", "c", "d"))
	dt.makeDomain(ds("c", "e"))
	dt.makeDomain(ds("d"))
	dt.makeDomain(ds("e"))
	var out testOut
	cat := Catalog{Writer: &out}
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.WriteDomains(fullTree); e != nil {
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
	dt.makeDomain(ds("a", "b"))
	var out testOut
	cat := Catalog{Writer: &out}
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteDomains(true); e == nil {
		t.Fatal("expected failure", out)
	} else {
		t.Log("okay:", e)
	}
}

// various white spacing and casing should become more friendly underscore case
func TestDomainCase(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("alpha   domain", "beta domain"))
	dt.makeDomain(ds("BetaDomain"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		a := cat.GetDomain("alpha_domain")
		if _, e := a.Resolve(); e != nil {
			t.Fatal(e)
		} else {
			got := a.resolved.names()
			if diff := pretty.Diff(got, []string{"beta_domain", "g"}); len(diff) > 0 {
				t.Fatal(got, diff)
			}
		}
	}
}
