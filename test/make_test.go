package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

// needs updating now that make is just a pattern call
func xTestMake(t *testing.T) {
	var kinds testutil.Kinds
	kinds.AddKinds((*GroupSettings)(nil))
	run := testpat.Runtime{
		testpat.Map{
			"group_settings": &Make,
		}, testutil.Runtime{
			Kinds: &kinds,
		},
	}
	op := &core.CallPattern{
		Pattern:   core.PatternName{Str: W("group_settings")},
		Arguments: core.NamedArgs("objects_with_articles", &core.FromBool{Val: B(true)}),
	}
	if obj, e := op.GetRecord(&run); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(g.RecordToValue(obj.Record()), map[string]interface{}{
		"name":          "",
		"label":         "",
		"innumerable":   "not_innumerable",
		"group_options": "objects_with_articles",
	}); len(diff) != 0 {
		t.Fatal(diff)
	}
}

var Make = testpat.Pattern{
	Name:   "group_settings",
	Labels: []string{"name", "label", "innumerable", "group_options"},
}
