package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks testutil.Kinds
	ks.AddKinds((*GroupCollation)(nil))
	if diff := pretty.Diff(ks.Builder, testutil.KindBuilder{
		Aspects: testutil.AspectMap{
			"innumerable":   true,
			"group options": true,
		},
		Fields: testutil.FieldMap{
			"innumerable": {
				{"not innumerable", "bool", "" /*"trait"*/},
				{"is innumerable", "bool", "" /*"trait"*/},
			},
			"group options": {
				{"without objects", "bool", "" /*"trait"*/},
				{"objects with articles", "bool", "" /*"trait"*/},
				{"objects without articles", "bool", "" /*"trait"*/},
			},
			"group settings": {
				{"name", "text", ""},
				{"label", "text", ""},
				{"innumerable", "text", "innumerable"},
				{"group options", "text", "group options"},
			},
			"grouped objects": {
				{"settings", "record", "group settings"},
				{"objects", "text_list", ""},
			},
			"group collation": {
				{"groups", "record_list", "grouped objects"},
			},
		}}); len(diff) > 0 {
		t.Fatal(pretty.Println(ks.Builder))
	}
}
