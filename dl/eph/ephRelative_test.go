package eph

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/literal"
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
	dt := domainTest{noShuffle: true}
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
	dt := domainTest{noShuffle: true}
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
	dt := domainTest{noShuffle: true}
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
	var out []Ephemera
	addKinds(&out,
		KindsOfRelation, "", // declare the relation table
		"k", "",
		"l", "k",
		"n", "k",
	)
	addNouns(&out,
		"a", "k",
		"b", "k",
		"c", "k",
		"d", "k",
		"e", "k",
		"f", "k",
		"l", "l",
		"n", "n",
		"z", "k",
	)
	addRelations(&out,
		"r_1_1", "k", tables.ONE_TO_ONE, "k",
		"r_1_x", "k", tables.ONE_TO_MANY, "k",
		"r_x_1", "k", tables.MANY_TO_ONE, "k",
		"r_x_x", "k", tables.MANY_TO_MANY, "k",
	)
	addRelatives(&out, nros...)
	return out
}

// kind, parent
func addKinds(out *[]Ephemera, kps ...string) {
	for i, cnt := 0, len(kps); i < cnt; i += 2 {
		k, p := kps[i], kps[i+1]
		*out = append(*out, &EphKinds{Kinds: k, From: p})
	}
}

// kind, name, affinity(key), class
func addFields(out *[]Ephemera, knacs ...string) {
	for i, cnt := 0, len(knacs); i < cnt; i += 4 {
		k, n, a, c := knacs[i], knacs[i+1], knacs[i+2], knacs[i+3]
		*out = append(*out, &EphFields{
			Kinds: k, Name: n, Affinity: Affinity{a}, Class: c,
		})
	}
}

// noun, kind
func addNouns(out *[]Ephemera, nks ...string) {
	for i, cnt := 0, len(nks); i < cnt; i += 2 {
		n, k := nks[i], nks[i+1]
		*out = append(*out, &EphNouns{Noun: n, Kind: k})
	}
}

// noun(string), field(string), value(literal)
func addValues(out *[]Ephemera, nfvs ...interface{}) {
	for i, cnt := 0, len(nfvs); i < cnt; i += 2 {
		n, f, v := nfvs[i].(string), nfvs[i+1].(string), nfvs[i+2].(literal.LiteralValue)
		*out = append(*out, &EphValues{Noun: n, Field: f, Value: v})
	}
}

// relation, kind, cardinality, otherKinds
func addRelations(out *[]Ephemera, rkcos ...string) {
	for i, cnt := 0, len(rkcos); i < cnt; i += 4 {
		r, k, c, o := rkcos[i], rkcos[i+1], rkcos[i+2], rkcos[i+3]
		var card EphCardinality
		switch c {
		case tables.ONE_TO_ONE:
			card = EphCardinality{EphCardinality_OneOne_Opt, &OneOne{k, o}}
		case tables.ONE_TO_MANY:
			card = EphCardinality{EphCardinality_OneMany_Opt, &OneMany{k, o}}
		case tables.MANY_TO_ONE:
			card = EphCardinality{EphCardinality_ManyOne_Opt, &ManyOne{k, o}}
		case tables.MANY_TO_MANY:
			card = EphCardinality{EphCardinality_ManyMany_Opt, &ManyMany{k, o}}
		default:
			panic("unknown cardinality")
		}
		*out = append(*out, &EphRelations{
			Rel:         r,
			Cardinality: card,
		})
	}
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
