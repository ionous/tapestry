package eph

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

func TestInitialFieldAssignment(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "l", From: "k"},
		&EphKinds{Kinds: "m", From: "l"},
		// some simple fields
		// the name of the field has to match the name of the aspect
		&EphKinds{Kinds: "k", Contain: []EphParams{
			{Name: "t", Affinity: Affinity{Affinity_Text}},
			{Name: "d", Affinity: Affinity{Affinity_Number}},
		}},
		// nouns with those fields
		&EphNouns{"apple", "k"},
		&EphNouns{"pear", "l"},
		&EphNouns{"toy boat", "m"},
		&EphNouns{"boat", "m"},
		// values using those fields
		&EphValues{"apple", "t", T("some text")},
		&EphValues{"pear", "d", I(123)},
		&EphValues{"toy", "d", I(321)},
		&EphValues{"boat", "t", T("more text")},
	)
	if cat, e := buildNouns(dt); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Value}
		if e := cat.WriteValues(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			`a:apple:t:"some text":x`,
			`a:boat:t:"more text":x`,
			`a:toy_boat:d:321:x`,
			`a:pear:d:123:x`,
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

func TestMissingField(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kinds: "k"},
		// a field
		&EphKinds{Kinds: "k", Contain: []EphParams{{Name: "d", Affinity: Affinity{Affinity_Number}}}},
		// a noun
		&EphNouns{"n", "k"},
		// and not that field
		&EphValues{"n", "t", T("no such field")},
	)
	if _, e := buildNouns(dt); e == nil || e.Error() != `field not found 'k.t'` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok", e)
	}
}

func TestInitialTraitAssignment(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kinds: "k"},
		&EphKinds{Kinds: "l", From: "k"},
		&EphKinds{Kinds: "m", From: "l"},
		// aspects
		&EphKinds{Kinds: kindsOf.Aspect.String()},
		&EphAspects{"a", dd("w", "x", "y")},
		&EphAspects{"b", dd("z")},
		// fields using those aspects:
		// the name of the field has to match the name of the aspect
		&EphKinds{Kinds: "k", Contain: []EphParams{
			AspectParam("a"),
			AspectParam("b"),
		}},
		// nouns with those aspects
		&EphNouns{"apple", "k"},
		&EphNouns{"pear", "l"},
		&EphNouns{"toy boat", "m"},
		&EphNouns{"boat", "m"},
		// values using those aspects or traits from those aspects:
		&EphValues{"apple", "a", T("y")}, // assign to the aspect directly
		&EphValues{"pear", "x", B(true)}, // assign to some traits indirectly
		&EphValues{"toy", "w", B(true)},
		&EphValues{"boat", "z", B(true)},
	)
	if cat, e := buildNouns(dt); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Value}
		if e := cat.WriteValues(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			`a:apple:a:"y":x`,
			`a:boat:b:"z":x`,
			`a:toy_boat:a:"w":x`,
			`a:pear:a:"x":x`,
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
