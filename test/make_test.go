package test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestMake(t *testing.T) {
	type panicTime struct {
		testutil.PanicRuntime
	}
	var testTime struct {
		panicTime
		testutil.Kinds
	}
	testTime.Kinds.AddKinds((*GroupSettings)(nil))
	op := &core.Make{Name: "group_settings",
		Arguments: &core.Arguments{[]*core.Argument{
			{"objects_with_articles", &core.FromBool{&core.Bool{true}}},
		}}}
	if obj, e := op.GetRecord(&testTime); e != nil {
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
