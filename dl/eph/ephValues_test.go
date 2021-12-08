package eph

import (
	"testing"

	"github.com/kr/pretty"
)

func TestInitialFieldAssignment(t *testing.T) {
	var out []Ephemera
	addKinds(&out,
		"k", "",
		"l", "k",
		"m", "l",
	)
	addFields(&out,
		"k", "t", Affinity_Text, "",
		"k", "d", Affinity_Number, "",
	)
	addNouns(&out,
		"apple", "k",
		"pear", "l",
		"toy boat", "m",
		"boat", "m",
	)
	addValues(&out,
		"apple", "t", T("some text"),
		"pear", "d", I(123),
		"toy", "d", I(321),
		"boat", "t", T("more text"),
	)
	var cat Catalog
	var dt domainTest
	dt.makeDomain(dd("a"), out...)
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_val}
		if e := cat.WriteValues(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:apple:t:some text:x",
			"a:boat:t:more text:x",
			"a:toy_boat:d:321:x",
			"a:pear:d:123:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// func TestInitialTraitAssignment(t *testing.T) {
// 	var out []Ephemera
// 	addKinds(&out,
// 		"k", "",
// 		"l", "k",
// 		"m", "l",
// 	)
// 	addFields(&out,
// 		"k", "t", Affinity_Text, "",
// 		"k", "d", Affinity_Number, "",
// 	)
// 	addNouns(&out,
// 		"apple", "k",
// 		"pear", "l",
// 		"toy boat", "m",
// 		"boat", "m",
// 	)
// 	addValues(&out,
// 		"apple", "t", T("some text"),
// 		"pear", "d", I(123),
// 		"toy", "d", I(321),
// 		"boat", "t", T("more text"),
// 	)
// 	var cat Catalog
// 	var dt domainTest
// 	dt.makeDomain(dd("a"), out...)
// 	if e := dt.addToCat(&cat); e != nil {
// 		t.Fatal(e)
// 	} else if e := cat.AssembleCatalog(PhaseActions{
// 		AncestryPhase: AncestryPhaseActions,
// 		NounPhase:     NounPhaseActions,
// 	}); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		out := testOut{mdl_val}
// 		if e := cat.WriteValues(&out); e != nil {
// 			t.Fatal(e)
// 		} else if diff := pretty.Diff(out[1:], testOut{
// 			"a:apple:t:some text:x",
// 			"a:boat:t:more text:x",
// 			"a:toy_boat:d:321:x",
// 			"a:pear:d:123:x",
// 		}); len(diff) > 0 {
// 			t.Log(pretty.Sprint(out))
// 			t.Fatal(diff)
// 		}
// 	}
// }
