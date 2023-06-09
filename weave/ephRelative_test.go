package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// follow along with relative test except add list of ephemera
func TestRelativeFormation(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			"a", "r_1_1", "a",
			"b", "r_1_1", "c",
			"c", "r_1_1", "b", // duplicates the relation
			"z", "r_1_1", "e",
			//
			"z", "r_1_x", "f",
			"c", "r_1_x", "a",
			"b", "r_1_x", "e",
			"c", "r_1_x", "c", // totally okay, one lhs can have multiple rhs
			"c", "r_1_x", "b",
			"z", "r_1_x", "d",
			"z", "r_1_x", "f", // a duplicate
			//
			"z", "r_x_1", "b",
			"f", "r_x_1", "f",
			"l", "r_x_1", "b",
			"b", "r_x_1", "a",
			"d", "r_x_1", "b",
			"c", "r_x_1", "d",
			"f", "r_x_1", "f", // duplicates the earlier f-f
			"e", "r_x_1", "f",
			//
			"a", "r_x_x", "a",
			"e", "r_x_x", "d",
			"a", "r_x_x", "b",
			"a", "r_x_x", "c",
			"f", "r_x_x", "d",
			"l", "r_x_x", "d",
			"a", "r_x_x", "b", // duplicates a, b
		)...,
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readPairs(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		"a:r_1_1:a:a",
		"a:r_1_1:b:c",
		"a:r_1_1:e:z",
		//
		"a:r_1_x:b:e",
		"a:r_1_x:c:a",
		"a:r_1_x:c:b",
		"a:r_1_x:c:c",
		"a:r_1_x:z:d",
		"a:r_1_x:z:f",
		//
		"a:r_x_1:b:a",
		"a:r_x_1:c:d",
		"a:r_x_1:d:b",
		"a:r_x_1:e:f",
		"a:r_x_1:f:f",
		"a:r_x_1:l:b",
		"a:r_x_1:z:b",
		//
		"a:r_x_x:a:a",
		"a:r_x_x:a:b",
		"a:r_x_x:a:c",
		"a:r_x_x:e:d",
		"a:r_x_x:f:d",
		"a:r_x_x:l:d",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
	// expects 4 warnings, one from each group
	for i := 0; i < 4; i++ {
		if ok, e := okError(t, warnings.shift(), `Duplicate relation`); !ok {
			t.Fatal(e)
		}
	}
}

// follow along with relative test except add list of ephemera
func TestRelativeOneOneViolation(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			"b", "r_1_1", "c",
			"b", "r_1_1", "d",
		)...,
	)
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	}
}

func TestRelativeOneManyViolation(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. one parent to many children
			"b", "r_1_x", "e",
			"c", "r_1_x", "e",
		)...,
	)

	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict`); !ok {
		t.Fatal(e)
	}
}

func TestRelativeManyOneViolation(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		newRelativeTest(
			// ex. many children to one parent
			"e", "r_x_1", "b",
			"e", "r_x_1", "c",
		)...,
	)

	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict`); !ok {
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
		newRelation("r_1_1", "k", tables.ONE_TO_ONE, "k"),
		newRelation("r_1_x", "k", tables.ONE_TO_MANY, "k"),
		newRelation("r_x_1", "k", tables.MANY_TO_ONE, "k"),
		newRelation("r_x_x", "k", tables.MANY_TO_MANY, "k"),
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
