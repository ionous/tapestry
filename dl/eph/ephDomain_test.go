package eph

import (
	"strconv"
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
		if e := a.Resolve(nil); e != nil {
			t.Fatal(e)
		} else {
			names := a.Resolved()
			if diff := pretty.Diff(names, []string{"b", "c", "d", "e", "g"}); len(diff) > 0 {
				t.Fatal(diff, names)
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
		if e := a.Resolve(nil); e == nil {
			t.Fatal(a.Resolved()) // we expect failure
		} else {
			t.Log("ok:", e)
		}
	}
}

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func TestDomainTable(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("a", "b", "d"))
	dt.makeDomain(ds("b", "c", "d"))
	dt.makeDomain(ds("c", "e"))
	dt.makeDomain(ds("d"))
	dt.makeDomain(ds("e"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		var out testOut
		if e := cat.WriteDomains(&out); e != nil {
			t.Fatal(e)
		} else {
			// domain name and the table
			got := out[mdl_domain]
			// we expect that we start generating dependencies in alphabetical domain order
			// so: "a" gets evaluated first, recursively
			if diff := pretty.Diff(got, []outEl{{
				// everything depends on g, so we see that first
				"g", "",
			}, {
				// a's first dep is "b", b's is "c", c's is "e"... so we see "e" next.
				// it only depends on g.
				"e", "g",
			}, {
				// c has no other dependencies other than e and g
				"c", "e,g",
			}, {
				// back in "b"'s dependencies, it's also dependent on "d"
				// and "d" -- like everything else -- depends on the global "g"
				"d", "g",
			}, {
				// we're done with the dependencies of b's dependencies
				// everything b needs, has been written into the table at this point
				"b", "c,d,e,g",
			}, {
				// and, likewise that's true for "a"
				"a", "b,c,d,e,g",
			}}); len(diff) > 0 {
				t.Log(pretty.Sprint(got))
				// t.Fatal(got, diff)
			}
		}
	}
}

func TestDomainWhenUndeclared(t *testing.T) {
	var dt domainTest
	// while we say "b" is a dependency of "a",
	// we never explicitly declare "b" --
	// and this should result in an error.
	dt.makeDomain(ds("a", "b"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		var out testOut
		if e := cat.WriteDomains(&out); e == nil {
			t.Fatal("expected failure", out)
		} else {
			t.Log("okay:", e)
		}
	}
}

// various white spacing and casing should become more friendly underscore case1
func TestDomainCase(t *testing.T) {
	var dt domainTest
	dt.makeDomain(ds("alpha   domain", "beta domain"))
	dt.makeDomain(ds("BetaDomain"))
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else {
		a := cat.GetDomain("alpha_domain")
		if e := a.Resolve(nil); e != nil {
			t.Fatal(e)
		} else {
			got, want := a.Resolved(), []string{"beta_domain", "g"}
			if diff := pretty.Diff(got, want); len(diff) > 0 {
				t.Fatal(got, diff)
			}
		}
	}
}

type domainTest struct {
	out []Ephemera
}

func ds(names ...string) []string {
	return names
}

func (dt *domainTest) makeDomain(names []string, add ...Ephemera) {
	dt.out = append(dt.out, &EphBeginDomain{
		Name:     names[0],
		Requires: names[1:],
	})
	dt.out = append(dt.out, add...)
	dt.out = append(dt.out, &EphEndDomain{
		Name: names[0],
	})
	return
}

func (dt *domainTest) addToCat(cat *Catalog) (err error) {
	g := cat.GetDomain("g")
	g.at = "global"
	cat.processing.Push(g)
	for i, el := range dt.out {
		if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			err = e
			break
		}
	}
	return
}
