package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"github.com/kr/pretty"
)

// allow creating both sides of the opposites
func TestOppositeAllowed(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
		&eph.Opposites{Word: "west", Opposite: "east"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	}
}

// disallow mismatched opposites
func TestOppositeConflict(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
		&eph.Opposites{Word: "unkindness", Opposite: "east"},
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		// "east" had opposite "west" wanted "east" as "unkindness"
		t.Log("ok:", e)
	}
}

func TestOppositeAssembly(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
	)
	dt.MakeDomain(dd("b", "a"),
		&eph.Opposites{Word: "west", Opposite: "east"},
		&eph.Opposites{Word: "north", Opposite: "south"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadOpposites(); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(out, []string{
			"a:east:west",
			"a:west:east",
			"b:north:south",
			"b:south:north",
		}); len(diff) > 0 {
			t.Log("got", len(out), out)
			t.Fatal(diff)
		}
	}
}
