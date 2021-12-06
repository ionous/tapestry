package eph

import (
	"errors"
	"testing"

	"github.com/kr/pretty"
)

// a single pattern declaration containing all its parts
func TestPatternSingle(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern}, // declare the patterns table
		&EphKinds{Kinds: "k"},            // a base for a parameter of k
		//
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}},
			Locals: []EphParams{{
				Name:     "l1",
				Affinity: Affinity{Affinity_NumList},
			}, {
				Name:     "l2",
				Affinity: Affinity{Affinity_Number},
			}},
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
	)
	expectFullResults(t, dt)
}

// redeclare the same pattern as in single, but with multiple commands
func TestPatternSeparateLocals(t *testing.T) {
	// not sure yet if order of fields in the pattern be sorted at create time...
	// ( ex. return value first, last, always? )
	dt := domainTest{noShuffle: true}
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern}, // declare the patterns table
		&EphKinds{Kinds: "k"},            // a base for a parameter of k
		//
		&EphPatterns{
			Name: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}}},
		&EphPatterns{
			Name: "p",
			Locals: []EphParams{{
				Name:     "l1",
				Affinity: Affinity{Affinity_NumList},
			}}},
		&EphPatterns{
			Name: "p",
			Locals: []EphParams{{
				Name:     "l2",
				Affinity: Affinity{Affinity_Number},
			}}},
	)
	expectFullResults(t, dt)
}

// kinds should allow fields across domains, and so should locals.
func TestPatternSeparateDomains(t *testing.T) {
	// not sure yet if order of fields in the pattern be sorted at create time...
	// ( ex. return value first, last, always? )
	dt := domainTest{noShuffle: true}
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern}, // declare the patterns table
		&EphKinds{Kinds: "k"},            // a base for a parameter of k
		//
		&EphPatterns{
			Name: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}}},
		&EphPatterns{
			Name: "p",
			Locals: []EphParams{{
				Name:     "l1",
				Affinity: Affinity{Affinity_NumList},
			}}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphPatterns{
			Name: "p",
			Locals: []EphParams{{
				Name:     "l2",
				Affinity: Affinity{Affinity_Number},
			}}},
	)
	expectFullResults(t, dt)
}

func expectFullResults(t *testing.T, dt domainTest) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		outkind := testOut{mdl_kind}
		if cat.WriteKinds(&outkind); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outkind[1:], testOut{
			"a:k::x",
			"a:pattern::x",
			"a:p:pattern:x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outkind))
			t.Fatal(diff)
		}
		outfields := testOut{mdl_field}
		if cat.WriteFields(&outfields); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outfields[1:], testOut{
			"a:p:success:bool::x",
			"a:p:p_1:text:k:x",
			"a:p:l_1:num_list::x",
			"a:p:l_2:number::x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outfields))
			t.Fatal(diff)
		}
		outpat := testOut{mdl_pat}
		if cat.WritePatterns(&outpat); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outpat[1:], testOut{
			// fix? run-in number suffixes?
			"a:p:p_1:success",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outpat))
			t.Fatal(diff)
		}
	}
}

// fail splitting the args or returns across domains.
func TestPatternSplitDomain(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern}, // declare the patterns table
		&EphPatterns{
			Name: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
			}}},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e == nil {
		t.Fatal("expected an error")
	} else {
		t.Log("okay", e)
	}
}

// fail multiple returns
func TestPatternMultipleReturn(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern},
		//
		&EphPatterns{
			Name: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
		&EphPatterns{
			Name: "p",
			Result: &EphParams{
				Name:     "more success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
	)
	var cat Catalog
	var conflict *Conflict
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected an redefined conflict; got", e)
	} else {
		t.Log("okay", e)
	}
}

// fail multiple arg sets
func TestPatternMultipleArgSets(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfPattern},
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
			}},
		},
		&EphPatterns{
			Name: "p",
			Args: []EphParams{{
				Name:     "p2",
				Affinity: Affinity{Affinity_Text},
			}},
		},
	)
	var cat Catalog
	var conflict *Conflict
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected an redefined conflict; got", e)
	} else {
		t.Log("okay", e)
	}
}
