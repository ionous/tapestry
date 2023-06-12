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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readKinds(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:a:",
		"a:b:a",
		"a:c:b,a",
		"a:e:b,a",
		"a:d:c,b,a",
		"a:f:e,b,a",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// within a single domain,
// should be able to define a kind as being *more* specific than another definition.
func TestKindRefining(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"a", "",
		"b", "a",
		"c", "a",
		"c", "b",
	)...)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readKinds(); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(out, []string{
			"a:a:",
			"a:b:a",
			"a:c:b,a",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// this is considered okay - it's in the same tree
func TestKindDelayedRefining(t *testing.T) {
	dt := newTestShuffle(t.Name(), false)
	defer dt.Close()
	dt.makeDomain(dd("d"), makeKinds(
		"a", "",
		"b", "",
		"c", "b",
		"c", "a", // <-- if delayed refining weren't allowed this would fail
		"a", "b", // with "Conflict can't redefine the ancestor of "c" as "a""
	)...)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readKinds(); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(out, []string{
			"d:b:",
			"d:a:b",
			"d:c:b,a",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

func TestKindMissing(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"), makeKinds(
		"c", "d",
		"b", "a",
		"a", "",
	)...)
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Missing kind "d" in domain "a"`); !ok {
		t.Fatal("unexpected error:", e)
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
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Missing a definition in domain "a" that would allow "d" to have the ancestor`); !ok {
		t.Fatal("unexpected error:", e)
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
