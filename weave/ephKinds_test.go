package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// test the kind mapping to kind list resolution
func TestKindTree(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"a", "",
		"b", "a",
		"c", "b",
		"e", "b",
		"f", "e",
		"d", "c",
	)...)
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if ks, e := cat.ResolveKinds(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{""}
		if e := ks.WriteTable(&out, "", true); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a::x",
			"b:a:x",
			"c:b,a:x",
			"d:c,b,a:x",
			"e:b,a:x",
			"f:e,b,a:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// this is considered okay - it's in the same tree
func TestKindDescendants(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"a", "",
		"b", "a",
		"c", "a",
		"c", "b",
	)...)
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if ks, e := cat.ResolveKinds(); e != nil {
		t.Fatal(e)
	} else {

		out := testOut{""}
		if e := ks.WriteTable(&out, "", true); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a::x",
			"b:a:x",
			"c:b,a:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// FIX: disabled; uses the catalog runtime now when finding plural fallbacks
// ( which doesnt exist here, causing a panic )
func xTestKindMissing(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"c", "d",
		"b", "a",
		"a", "",
	)...)
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok:", e)
	}
}

func TestKindSingleParent(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"a", "",
		"b", "a",
		"c", "a",
		"d", "b",
		"d", "c",
	)...)
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok:", e)
	}
}

func makeKinds(strs ...string) []eph.Ephemera {
	kinds := make([]eph.Ephemera, len(strs)/2)
	for i := range kinds {
		kinds[i] = &eph.Kinds{
			Kind:     strs[i*2+0],
			Ancestor: strs[i*2+1],
		}
	}
	return kinds
}
