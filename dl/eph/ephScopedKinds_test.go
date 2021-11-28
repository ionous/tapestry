package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// the basic TestKind(s) test kinds in a single domain
// so we want to to make sure we can handle multiple domains too
func TestSameScopedKinds(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "m", From: "k"}, // parent/root domain
	)
	dt.makeDomain(dd("c", "b"),
		&EphKinds{Kinds: "q", From: "j"}, // same domain
		&EphKinds{Kinds: "n", From: "k"}, // root domain
		&EphKinds{Kinds: "j", From: "m"}, // parent domain
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if ks, e := cat.ResolveDomains(); e != nil {
		t.Fatal(e)
	} else {
		var out testOut
		for _, deps := range ks {
			if e := cat.AssembleDomain(deps, PhaseActions{
				AncestryPhase: func(c *Catalog, d *Domain) (err error) {
					if ks, e := d.ResolveKinds(); e != nil {
						err = e
					} else if e := ks.WriteTable(&out, "", true); e != nil {
						err = e
					}
					return
				},
			}); e != nil {
				t.Fatal(e)
				break
			}
		}
		if diff := pretty.Diff(out, testOut{
			"k:", "m:k", "j:m,k", "q:j,m,k", "n:k",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
