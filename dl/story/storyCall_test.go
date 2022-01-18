package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"

	"github.com/kr/pretty"
)

// verifies this expands a pattern call and that it generates a pattern reference.
func TestDetermineNum(t *testing.T) {
	var call rt.NumberEval
	if e := story.Decode(rt.NumberEval_Slot{&call}, []byte(`{"Factorial num:":{"FromNum:": 3}}`), tapestry.AllSignatures); e != nil {
		t.Fatal(e)
	} else {
		call := call.(*core.CallPattern)
		if diff := pretty.Diff(call, &core.CallPattern{
			Pattern: core.PatternName{Str: "factorial"},
			Arguments: core.CallArgs{
				Args: []core.CallArg{{
					Name: "num",
					From: &core.FromNum{
						Val: F(3),
					}}}}}); len(diff) > 0 {
			t.Fatal(diff)
		} else {
			refs := story.ImportPattern(call)
			if diff := pretty.Diff(refs, &eph.EphRefs{
				Refs: []eph.Ephemera{
					&eph.EphKinds{
						Kinds: "factorial",
						From:  kindsOf.Pattern.String(),
						Contain: []eph.EphParams{{
							Affinity: eph.Affinity{eph.Affinity_Number},
							Name:     "num",
						}},
					}},
			}); len(diff) > 0 {
				t.Fatal(diff)
			}
		}
	}
}
