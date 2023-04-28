package eph

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// a single pattern declaration containing all its parts
func TestPatternSingle(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphKinds{Kind: "k"},                      // a base for a parameter of k
		//
		&EphPatterns{
			PatternName: "p",
			Params: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}},
			Locals: []EphParams{{
				Name:     "l1",
				Affinity: Affinity{Affinity_NumList},
			}, {
				Name:      "l2",
				Affinity:  Affinity{Affinity_Number},
				Initially: &assign.FromNumber{Value: I(10)},
			}},
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
	)
	expectFullResults(t, &dt)
}

// redeclare the same pattern as in single, but with multiple commands
func TestPatternSeparateLocals(t *testing.T) {
	// not sure yet if order of fields in the pattern be sorted at create time...
	// ( ex. return value first, last, always? )
	dt := domainTest{noShuffle: true}
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphKinds{Kind: "k"},                      // a base for a parameter of k
		//
		&EphPatterns{
			PatternName: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
		&EphPatterns{
			PatternName: "p",
			Params: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}}},
		&EphPatterns{
			PatternName: "p",
			Locals: []EphParams{{
				Name: "l1",
				// FIX: rather than give an affinity
				// assign an empty assignment of the right type
				Affinity: Affinity{Affinity_NumList},
			}}},
		&EphPatterns{
			PatternName: "p",
			Locals: []EphParams{{
				Name:      "l2",
				Affinity:  Affinity{Affinity_Number},
				Initially: &assign.FromNumber{Value: I(10)},
			}}},
	)
	expectFullResults(t, &dt)
}

// kinds should allow fields across domains, and so should locals.
func TestPatternSeparateDomains(t *testing.T) {
	// not sure yet if order of fields in the pattern be sorted at create time...
	// ( ex. return value first, last, always? )
	dt := domainTest{noShuffle: true}
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphKinds{Kind: "k"},                      // a base for a parameter of k
		//
		&EphPatterns{
			PatternName: "p",
			Params: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
				Class:    "k",
			}}},
		&EphPatterns{
			PatternName: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
		&EphPatterns{
			PatternName: "p",
			Locals: []EphParams{{
				Name:     "l1",
				Affinity: Affinity{Affinity_NumList},
			}}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphPatterns{
			PatternName: "p",
			Locals: []EphParams{{
				Name:      "l2",
				Affinity:  Affinity{Affinity_Number},
				Initially: &assign.FromNumber{Value: I(10)},
			}}},
	)
	expectFullResults(t, &dt)
}

func expectFullResults(t *testing.T, dt *domainTest) {
	if cat, e := buildAncestors(dt); e != nil {
		t.Fatal(e)
	} else {
		outkind := testOut{mdl.Kind}
		if e := cat.WriteKinds(&outkind); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outkind[1:], testOut{
			"a:k::x",
			"a:patterns::x",
			"a:p:patterns:x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outkind))
			t.Fatal(diff)
		}
		outfields := testOut{mdl.Field}
		if e := cat.WriteFields(&outfields); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outfields[1:], testOut{
			"a:p:p_1:text:k:x",
			"a:p:success:bool::x",
			"a:p:l_1:num_list::x",
			"a:p:l_2:number::x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outfields))
			t.Fatal(diff)
		}
		outpat := testOut{mdl.Pat}
		if e := cat.WritePatterns(&outpat); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outpat[1:], testOut{
			// fix? run-in number suffixes?
			"a:p:p_1:success",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outpat))
			t.Fatal(diff)
		}
		outlocals := testOut{mdl.Assign}
		if e := cat.WriteLocals(&outlocals); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outlocals[1:], testOut{
			`a:p:l_2:{"FromNumber:":10}`,
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outlocals))
			t.Fatal(diff)
		}
	}
}

// fail splitting the args or returns across domains.
func TestPatternSplitDomain(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphPatterns{
			PatternName: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphPatterns{
			PatternName: "p",
			Params: []EphParams{{
				Name:     "p1",
				Affinity: Affinity{Affinity_Text},
			}}},
	)
	if _, e := buildAncestors(&dt); e == nil {
		t.Fatal("expected an error")
	} else {
		t.Log("okay", e)
	}
}

// fail multiple returns
func TestPatternMultipleReturn(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()},
		//
		&EphPatterns{
			PatternName: "p",
			Result: &EphParams{
				Name:     "success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
		&EphPatterns{
			PatternName: "p",
			Result: &EphParams{
				Name:     "more success",
				Affinity: Affinity{Affinity_Bool},
			},
		},
	)
	var conflict *Conflict
	if _, e := buildAncestors(&dt); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected an redefined conflict; got", e)
	} else {
		t.Log("okay", e)
	}
}

// fail multiple arg sets: args are now written individually so this is allowed.
// fix: stop args after writing locals or returns?
// func TestPatternMultipleArgSets(t *testing.T) {
// 	var dt domainTest
// defer dt.Close()
// 	dt.makeDomain(dd("a"),
// 		&EphKinds{Kind: kindsOf.Pattern.String()},
// 		&EphPatterns{
// 			PatternName: "p",
// 			Params: []EphParams{{
// 				Name:     "p1",
// 				Affinity: Affinity{Affinity_Text},
// 			}},
// 		},
// 		&EphPatterns{
// 			PatternName: "p",
// 			Params: []EphParams{{
// 				Name:     "p2",
// 				Affinity: Affinity{Affinity_Text},
// 			}},
// 		},
// 	)
// 	var conflict *Conflict
// 	if _, e := buildAncestors(&dt); e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
// 		t.Fatal("expected an redefined conflict; got", e)
// 	} else {
// 		t.Log("okay", e)
// 	}
// }

// fail conflicting assignment
func TestPatternConflictingInit(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphPatterns{
			PatternName: "p",
			Locals: []EphParams{{
				Name:      "n",
				Affinity:  Affinity{Affinity_Number},
				Initially: &assign.FromText{Value: T("mismatched")},
			}},
		},
	)
	if _, e := buildAncestors(&dt); e == nil || e.Error() != `mismatched affinity of initial value (a text) for field "n" (a number)` {
		t.Fatal("expected an error; got:", e)
	} else {
		t.Log("ok", e)
	}
}

// a simple pattern with no return
func TestPatternNoResults(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&EphPatterns{PatternName: "p"},
	)
	if cat, e := buildAncestors(&dt); e != nil {
		t.Fatal(e)
	} else {
		outkind := testOut{mdl.Kind}
		if e := cat.WriteKinds(&outkind); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outkind[1:], testOut{
			"a:patterns::x",
			"a:p:patterns:x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outkind))
			t.Fatal(diff)
		}
		outfields := testOut{mdl.Field}
		if e := cat.WriteFields(&outfields); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outfields[1:], testOut{
			// this might be nice, but would requiring changing pattern calls
			// *and* pattern tests ( which dont have first blank elements )
			//"a:p::bool::x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outfields))
			t.Fatal(diff)
		}
		outpat := testOut{mdl.Pat}
		if e := cat.WritePatterns(&outpat); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outpat[1:], testOut{
			"a:p::",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outpat))
			t.Fatal(diff)
		}
		outlocals := testOut{mdl.Assign}
		if e := cat.WriteLocals(&outlocals); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(outlocals[1:], testOut{
			/**/
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(outlocals))
			t.Fatal(diff)
		}
	}
}
