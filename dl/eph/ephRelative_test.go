package eph

import (
	"testing"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/kr/pretty"
)

// follow along with relative test except add list of ephemera
func TestRelativeFormation(t *testing.T) {
	var dt domainTest
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
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_pair}
		if e := cat.WritePairs(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"a:a:r_1_1:a:x",
			"a:b:r_1_1:c:x",
			"a:e:r_1_1:z:x",
			//
			"a:b:r_1_x:e:x",
			"a:c:r_1_x:a:x",
			"a:c:r_1_x:b:x",
			"a:c:r_1_x:c:x",
			"a:z:r_1_x:d:x",
			"a:z:r_1_x:f:x",
			//
			"a:b:r_x_1:a:x",
			"a:c:r_x_1:d:x",
			"a:d:r_x_1:b:x",
			"a:e:r_x_1:f:x",
			"a:f:r_x_1:f:x",
			"a:l:r_x_1:b:x",
			"a:z:r_x_1:b:x",
			//
			"a:a:r_x_x:a:x",
			"a:a:r_x_x:b:x",
			"a:a:r_x_x:c:x",
			"a:e:r_x_x:d:x",
			"a:f:r_x_x:d:x",
			"a:l:r_x_x:d:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// follow along with relative test except add list of ephemera
func TestRelativeOneOneViolation(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		newRelativeTest(
			"b", "r_1_1", "c",
			"b", "r_1_1", "d",
		)...,
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestRelativeOneManyViolation(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. one parent to many children
			"b", "r_1_x", "e",
			"c", "r_1_x", "e",
		)...,
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func TestRelativeManyOneViolation(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. many children to one parent
			"e", "r_x_1", "b",
			"e", "r_x_1", "c",
		)...,
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
		NounPhase:     NounPhaseActions,
	}); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok", e)
	}
}

func newRelativeTest(nros ...string) []Ephemera {
	out := []Ephemera{
		&EphKinds{KindsOfRelation, ""}, // declare the relation table
		&EphKinds{"k", ""},
		&EphKinds{"l", "k"},
		&EphKinds{"n", "k"},
		// nouns:
		&EphNouns{"a", "k"},
		&EphNouns{"b", "k"},
		&EphNouns{"c", "k"},
		&EphNouns{"d", "k"},
		&EphNouns{"e", "k"},
		&EphNouns{"f", "k"},
		&EphNouns{"l", "l"},
		&EphNouns{"n", "n"},
		&EphNouns{"z", "k"},
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
