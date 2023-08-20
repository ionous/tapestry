package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func TestCounterSearch(t *testing.T) {
	c := &core.CallTrigger{}
	if !rules.FilterHasCounter(c) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{Test: []rt.BoolEval{c}}
		if !rules.FilterHasCounter(allTrue) {
			t.Fatal("all true")
		} else {
			not := &core.Not{Test: allTrue}
			if !rules.FilterHasCounter(not) {
				t.Fatal("not")
			} else {
				empty := &core.Not{Test: &core.AllTrue{}}
				if rules.FilterHasCounter(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
