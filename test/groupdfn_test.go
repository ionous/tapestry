package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

// adding a test struct as a kind should produce a certain set of fields and aspects
func TestKindsForType(t *testing.T) {
	var ks testutil.Kinds
	ks.AddKinds((*GroupCollation)(nil))
	// have to ask for the type to get the kinds to cache
	if _, e := ks.GetKindByName("group collation"); e != nil {
		t.Fail()
	} else if diff := pretty.Diff(ks.Builder, testutil.KindBuilder{
		Parents: nil,
		Aspects: []rt.Aspect{{
			Name: "group options",
			Traits: []string{
				"without objects",
				"objects with articles",
				"objects without articles",
			}},
		}, Fields: testutil.FieldMap{
			"group collation": {
				{Name: "groups", Affinity: "record_list", Type: "grouped objects"},
			},
			"grouped objects": {
				{Name: "settings", Affinity: "record", Type: "group settings"},
				{Name: "objects", Affinity: "text_list"},
			},
			"group settings": {
				{Name: "name", Affinity: "text"},
				{Name: "label", Affinity: "text"},
				{Name: "group options", Affinity: "text", Type: "group options"},
			},
		}}); len(diff) > 0 {
		t.Fatal(pretty.Println(ks.Builder))
	}
}
