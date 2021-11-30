package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// 1. add some fields to a kind
// 2. we can define a kind in one domain, and its fields in another
// 3. we can redefine fields in the same domain, and in another
// 4. fields conflict in sub-domains
// 5. rivals declarations are fine so long as they match
// 6. classes cant refer to kinds that dont exist.
// -----------------------------------------------
func TestFields(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "t", Class: "k", Affinity: Affinity{Affinity_Text}},
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	var out testOut
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"a:k:t:text:k:2",
		"a:k:n:number::3",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func buildFields(dt domainTest) (ret *Catalog, err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: PhaseAction{
			PhaseFlags{NoDuplicates: true},
			func(d *Domain) error {
				_, e := d.ResolveKinds()
				return e
			},
		},
	}); e != nil {
		err = e
	} else {
		ret = &cat
	}
	return
}
