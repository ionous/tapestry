package story_test

import (
	"testing"
)

// read the factorial story
// todo: verify that the db gets the expected elements
func Fix_TestFactorialImport(t *testing.T) {
	// db := testdb.Open(t.Name(), testdb.Memory, "")
	// k := weave.NewCatalog(db)
	// if e := story.ImportStory(k, t.Name(), debug.FactorialStory); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	// the hierarchical story as a flat list of commands used by the assembler
	// 	// fix, future: "get var" and "assign" in the scope of a pattern should be generating parameter refs
	// 	expect := []eph.Ephemera{
	// 		&eph.AssertDomainStart{
	// 			Name: "factorial",
	// 		},
	// 		// fix? disabled refs for now.
	// 		// referencing a call to the factorial pattern
	// 		// &eph.Refs{Refs: []eph.Ephemera{
	// 		// 	&eph.Kinds{
	// 		// 		Kind: "factorial",
	// 		// 		// Ancestor:  kindsOf.Pattern.String() -- see note in ImportCall
	// 		// 		Contain: []eph.Params{{
	// 		// 			Affinity: eph.affine.Number},
	// 		// 			Name:     "num",
	// 		// 		}},
	// 		// 	},
	// 		// }},
	// 		&eph.AssertDomainEnd{
	// 			Name: "factorial",
	// 		},
	// 		// one test rule
	// 		&eph.Checks{
	// 			Name: "factorial",
	// 			Exe:  debug.FactorialCheck,
	// 		},
	// 		// a pattern definition including one parameter
	// 		&eph.Patterns{
	// 			PatternName: "factorial",
	// 		},
	// 		&eph.Patterns{
	// 			PatternName: "factorial",
	// 			Result: &eph.Params{
	// 				Affinity: eph.affine.Number},
	// 				Name:     "num",
	// 			},
	// 		},
	// 		&eph.Patterns{
	// 			PatternName: "factorial",
	// 			Params: []eph.Params{{
	// 				Affinity: eph.affine.Number},
	// 				Name:     "num",
	// 			}},
	// 		},
	// 		// a pattern return ( which happens to use the same input var )
	// 		// &eph.Patterns{
	// 		// 	Name: "factorial",
	// 		// },
	// 		// our lowest priority rule
	// 		&eph.Rules{
	// 			PatternName: "factorial",
	// 			Filter:      &core.Always{},
	// 			When:        eph.Timing{eph.Timing_During},
	// 			Exe:         debug.FactorialMulMinusOne,
	// 		},
	// 		// the story happens to declare the return value twice
	// 		// once before each rule.... that's fine it will be logged but it wont fail.
	// 		// &eph.Patterns{
	// 		// 	Name: "factorial",
	// 		// 	Result: &eph.Params{
	// 		// 		Affinity: eph.affine.Number},
	// 		// 		Name:     "num",
	// 		// 	},
	// 		// },
	// 		// our highest priority rule ( that tests for zero )
	// 		&eph.Rules{
	// 			PatternName: "factorial",
	// 			Filter:      debug.FactorialIsZero,
	// 			When:        eph.Timing{eph.Timing_During},
	// 			Exe:         debug.FactorialUseOne,
	// 		},
	// 	}
	// 	if diff := pretty.Diff(els, expect); len(diff) > 0 {
	// 		t.Error(diff, pretty.Sprint(els))
	// 	}
	// }
}
