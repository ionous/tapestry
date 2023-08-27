package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/scope"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

// fix? i wonder if now "collation" can be the list of groups.
func TestGrouping(t *testing.T) {
	//errutil.Panic = true
	var kinds testutil.Kinds
	var objs testutil.Objects
	kinds.AddKinds((*Things)(nil), (*Locals)(nil), (*AssignGrouping)(nil), (*CollateGroups)(nil), (*MatchGroups)(nil))
	numObjects := objs.AddObjects(kinds.Kind("things"), "mildred", "apple", "pen", "thing-1", "thing-2")

	// create a new value of type "locals" containing "Objects:objectNames"
	// Objects   []string
	// Settings  []GroupSettings
	// Collation GroupCollation
	locals := kinds.NewRecord("locals", "objects", objs.Names())
	lt := testpat.Runtime{
		Runtime: testutil.Runtime{
			Kinds:     &kinds,
			ObjectMap: objs,
			Chain: scope.MakeChain(
				scope.FromRecord(locals),
			),
		},
		Map: testpat.Map{
			"assign grouping": &assignGrouping,
			"collate groups":  &collateGroups,
			"match groups":    &matchGroups,
		},
	}
	// generate one GroupSettings record for each object
	if e := runGroupTogther.Execute(&lt); e != nil {
		t.Fatal("groupTogther", e)
	} else {
		if settings, e := locals.GetNamedField("settings"); e != nil {
			t.Fatal(e)
		} else if cnt := settings.Len(); cnt != numObjects {
			t.Fatal("expected", numObjects, "settings, have", cnt)
		} else {
			// reduce the settings into unique groups
			if e := runCollateGroups.Execute(&lt); e != nil {
				t.Fatal("collateGroups", e)
			} else if collation, e := locals.GetNamedField("collation"); e != nil {
				t.Fatal(e)
			} else if groups, e := collation.FieldByName("groups"); e != nil {
				t.Fatal(e)
			} else {
				expect := []interface{}{
					map[string]interface{}{
						"settings": map[string]interface{}{
							"name":          "apple",
							"label":         "",
							"innumerable":   "not innumerable",
							"group options": "without objects",
						},
						"objects": []string{"apple", "mildred", "pen"},
					},
					map[string]interface{}{
						"settings": map[string]interface{}{
							"name":          "thing-1", // COUNTER:#
							"label":         "thingies",
							"innumerable":   "not innumerable",
							"group options": "without objects",
						},
						"objects": []string{"thing-1", "thing-2"}, // COUNTER:#
					},
				}
				got := g.RecordsToValue(groups.Records())
				if diff := pretty.Diff(expect, got); len(diff) > 0 {
					t.Log(pretty.Println(g.RecordsToValue(settings.Records())))
					t.Log(pretty.Sprint(got))
					t.Fatal(diff)
				}
			}
		}
	}
}
