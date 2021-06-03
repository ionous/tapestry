package story

import (
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
)

func TestSearchForCounters(t *testing.T) {
	c := &core.CallTrigger{}
	if !searchForCounters(r.ValueOf(c)) {
		t.Fatal("core")
	} else {
		allTrue := &core.AllTrue{[]rt.BoolEval{c}}
		if !searchForCounters(r.ValueOf(allTrue)) {
			t.Fatal("all true")
		} else {
			not := &core.Not{allTrue}
			if !searchForCounters(r.ValueOf(not)) {
				t.Fatal("not")
			}
		}
	}
}
