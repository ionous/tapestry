package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if outkind, e := dt.readKinds(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outkind, []string{
		"a:k:",
		"a:patterns:",
		"a:p:patterns",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outkind))
		t.Fatal(diff)
	} else if outfields, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outfields, []string{
		"a:p:l_1:num_list:", // field output gets sorted by name
		"a:p:l_2:number:",
		"a:p:p_1:text:k",
		"a:p:success:bool:",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outfields))
		t.Fatal(diff)
	} else if outpat, e := dt.readPatterns(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outpat, []string{
		"a:p:p_1:success",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outpat))
		t.Fatal(diff)
	} else if outlocals, e := dt.readLocals(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outlocals, []string{
		`a:p:l_2:{"FromNumber:":10}`,
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outlocals))
		t.Fatal(diff)
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
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `unexpected result`); !ok {
		t.Fatal("expected error; got:", e)
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
// if _, e := dt.Assemble();  e == nil || !errors.As(e, &conflict) || conflict.Reason != Redefined {
// 		t.Fatal("expected an redefined conflict; got", e)
// 	} else {
// 		t.Log("okay", e)
// 	}
// }

// fail conflicting assignment:
// shouldnt be able to write text data into a number field
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
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict mismatched assignment for field "n" of kind "p"`); !ok {
		t.Fatal("expected error; got:", e)
	}
}

// a simple pattern with no parameters and no return
func TestPatternNoResults(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
		&eph.Patterns{PatternName: "p"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if outkind, e := dt.readKinds(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outkind, []string{
		"a:patterns:",
		"a:p:patterns",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outkind))
		t.Fatal(diff)
	} else if outfields, e := dt.readFields(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outfields, []string{
		// this might be nice, but would requiring changing pattern calls
		// *and* pattern tests ( which dont have first blank elements )
		//"a:p::bool:",
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outfields))
		t.Fatal(diff)
	} else if outpat, e := dt.readPatterns(); e != nil {
		t.Fatal(e)
	} else if len(outpat) > 0 {
		t.Fatal("expected no labels; got:", pretty.Sprint(outpat))
	} else if outlocals, e := dt.readLocals(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(outlocals, []string{
		/**/
	}); len(diff) > 0 {
		t.Log("got:", pretty.Sprint(outlocals))
		t.Fatal(diff)
	}
}
