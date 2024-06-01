package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func TestCounterSearch(t *testing.T) {
	c := &call.CallTrigger{}
	if !rules.FilterHasCounter(c) {
		t.Fatal("core")
	} else {
		allTrue := &logic.IsAll{Test: []rt.BoolEval{c}}
		if !rules.FilterHasCounter(allTrue) {
			t.Fatal("all true")
		} else {
			not := &logic.Not{Test: allTrue}
			if !rules.FilterHasCounter(not) {
				t.Fatal("not")
			} else {
				empty := &logic.Not{Test: &logic.IsAll{}}
				if rules.FilterHasCounter(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
