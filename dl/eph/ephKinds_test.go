package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// test the kind mapping to kind list resolution
// we dont have to be exhaustive because its built on the dependency system which is already tested
func TestKinds(t *testing.T) {
	var ks Kinds
	for i, add := 0, []string{
		"device", "prop",
		"animal", "actor",
		"actor", "thing",
		"thing", "kind",
		"prop", "thing",
		"kind", "",
	}; i < len(add); i += 2 {
		ks.AddKind(add[i], add[i+1])
	}
	var res ResolvedKinds
	if e := ks.ResolveKinds(&res); e != nil {
		t.Fatal(e)
	} else {
		SortKinds(res)
		// [{kind []} {thing [kind]}
		if diff := pretty.Diff(res, ResolvedKinds{
			{"kind", nil},
			{"thing", dd("kind")},
			{"actor", dd("thing", "kind")},
			{"prop", dd("thing", "kind")},
			{"animal", dd("actor", "thing", "kind")},
			{"device", dd("prop", "thing", "kind")},
		}); len(diff) > 0 {
			t.Log(pretty.Print(res))
			t.Fatal(diff)
		}
	}

}
