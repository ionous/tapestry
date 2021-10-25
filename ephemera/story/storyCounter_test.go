package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/rt"
)

func Test_SearchForCounters(t *testing.T) {
	c := &core.CallTrigger{}
	if !story.SearchForCounters(c) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{[]rt.BoolEval{c}}
		if !story.SearchForCounters(allTrue) {
			t.Fatal("all true")
		} else {
			not := &core.Not{allTrue}
			if !story.SearchForCounters(not) {
				t.Fatal("not")
			} else {
				empty := &core.Not{&core.AllTrue{}}
				if story.SearchForCounters(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
