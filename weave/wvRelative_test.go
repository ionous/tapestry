package weave_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/kr/pretty"
)

// follow along with relative test except add list of ephemera
func TestRelativeFormation(t *testing.T) {
	var warnings mdl.Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		newRelativeTest(
			"a", "r 1 1", "a",
			"b", "r 1 1", "c",
			"c", "r 1 1", "b", // duplicates the relation
			"z", "r 1 1", "e",
			//
			"z", "r 1 x", "f",
			"c", "r 1 x", "a",
			"b", "r 1 x", "e",
			"c", "r 1 x", "c", // totally okay, one lhs can have multiple rhs
			"c", "r 1 x", "b",
			"z", "r 1 x", "d",
			"z", "r 1 x", "f", // a duplicate
			//
			"z", "r x 1", "b",
			"f", "r x 1", "f",
			"l", "r x 1", "b",
			"b", "r x 1", "a",
			"d", "r x 1", "b",
			"c", "r x 1", "d",
			"f", "r x 1", "f", // duplicates the earlier f-f
			"e", "r x 1", "f",
			//
			"a", "r x x", "a",
			"e", "r x x", "d",
			"a", "r x x", "b",
			"a", "r x x", "c",
			"f", "r x x", "d",
			"l", "r x x", "d",
			"a", "r x x", "b", // duplicates a, b
		)...,
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.ReadPairs(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:r 1 1:a:a",
		"a:r 1 1:b:c",
		"a:r 1 1:e:z",
		//
		"a:r 1 x:b:e",
		"a:r 1 x:c:a",
		"a:r 1 x:c:b",
		"a:r 1 x:c:c",
		"a:r 1 x:z:d",
		"a:r 1 x:z:f",
		//
		"a:r x 1:b:a",
		"a:r x 1:c:d",
		"a:r x 1:d:b",
		"a:r x 1:e:f",
		"a:r x 1:f:f",
		"a:r x 1:l:b",
		"a:r x 1:z:b",
		//
		"a:r x x:a:a",
		"a:r x x:a:b",
		"a:r x x:a:c",
		"a:r x x:e:d",
		"a:r x x:f:d",
		"a:r x x:l:d",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
	// expects 4 warnings, one from each group
	for i := 0; i < 4; i++ {
		if e := warnings.Expect(`Duplicate relation`); e != nil {
			t.Fatal(e)
		}
	}
}

// follow along with relative test except add list of ephemera
func TestRelativeOneOneViolation(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		newRelativeTest(
			"b", "r 1 1", "c",
			"b", "r 1 1", "d",
		)...,
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	}
}

func TestRelativeOneManyViolation(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		newRelativeTest(
			// ex. one parent to many children
			"b", "r 1 x", "e",
			"c", "r 1 x", "e",
		)...,
	)

	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	}
}

func TestRelativeManyOneViolation(t *testing.T) {
	dt := testweave.NewWeaver(t.Name())
	defer dt.Close()
	dt.MakeDomain(dd("a"),
		newRelativeTest(
			// ex. many children to one parent
			"e", "r x 1", "b",
			"e", "r x 1", "c",
		)...,
	)
	_, e := dt.Assemble()
	if ok, e := testweave.OkayError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	}
}

func newRelativeTest(nros ...string) []eph.Ephemera {
	out := []eph.Ephemera{
		&eph.Kinds{Kind: kindsOf.Relation.String()}, // declare the relation table
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "l", Ancestor: "k"},
		&eph.Kinds{Kind: "n", Ancestor: "k"},
		// nouns:
		&eph.Nouns{Noun: "a", Kind: "k"},
		&eph.Nouns{Noun: "b", Kind: "k"},
		&eph.Nouns{Noun: "c", Kind: "k"},
		&eph.Nouns{Noun: "d", Kind: "k"},
		&eph.Nouns{Noun: "e", Kind: "k"},
		&eph.Nouns{Noun: "f", Kind: "k"},
		&eph.Nouns{Noun: "l", Kind: "l"},
		&eph.Nouns{Noun: "n", Kind: "n"},
		&eph.Nouns{Noun: "z", Kind: "k"},
		// some relations those nouns can participate in.
		newRelation("r 1 1", "k", tables.ONE_TO_ONE, "k"),
		newRelation("r 1 x", "k", tables.ONE_TO_MANY, "k"),
		newRelation("r x 1", "k", tables.MANY_TO_ONE, "k"),
		newRelation("r x x", "k", tables.MANY_TO_MANY, "k"),
	}
	addRelatives(&out, nros...)
	return out
}

// add noun, stem/rel, otherNoun ephemera
func addRelatives(out *[]eph.Ephemera, nros ...string) {
	for i, cnt := 0, len(nros); i < cnt; i += 3 {
		n, r, o := nros[i], nros[i+1], nros[i+2]
		*out = append(*out, &eph.Relatives{
			Noun:      n,
			Rel:       r,
			OtherNoun: o,
		})
	}
}
