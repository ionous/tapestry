package story

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
)

func TestSearchForCounters(t *testing.T) {
	c := &core.CallTrigger{}
	if !searchForCounters(c) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{[]rt.BoolEval{c}}
		if !searchForCounters(allTrue) {
			t.Fatal("all true")
		} else {
			not := &core.Not{allTrue}
			if !searchForCounters(not) {
				t.Fatal("not")
			} else {
				empty := &core.Not{&core.AllTrue{}}
				if searchForCounters(empty) {
					t.Fatal("should have no counters")
				}
			}
		}
	}
}
