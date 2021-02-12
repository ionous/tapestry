package test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

// fix? i wonder if now "collation" can be the list of groups.
func TestGrouping(t *testing.T) {
	var kinds testutil.Kinds
	kinds.AddKinds((*Things)(nil), (*Values)(nil))
	objectNames := sliceOf.String("mildred", "apple", "pen", "thing_1", "thing_2")
	if objs, e := testutil.Objects(kinds.Kind("Things"), objectNames...); e != nil {
		t.Fatal(e)
	} else {
		// create a new value of type "Values" containing "Objects:objectNames"
		values := kinds.New("Values", "Objects", objectNames)
		lt := testutil.Runtime{
			Kinds:     &kinds,
			ObjectMap: objs,
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
			PatternMap: testutil.PatternMap{
				"assignGrouping": &assignGrouping,
				"collateGroups":  &collateGroups,
				"matchGroups":    &matchGroups,
			},
		}
		if e := runGroupTogther.Execute(&lt); e != nil {
			t.Fatal("groupTogther", e)
		} else if e := runCollateGroups.Execute(&lt); e != nil {
			t.Fatal("collateGroups", e)
		} else if collation, e := values.GetNamedField("Collation"); e != nil {
			t.Fatal(e)
		} else if groups, e := collation.FieldByName("Groups"); e != nil {
			t.Fatal(e)
		} else {
			expect := []interface{}{
				map[string]interface{}{
					"Settings": map[string]interface{}{
						"Name":         "mildred",
						"Label":        "",
						"Innumerable":  "NotInnumerable",
						"GroupOptions": "WithoutObjects",
					},
					"Objects": []string{"mildred", "apple", "pen"},
				},
				map[string]interface{}{
					"Settings": map[string]interface{}{
						"Name":         "thing_1", // COUNTER:#
						"Label":        "thingies",
						"Innumerable":  "NotInnumerable",
						"GroupOptions": "WithoutObjects",
					},
					"Objects": []string{"thing_1", "thing_2"}, // COUNTER:#
				},
			}
			got := g.RecordsToValue(groups.Records())
			if diff := pretty.Diff(expect, got); len(diff) > 0 {
				t.Log(pretty.Sprint(got))
				t.Fatal(diff)
			}
		}
	}
}

func logGroups(t *testing.T, groups []*g.Record) {
	t.Log("groups", len(groups), pretty.Sprint(g.RecordsToValue(groups)))
}
