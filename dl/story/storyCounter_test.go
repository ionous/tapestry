package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"

	"git.sr.ht/~ionous/tapestry/rt"
)

func Test_SearchForCounters(t *testing.T) {
	c := &core.CallTrigger{}
	if !story.SearchForCounters(c) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{Test: []rt.BoolEval{c}}
		if !story.SearchForCounters(allTrue) {
			t.Fatal("all true")
		} else {
			not := &core.Not{Test: allTrue}
			if !story.SearchForCounters(not) {
				t.Fatal("not")
			} else {
				empty := &core.Not{Test: &core.AllTrue{}}
				if story.SearchForCounters(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
