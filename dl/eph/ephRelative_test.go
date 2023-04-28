package eph

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// follow along with relative test except add list of ephemera
func TestRelativeFormation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			"a", "r_1_1", "a",
			"b", "r_1_1", "c",
			"c", "r_1_1", "b",
			"z", "r_1_1", "e",
			// these pass the original test, but are in fact failures.
			//	"b", "r_1_1", "d",
			//	"c", "r_1_1", "a",
			//	"z", "r_1_1", "f",
			//
			"z", "r_1_x", "f",
			"c", "r_1_x", "a",
			"b", "r_1_x", "e",
			"c", "r_1_x", "c",
			"c", "r_1_x", "b",
			"z", "r_1_x", "d",
			"z", "r_1_x", "f",
			//
			"z", "r_x_1", "b",
			"f", "r_x_1", "f",
			"l", "r_x_1", "b",
			"b", "r_x_1", "a",
			"d", "r_x_1", "b",
			"c", "r_x_1", "d",
			"f", "r_x_1", "f",
			"e", "r_x_1", "f",
			//
			"a", "r_x_x", "a",
			"e", "r_x_x", "d",
			"a", "r_x_x", "b",
			"a", "r_x_x", "c",
			"f", "r_x_x", "d",
			"l", "r_x_x", "d",
			"a", "r_x_x", "b",
		)...,
	)
	if cat, e := buildNouns(&dt); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Pair}
		if e := cat.WritePairs(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:r_1_1:a:a:x",
			"a:r_1_1:b:c:x",
			"a:r_1_1:e:z:x",
			//
			"a:r_1_x:b:e:x",
			"a:r_1_x:c:a:x",
			"a:r_1_x:c:b:x",
			"a:r_1_x:c:c:x",
			"a:r_1_x:z:d:x",
			"a:r_1_x:z:f:x",
			//
			"a:r_x_1:b:a:x",
			"a:r_x_1:c:d:x",
			"a:r_x_1:d:b:x",
			"a:r_x_1:e:f:x",
			"a:r_x_1:f:f:x",
			"a:r_x_1:l:b:x",
			"a:r_x_1:z:b:x",
			//
			"a:r_x_x:a:a:x",
			"a:r_x_x:a:b:x",
			"a:r_x_x:a:c:x",
			"a:r_x_x:e:d:x",
			"a:r_x_x:f:d:x",
			"a:r_x_x:l:d:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// follow along with relative test except add list of ephemera
func TestRelativeOneOneViolation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			"b", "r_1_1", "c",
			"b", "r_1_1", "d",
		)...,
	)
	if _, e := buildNouns(&dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestRelativeOneManyViolation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. one parent to many children
			"b", "r_1_x", "e",
			"c", "r_1_x", "e",
		)...,
	)
	if _, e := buildNouns(&dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestRelativeManyOneViolation(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. many children to one parent
			"e", "r_x_1", "b",
			"e", "r_x_1", "c",
		)...,
	)
	if _, e := buildNouns(&dt); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func newRelativeTest(nros ...string) []Ephemera {
	out := []Ephemera{
		&EphKinds{Kind: kindsOf.Relation.String()}, // declare the relation table
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "l", Ancestor: "k"},
		&EphKinds{Kind: "n", Ancestor: "k"},
		// nouns:
		&EphNouns{Noun: "a", Kind: "k"},
		&EphNouns{Noun: "b", Kind: "k"},
		&EphNouns{Noun: "c", Kind: "k"},
		&EphNouns{Noun: "d", Kind: "k"},
		&EphNouns{Noun: "e", Kind: "k"},
		&EphNouns{Noun: "f", Kind: "k"},
		&EphNouns{Noun: "l", Kind: "l"},
		&EphNouns{Noun: "n", Kind: "n"},
		&EphNouns{Noun: "z", Kind: "k"},
		// some relations those nouns can participate in.
		newRelation("r_1_1", "k", tables.ONE_TO_ONE, "k"),
		newRelation("r_1_x", "k", tables.ONE_TO_MANY, "k"),
		newRelation("r_x_1", "k", tables.MANY_TO_ONE, "k"),
		newRelation("r_x_x", "k", tables.MANY_TO_MANY, "k"),
	}
	addRelatives(&out, nros...)
	return out
}

// add noun, stem/rel, otherNoun ephemera
func addRelatives(out *[]Ephemera, nros ...string) {
	for i, cnt := 0, len(nros); i < cnt; i += 3 {
		n, r, o := nros[i], nros[i+1], nros[i+2]
		*out = append(*out, &EphRelatives{
			Noun:      n,
			Rel:       r,
			OtherNoun: o,
		})
	}
}
