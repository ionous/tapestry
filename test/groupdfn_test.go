package test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks testutil.Kinds
	ks.AddKinds((*GroupCollation)(nil))
	if diff := pretty.Diff(ks.Fields, testutil.FieldMap{
		"innumerable": {
			{"not_innumerable", "bool", "trait"},
			{"is_innumerable", "bool", "trait"},
		},
		"group_options": {
			{"without_objects", "bool", "trait"},
			{"objects_with_articles", "bool", "trait"},
			{"objects_without_articles", "bool", "trait"},
		},
		"group_settings": {
			{"name", "text", "string"},
			{"label", "text", "string"},
			{"innumerable", "text", "aspect"},
			{"group_options", "text", "aspect"},
		},
		"grouped_objects": {
			{"settings", "record", "group_settings"},
			{"objects", "text_list", "string"},
		},
		"group_collation": {
			{"groups", "record_list", "grouped_objects"},
		},
	}); len(diff) > 0 {
		t.Fatal(pretty.Println(ks.Fields))
	}
}
