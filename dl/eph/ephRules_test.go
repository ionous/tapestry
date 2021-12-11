package eph

import (
	"strconv"
	"testing"

	"github.com/kr/pretty"
)

func TestPatternRules(t *testing.T) {
	dt := domainTest{noShuffle: true}
	dt.makeDomain(dd("a"),
		append(append([]Ephemera{
			&EphKinds{Kinds: KindsOfPattern}, // declare the patterns table
			&EphPatterns{
				Name: "p",
				Result: &EphParams{
					Name:     "success",
					Affinity: Affinity{Affinity_Bool},
				}},
		}, makeRules("p", "d", EphTiming_During, 3)...),
			makeRules("p", "b", EphTiming_Before, 2)...)...,
	)
	dt.makeDomain(dd("b", "a"),
		makeRules("p", "a", EphTiming_After, 3)...,
	)
	//
	if cat, e := buildAncestors(dt); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_rule}
		if e := cat.WriteRules(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			// domain, pattern, phase, filter, prog, at
			"b:p:3:filter-a2:prog-a2:x", // for now at least, we list domains backwards for rules.
			"b:p:3:filter-a1:prog-a1:x",
			"b:p:3:filter-a0:prog-a0:x",
			//
			"a:p:1:filter-b1:prog-b1:x", // even though the "b" items were specified second, "before" rules get listed first.
			"a:p:1:filter-b0:prog-b0:x",
			//
			"a:p:2:filter-d2:prog-d2:x",
			"a:p:2:filter-d1:prog-d1:x",
			"a:p:2:filter-d0:prog-d0:x",
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

func makeRules(pattern, group, timing string, cnt int) (ret []Ephemera) {
	for i := 0; i < cnt; i++ {
		ret = append(ret, &EphRules{
			Name:   pattern,
			Filter: "filter-" + group + strconv.Itoa(i),
			When:   EphTiming{timing},
			Prog:   "prog-" + group + strconv.Itoa(i),
		})
	}
	return
}
