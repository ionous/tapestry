package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// test the kind mapping to kind list resolution
// probably didnt have to be exhaustive because its built on the dependency system which is already tested
func TestKindTree(t *testing.T) {
	ks := makeKinds(
		"a", "",
		"b", "a",
		"c", "b",
		"e", "b",
		"f", "e",
		"d", "c",
	)
	var res ResolvedKinds
	if e := ks.ResolveKinds(&res); e != nil {
		t.Fatal(e)
	} else {
		SortKinds(res)
		if diff := pretty.Diff(res, ResolvedKinds{
			{"a", nil},
			{"b", dd("a")},
			{"c", dd("a", "b")},
			{"e", dd("a", "b")},
			{"d", dd("a", "b", "c")},
			{"f", dd("a", "b", "e")},
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(res))
			t.Fatal(diff)
		}
	}
}

// this is considered okay - it's in the same tree
func TestKindDescendants(t *testing.T) {
	ks := makeKinds(
		"a", "",
		"b", "a",
		"c", "a",
		"c", "b",
	)
	var res ResolvedKinds
	if e := ks.ResolveKinds(&res); e != nil {
		t.Log(res)
		t.Fatal(e)
	} else {
		SortKinds(res)
		if diff := pretty.Diff(res, ResolvedKinds{
			{"a", nil},
			{"b", dd("a")},
			{"c", dd("a", "b")},
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(res))
			t.Fatal(diff)
		}
	}
}

func TestKindMissing(t *testing.T) {
	ks := makeKinds(
		"c", "d",
		"b", "a",
		"a", "",
	)
	var res ResolvedKinds
	if e := ks.ResolveKinds(&res); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok:", e)
	}
}

func TestKindConflict(t *testing.T) {
	ks := makeKinds(
		"a", "",
		"b", "a",
		"c", "a",
		"d", "b",
		"d", "c",
	)
	var res ResolvedKinds
	if e := ks.ResolveKinds(&res); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok:", e)
	}
}

func makeKinds(strs ...string) Kinds {
	var ks Kinds
	for i, cnt := 0, len(strs); i < cnt; i += 2 {
		ks.AddKind(strs[i], strs[i+1])
	}
	return ks
}
