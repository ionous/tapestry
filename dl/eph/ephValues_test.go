package eph

import (
	"testing"

	"github.com/kr/pretty"
)

func TestInitialFieldAssignment(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{"k", ""},
		&EphKinds{"l", "k"},
		&EphKinds{"m", "l"},
		// some simple fields
		// the name of the field has to match the name of the aspect
		&EphFields{"k", Affinity{Affinity_Text}, "t", ""},
		&EphFields{"k", Affinity{Affinity_Number}, "d", ""},
		// nouns with those fields
		&EphNouns{"apple", "k"},
		&EphNouns{"pear", "l"},
		&EphNouns{"toy boat", "m"},
		&EphNouns{"boat", "m"},
		// values using those fields
		&EphValues{"apple", "t", T("some text")},
		&EphValues{"pear", "d", I(123)},
		&EphValues{"toy", "d", I(321)},
		&EphValues{"boat", "t", T("more text")},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_value}
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

func TestInitialTraitAssignment(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{"k", ""},
		&EphKinds{"l", "k"},
		&EphKinds{"m", "l"},
		// aspects
		&EphKinds{KindsOfAspect, ""},
		&EphAspects{"a", dd("w", "x", "y")},
		&EphAspects{"b", dd("z")},
		// fields using those aspects:
		// the name of the field has to match the name of the aspect
		&EphFields{"k", Affinity{Affinity_Text}, "a", "a"},
		&EphFields{"k", Affinity{Affinity_Text}, "b", "b"},
		// nouns with those aspects
		&EphNouns{"apple", "k"},
		&EphNouns{"pear", "l"},
		&EphNouns{"toy boat", "m"},
		&EphNouns{"boat", "m"},
		// values using those aspects or traits from those aspects:
		&EphValues{"apple", "a", T("y")}, // assign to the aspect directly
		&EphValues{"pear", "x", B(true)}, // assign to some traits indirectly
		&EphValues{"toy", "w", B(true)},
		&EphValues{"boat", "z", B(true)},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_value}
		if e := cat.WriteValues(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:apple:a:y:x",
			"a:boat:b:z:x",
			"a:toy_boat:a:w:x",
			"a:pear:a:x:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
