package eph

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

func TestValueFieldAssignment(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "l", Ancestor: "k"},
		&EphKinds{Kind: "m", Ancestor: "l"},
		// some simple fields
		// the name of the field has to match the name of the aspect
		&EphKinds{Kind: "k", Contain: []EphParams{
			{Name: "t", Affinity: Affinity{Affinity_Text}},
			{Name: "d", Affinity: Affinity{Affinity_Number}},
		}},
		// nouns with those fields
		&EphNouns{Noun: "apple", Kind: "k"},
		&EphNouns{Noun: "pear", Kind: "l"},
		&EphNouns{Noun: "toy boat", Kind: "m"},
		&EphNouns{Noun: "boat", Kind: "m"},
		// values using those fields
		&EphValues{Noun: "apple", Field: "t", Value: T("some text")},
		&EphValues{Noun: "pear", Field: "d", Value: I(123)},
		&EphValues{Noun: "toy", Field: "d", Value: I(321)},
		&EphValues{Noun: "boat", Field: "t", Value: T("more text")},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kind: "k"},
		// a field
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "d", Affinity: Affinity{Affinity_Number}}}},
		// a noun
		&EphNouns{Noun: "n", Kind: "k"},
		// and not that field
		&EphValues{Noun: "n", Field: "t", Value: T("no such field")},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil || e.Error() != `field not found 'k.t'` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok", e)
	}
}

func TestValueTraitAssignment(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "l", Ancestor: "k"},
		&EphKinds{Kind: "m", Ancestor: "l"},
		// aspects
		&EphKinds{Kind: kindsOf.Aspect.String()},
		&EphAspects{Aspects: "a", Traits: dd("w", "x", "y")},
		&EphAspects{Aspects: "b", Traits: dd("z")},
		// fields using those aspects:
		// the name of the field has to match the name of the aspect
		&EphKinds{Kind: "k", Contain: []EphParams{
			AspectParam("a"),
			AspectParam("b"),
		}},
		// nouns with those aspects
		&EphNouns{Noun: "apple", Kind: "k"},
		&EphNouns{Noun: "pear", Kind: "l"},
		&EphNouns{Noun: "toy boat", Kind: "m"},
		&EphNouns{Noun: "boat", Kind: "m"},
		// values using those aspects or traits from those aspects:
		&EphValues{Noun: "apple", Field: "a", Value: T("y")}, // assign to the aspect directly
		&EphValues{Noun: "pear", Field: "x", Value: B(true)}, // assign to some traits indirectly
		&EphValues{Noun: "toy", Field: "w", Value: B(true)},
		&EphValues{Noun: "boat", Field: "z", Value: B(true)},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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

func TestValuePaths(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// declare the existence of records
		&EphKinds{Kind: kindsOf.Record.String()},
		// a record with some fields
		&EphKinds{Kind: "inner", Ancestor: kindsOf.Record.String(), Contain: []EphParams{
			{Name: "num", Affinity: Affinity{Affinity_Number}},
			{Name: "text", Affinity: Affinity{Affinity_Text}},
		}},
		// a record holding that record
		&EphKinds{Kind: "outer", Ancestor: kindsOf.Record.String(), Contain: []EphParams{
			// we use the shortcut: a field named _ of type record will (attempt) to be a kind of that record.
			{Name: "inner", Affinity: Affinity{Affinity_Record}},
		}},
		//  a proper kind holding the record of records
		&EphKinds{Kind: "k", Contain: []EphParams{
			{Name: "outer", Affinity: Affinity{Affinity_Record}},
		}},
		// a noun of that kind, with the record of records.
		&EphNouns{Noun: "test", Kind: "k"},
		// values targeting a field inside the record
		&EphValues{Noun: "test", Field: "text", Value: T("some text"), Path: []string{
			"outer", "inner",
		}},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Value}
		if e := cat.WriteValues(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			// `a:test:outer:{"Fields:":[{"Field field:value:":["inner",{"Fields:":[{"Field field:value:":["text","some text"]}]}]}]}:x`,
			`a:test:outer:{"inner":{"text":"some text"}}:x`,
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
