package rules

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"

	"git.sr.ht/~ionous/tapestry/rt"
)

func Test_SearchForCounters(t *testing.T) {
	c := &core.CallTrigger{}
	if !SearchForCounters(c) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{Test: []rt.BoolEval{c}}
		if !SearchForCounters(allTrue) {
			t.Fatal("all true")
		} else {
			not := &core.Not{Test: allTrue}
			if !SearchForCounters(not) {
				t.Fatal("not")
			} else {
				empty := &core.Not{Test: &core.AllTrue{}}
				if SearchForCounters(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
