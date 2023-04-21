package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"

	"github.com/kr/pretty"
)

// verifies this expands a pattern call and that it generates a pattern reference.
func TestDetermineNum(t *testing.T) {
	var call rt.NumberEval
	if e := story.Decode(rt.NumberEval_Slot{&call}, []byte(`{"Factorial num:":{"FromNumber:": 3}}`), tapestry.AllSignatures); e != nil {
		t.Fatal(e)
	} else {
		call := call.(*assign.CallPattern)
		if diff := pretty.Diff(call, &assign.CallPattern{
			PatternName: "factorial",
			Arguments: []assign.Arg{{
				Name:  "num",
				Value: &assign.FromNumber{Value: F(3)},
			}}}); len(diff) > 0 {
			t.Fatal(diff)
		} else {
			// disabling refs for now....
			// maybe instead could just request from the db that something exists in a scheduled post-weave check.
			// refs := story.ImportCall(call)
			// if diff := pretty.Diff(refs, &eph.EphRefs{
			// 	Refs: []eph.Ephemera{
			// 		&eph.EphKinds{
			// 			Kind: "factorial",
			// 			// From:  kindsOf.Pattern.String() -- see note in ImportCall
			// 			Contain: []eph.EphParams{{
			// 				Affinity: eph.Affinity{eph.Affinity_Number},
			// 				Name:     "num",
			// 			}},
			// 		}},
			// }); len(diff) > 0 {
			// 	t.Fatal(diff)
			// }
		}
	}
}
