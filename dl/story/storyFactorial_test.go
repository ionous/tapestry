package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/test/debug"
	"github.com/kr/pretty"
)

// read the factorial story
func TestFactorialImport(t *testing.T) {
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := k.ImportStory(t.Name(), debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else {
		// the hierarchical story as a flat list of commands used by the assembler
		// fix, future: "get var" and "assign" in the scope of a pattern should be generating parameter refs
		// fix, future: the debug story shouldnt be using core.CallArg, it should be using the story.Args
		//              so that we can see the code generate parameter references
		//              ( alt: there should be no story args, and marshal should do the work )
		//              cant change it right now b/c other tests depend on it.
		expect := []eph.Ephemera{
			// one test expectation
			&eph.EphChecks{
				Name:   "factorial",
				Expect: T("6"),
			},
			// one test rule
			&eph.EphChecks{
				Name: "factorial",
				Exe:  debug.FactorialCheck.Exe[0],
			},
			// a pattern definition including one parameter
			&eph.EphPatterns{
				Name: "factorial",
				Params: []eph.EphParams{{
					Affinity: eph.Affinity{eph.Affinity_Number},
					Name:     "num",
				}},
			},
			// a pattern return ( which happens to use the same input var )
			&eph.EphPatterns{
				Name: "factorial",
				Result: &eph.EphParams{
					Affinity: eph.Affinity{eph.Affinity_Number},
					Name:     "num",
				},
			},
			// our lowest priority rule
			&eph.EphRules{
				Name:   "factorial",
				Filter: &core.Always{},
				When:   eph.EphTiming{eph.EphTiming_During},
				Exe:    debug.FactorialMulMinusOne.Exe[0],
			},
			// the story happens to declare the return value twice
			// once before each rule.... that's fine it will be logged but it wont fail.
			&eph.EphPatterns{
				Name: "factorial",
				Result: &eph.EphParams{
					Affinity: eph.Affinity{eph.Affinity_Number},
					Name:     "num",
				},
			},
			// our highest priority rule ( that tests for zero )
			&eph.EphRules{
				Name:   "factorial",
				Filter: debug.FactorialIsZero,
				When:   eph.EphTiming{eph.EphTiming_During},
				Exe:    debug.FactorialUseOne.Exe[0],
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Error(diff, pretty.Sprint(els))
		}
	}
}
