package eph

import (
	"strconv"
	"testing"

	"github.com/kr/pretty"
)

// test the kind mapping to kind list resolution
// probably didnt have to be exhaustive because its built on the dependency system which is already tested
func TestKindTree(t *testing.T) {
	ks := makeKinds(t,
		"a", "",
		"b", "a",
		"c", "b",
		"e", "b",
		"f", "e",
		"d", "c",
	)
	if ks, e := ks.ResolveKinds(); e != nil {
		t.Fatal(e)
	} else {
		var out testOut
		if e := ks.WriteTable(&out, "", true); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out, testOut{
			"a:",
			"b:a",
			"c:b,a",
			"d:c,b,a",
			"e:b,a",
			"f:e,b,a",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// this is considered okay - it's in the same tree
func TestKindDescendants(t *testing.T) {
	ks := makeKinds(t,
		"a", "",
		"b", "a",
		"c", "a",
		"c", "b",
	)
	if res, e := ks.ResolveKinds(); e != nil {
		t.Log(res)
		t.Fatal(e)
	} else {
		var out testOut
		if e := res.WriteTable(&out, "", true); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out, testOut{
			"a:",
			"b:a",
			"c:b,a",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(res))
			t.Fatal(diff)
		}
	}
}

func TestKindMissing(t *testing.T) {
	ks := makeKinds(t,
		"c", "d",
		"b", "a",
		"a", "",
	)
	if res, e := ks.ResolveKinds(); e == nil {
		var out testOut
		res.WriteTable(&out, "", true)
		t.Fatal("expected error", out)
	} else {
		t.Log("ok:", e)
	}
}

func TestKindSingleParent(t *testing.T) {
	ks := makeKinds(t,
		"a", "",
		"b", "a",
		"c", "a",
		"d", "b",
		"d", "c",
	)
	if res, e := ks.ResolveKinds(); e == nil {
		var out testOut
		res.WriteTable(&out, "", true)
		t.Fatal("expected error", out)
	} else {
		t.Log("ok:", e)
	}
}

func makeKinds(t *testing.T, strs ...string) *Domain {
	d := Domain{name: "kinds", at: t.Name()}
	for i, cnt := 0, len(strs); i < cnt; i += 2 {
		a := d.EnsureKind(strs[i], strconv.Itoa(i))
		if b := strs[i+1]; len(b) > 0 {
			a.AddRequirement(b)
		}
	}
	return &d
}
