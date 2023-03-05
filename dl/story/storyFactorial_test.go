package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/support/asm"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"
)

// read the factorial story
func TestFactorialImport(t *testing.T) {
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := asm.ImportStory(k, t.Name(), debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else {
		// the hierarchical story as a flat list of commands used by the assembler
		// fix, future: "get var" and "assign" in the scope of a pattern should be generating parameter refs
		expect := []eph.Ephemera{
			&eph.EphBeginDomain{
				Name: "factorial",
			},
			// referencing a call to the factorial pattern
			&eph.EphRefs{Refs: []eph.Ephemera{
				&eph.EphKinds{
					Kind: "factorial",
					// Ancestor:  kindsOf.Pattern.String() -- see note in ImportCall
					Contain: []eph.EphParams{{
						Affinity: eph.Affinity{eph.Affinity_Number},
						Name:     "num",
					}},
				},
			}},
			&eph.EphEndDomain{},
			// one test rule
			&eph.EphChecks{
				Name: "factorial",
				Exe:  debug.FactorialCheck,
			},
			// a pattern definition including one parameter
			&eph.EphPatterns{
				PatternName: "factorial",
				Params: []eph.EphParams{{
					Affinity: eph.Affinity{eph.Affinity_Number},
					Name:     "num",
				}},
				Result: &eph.EphParams{
					Affinity: eph.Affinity{eph.Affinity_Number},
					Name:     "num",
				},
			},
			// a pattern return ( which happens to use the same input var )
			// &eph.EphPatterns{
			// 	Name: "factorial",
			// },
			// our lowest priority rule
			&eph.EphRules{
				PatternName: "factorial",
				Filter:      &core.Always{},
				When:        eph.EphTiming{eph.EphTiming_During},
				Exe:         debug.FactorialMulMinusOne,
			},
			// the story happens to declare the return value twice
			// once before each rule.... that's fine it will be logged but it wont fail.
			// &eph.EphPatterns{
			// 	Name: "factorial",
			// 	Result: &eph.EphParams{
			// 		Affinity: eph.Affinity{eph.Affinity_Number},
			// 		Name:     "num",
			// 	},
			// },
			// our highest priority rule ( that tests for zero )
			&eph.EphRules{
				PatternName: "factorial",
				Filter:      debug.FactorialIsZero,
				When:        eph.EphTiming{eph.EphTiming_During},
				Exe:         debug.FactorialUseOne,
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Error(diff, pretty.Sprint(els))
		}
	}
}
