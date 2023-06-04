package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// allow creating both sides of the opposites
func TestOppositeAllowed(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
		&eph.Opposites{Word: "west", Opposite: "east"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	}
}

// disallow mismatched opposites
func TestOppositeConflict(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
		&eph.Opposites{Word: "unkindness", Opposite: "east"},
	)
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `conflict`); !ok {
		t.Fatal("expected error; got:", e)
	} else {
		// "east" had opposite "west" wanted "east" as "unkindness"
		t.Log("ok:", e)
	}
}

func TestOppositeAssembly(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Opposites{Word: "east", Opposite: "west"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Opposites{Word: "west", Opposite: "east"},
		&eph.Opposites{Word: "north", Opposite: "south"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readOpposites(); e != nil {
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