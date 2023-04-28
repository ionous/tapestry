package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

func TestValueFieldAssignment(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "l", Ancestor: "k"},
		&eph.Kinds{Kind: "m", Ancestor: "l"},
		// some simple fields
		// the name of the field has to match the name of the aspect
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			{Name: "t", Affinity: Affinity{eph.Affinity_Text}},
			{Name: "d", Affinity: Affinity{eph.Affinity_Number}},
		}},
		// nouns with those fields
		&eph.Nouns{Noun: "apple", Kind: "k"},
		&eph.Nouns{Noun: "pear", Kind: "l"},
		&eph.Nouns{Noun: "toy boat", Kind: "m"},
		&eph.Nouns{Noun: "boat", Kind: "m"},
		// values using those fields
		&eph.Values{Noun: "apple", Field: "t", Value: T("some text")},
		&eph.Values{Noun: "pear", Field: "d", Value: I(123)},
		&eph.Values{Noun: "toy", Field: "d", Value: I(321)},
		&eph.Values{Noun: "boat", Field: "t", Value: T("more text")},
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
		&eph.Kinds{Kind: "k"},
		// a field
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "d", Affinity: Affinity{eph.Affinity_Number}}}},
		// a noun
		&eph.Nouns{Noun: "n", Kind: "k"},
		// and not that field
		&eph.Values{Noun: "n", Field: "t", Value: T("no such field")},
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
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "l", Ancestor: "k"},
		&eph.Kinds{Kind: "m", Ancestor: "l"},
		// aspects
		&eph.Kinds{Kind: kindsOf.Aspect.String()},
		&eph.Aspects{Aspects: "a", Traits: dd("w", "x", "y")},
		&eph.Aspects{Aspects: "b", Traits: dd("z")},
		// fields using those aspects:
		// the name of the field has to match the name of the aspect
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			AspectParam("a"),
			AspectParam("b"),
		}},
		// nouns with those aspects
		&eph.Nouns{Noun: "apple", Kind: "k"},
		&eph.Nouns{Noun: "pear", Kind: "l"},
		&eph.Nouns{Noun: "toy boat", Kind: "m"},
		&eph.Nouns{Noun: "boat", Kind: "m"},
		// values using those aspects or traits from those aspects:
		&eph.Values{Noun: "apple", Field: "a", Value: T("y")}, // assign to the aspect directly
		&eph.Values{Noun: "pear", Field: "x", Value: B(true)}, // assign to some traits indirectly
		&eph.Values{Noun: "toy", Field: "w", Value: B(true)},
		&eph.Values{Noun: "boat", Field: "z", Value: B(true)},
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
		&eph.Kinds{Kind: kindsOf.Record.String()},
		// a record with some fields
		&eph.Kinds{Kind: "inner", Ancestor: kindsOf.Record.String(), Contain: []eph.Params{
			{Name: "num", Affinity: Affinity{eph.Affinity_Number}},
			{Name: "text", Affinity: Affinity{eph.Affinity_Text}},
		}},
		// a record holding that record
		&eph.Kinds{Kind: "outer", Ancestor: kindsOf.Record.String(), Contain: []eph.Params{
			// we use the shortcut: a field named _ of type record will (attempt) to be a kind of that record.
			{Name: "inner", Affinity: Affinity{eph.Affinity_Record}},
		}},
		//  a proper kind holding the record of records
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			{Name: "outer", Affinity: Affinity{eph.Affinity_Record}},
		}},
		// a noun of that kind, with the record of records.
		&eph.Nouns{Noun: "test", Kind: "k"},
		// values targeting a field inside the record
		&eph.Values{Noun: "test", Field: "text", Value: T("some text"), Path: []string{
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
