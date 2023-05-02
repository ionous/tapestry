package weave

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables"
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
	if _, e := dt.Assemble(); e == nil {
		t.Fatal("expected failure")
	} else if !strings.Contains(e.Error(), "conflict") {
		t.Fatal(e)
	} else {
		// "east" had opposite "west" wanted "east" as "unkindness"
		t.Log("ok", e)
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
	} else if out, e := tables.ScanStrings(dt.db, `
select md.domain ||':'|| mp.oneWord ||':'|| mp.otherWord
from mdl_rev mp 
join mdl_domain md 
where md.rowid == mp.domain
order by md.domain, mp.oneWord`,
	); e != nil {
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
