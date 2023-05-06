package weave

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// a single pattern declaration containing all its parts
func TestPatternSingle(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Kinds{Kind: "k"},                      // a base for a parameter of k
		//
		&eph.Patterns{
			PatternName: "p",
			Params: []eph.Params{{
				Name:     "p1",
				Affinity: affine.Text,
				Class:    "k",
			}},
			Locals: []eph.Params{{
				Name:     "l1",
				Affinity: affine.NumList,
			}, {
				Name:      "l2",
				Affinity:  affine.Number,
				Initially: &assign.FromNumber{Value: I(10)},
			}},
			Result: &eph.Params{
				Name:     "success",
				Affinity: affine.Bool,
			},
		},
	)
	expectFullResults(t, dt)
}

// redeclare the same pattern as in single, but with multiple commands
func TestPatternSeparateLocals(t *testing.T) {
	// not sure yet if order of fields in the pattern be sorted at create time...
	// ( ex. return value first, last, always? )
	dt := newTestShuffle(t.Name(), false)
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Kinds{Kind: "k"},                      // a base for a parameter of k
		//
		&eph.Patterns{
			PatternName: "p",
			Result: &eph.Params{
				Name:     "success",
				Affinity: affine.Bool,
			}},
		&eph.Patterns{
			PatternName: "p",
			Params: []eph.Params{{
				Name:     "p1",
				Affinity: affine.Text,
				Class:    "k",
			}}},
		&eph.Patterns{
			PatternName: "p",
			Locals: []eph.Params{{
				Name: "l1",
				// FIX: rather than give an affinity
				// assign an empty assignment of the right type
				Affinity: affine.NumList,
			}}},
		&eph.Patterns{
			PatternName: "p",
			Locals: []eph.Params{{
				Name:      "l2",
				Affinity:  affine.Number,
				Initially: &assign.FromNumber{Value: I(10)},
			}}},
	)
	expectFullResults(t, dt)
}

// kinds should allow fields across domains, and so should locals.
func TestPatternSeparateDomains(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Kinds{Kind: "k"},                      // a base for a parameter of k
		//
		&eph.Patterns{
			PatternName: "p",
			Params: []eph.Params{{
				Name:     "p1",
				Affinity: affine.Text,
				Class:    "k",
			}}},
		&eph.Patterns{
			PatternName: "p",
			Result: &eph.Params{
				Name:     "success",
				Affinity: affine.Bool,
			}},
		&eph.Patterns{
			PatternName: "p",
			Locals: []eph.Params{{
				Name:     "l1",
				Affinity: affine.NumList,
			}}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Patterns{
			PatternName: "p",
			Locals: []eph.Params{{
				Name:      "l2",
				Affinity:  affine.Number,
				Initially: &assign.FromNumber{Value: I(10)},
			}}},
	)
	expectFullResults(t, dt)
}

func expectFullResults(t *testing.T, dt *domainTest) {
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if outkind, e := readKinds(cat.db); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outkind, []string{
		"a:k:",
		"a:patterns:",
		"a:p:patterns",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outkind))
		t.Fatal(diff)
	} else {
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
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Patterns{
			PatternName: "p",
			Result: &eph.Params{
				Name:     "success",
				Affinity: affine.Bool,
			}},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Patterns{
			PatternName: "p",
			Params: []eph.Params{{
				Name:     "p1",
				Affinity: affine.Text,
			}}},
	)

	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected an error")
	} else {
		t.Log("okay", e)
	}
}

// fail multiple returns
func TestPatternMultipleReturn(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()},
		//
		&eph.Patterns{
			PatternName: "p",
			Result: &eph.Params{
				Name:     "success",
				Affinity: affine.Bool,
			},
		},
		&eph.Patterns{
			PatternName: "p",
			Result: &eph.Params{
				Name:     "more success",
				Affinity: affine.Bool,
			},
		},
	)
	var conflict *Conflict

	_, e := dt.Assemble()
	if e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
		t.Fatal("expected an redefined conflict; got", e)
	} else {
		t.Log("okay", e)
	}
}

// fail multiple arg sets: args are now written individually so this is allowed.
// fix: stop args after writing locals or returns?
// func TestPatternMultipleArgSets(t *testing.T) {
// 	dt:= newTest(t.Name())
// defer dt.Close()
// 	dt.makeDomain(dd("a"),
// 		&eph.Kinds{Kind: kindsOf.Pattern.String()},
// 		&eph.Patterns{
// 			PatternName: "p",
// 			Params: []eph.Params{{
// 				Name:     "p1",
// 				Affinity: affine.Text},
// 			}},
// 		},
// 		&eph.Patterns{
// 			PatternName: "p",
// 			Params: []eph.Params{{
// 				Name:     "p2",
// 				Affinity: affine.Text},
// 			}},
// 		},
// 	)
// 	var conflict *Conflict
//
// if cat, e := dt.Assemble();  e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
// 		t.Fatal("expected an redefined conflict; got", e)
// 	} else {
// 		t.Log("okay", e)
// 	}
// }

// fail conflicting assignment
func TestPatternConflictingInit(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Patterns{
			PatternName: "p",
			Locals: []eph.Params{{
				Name:      "n",
				Affinity:  affine.Number,
				Initially: &assign.FromText{Value: T("mismatched")},
			}},
		},
	)
	if _, e := dt.Assemble(); e == nil || e.Error() != `mismatched affinity of initial value (a text) for field "n" (a number)` {
		t.Fatal("expected an error; got:", e)
	} else {
		t.Log("ok", e)
	}
}

// a simple pattern with no return
func TestPatternNoResults(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Patterns{PatternName: "p"},
	)
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if outkind, e := readKinds(cat.db); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outkind, []string{
		"a:patterns:",
		"a:p:patterns",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outkind))
		t.Fatal(diff)
	} else {
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
