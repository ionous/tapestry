package test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testpat"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

// fix? i wonder if now "collation" can be the list of groups.
func TestGrouping(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects
	kinds.AddKinds((*Things)(nil), (*Values)(nil), (*AssignGrouping)(nil), (*CollateGroups)(nil), (*MatchGroups)(nil))
	objs.AddObjects(kinds.Kind("things"), "mildred", "apple", "pen", "thing_1", "thing_2")

	// create a new value of type "values" containing "Objects:objectNames"
	values := kinds.NewRecord("values", "objects", objs.Names())
	lt := testpat.Runtime{
		Runtime: testutil.Runtime{
			Kinds:     &kinds,
			ObjectMap: objs,
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
		},
		Map: testpat.Map{
			"assign_grouping": &assignGrouping,
			"collate_groups":  &collateGroups,
			"match_groups":    &matchGroups,
		},
	}
	if e := runGroupTogther.Execute(&lt); e != nil {
		t.Fatal("groupTogther", e)
	} else if e := runCollateGroups.Execute(&lt); e != nil {
		t.Fatal("collateGroups", e)
	} else if collation, e := values.GetNamedField("collation"); e != nil {
		t.Fatal(e)
	} else if groups, e := collation.FieldByName("groups"); e != nil {
		t.Fatal(e)
	} else {
		expect := []interface{}{
			map[string]interface{}{
				"settings": map[string]interface{}{
					"name":          "apple",
					"label":         "",
					"innumerable":   "not_innumerable",
					"group_options": "without_objects",
				},
				"objects": []string{"apple", "mildred", "pen"},
			},
			map[string]interface{}{
				"settings": map[string]interface{}{
					"name":          "thing_1", // COUNTER:#
					"label":         "thingies",
					"innumerable":   "not_innumerable",
					"group_options": "without_objects",
				},
				"objects": []string{"thing_1", "thing_2"}, // COUNTER:#
			},
		}
		got := g.RecordsToValue(groups.Records())
		if diff := pretty.Diff(expect, got); len(diff) > 0 {
			t.Log(pretty.Sprint(got))
			t.Fatal(diff)
		}
	}
}

func logGroups(t *testing.T, groups []*g.Record) {
	t.Log("groups", len(groups), pretty.Sprint(g.RecordsToValue(groups)))
}
