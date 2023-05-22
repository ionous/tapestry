package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

func TestValueFieldAssignment(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&eph.Kinds{Kind: "k"},
		&eph.Kinds{Kind: "l", Ancestor: "k"},
		&eph.Kinds{Kind: "m", Ancestor: "l"},
		// some simple fields
		// the name of the field has to match the name of the aspect
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			{Name: "t", Affinity: affine.Text},
			{Name: "d", Affinity: affine.Number},
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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readValues(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		`a:apple:t:"some text"`,
		`a:boat:t:"more text"`,
		`a:toy_boat:d:321`,
		`a:pear:d:123`,
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func TestMissingField(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// some random set of kinds
		&eph.Kinds{Kind: "k"},
		// a field
		&eph.Kinds{Kind: "k", Contain: []eph.Params{{Name: "d", Affinity: affine.Number}}},
		// a noun
		&eph.Nouns{Noun: "n", Kind: "k"},
		// and not that field
		&eph.Values{Noun: "n", Field: "t", Value: T("no such field")},
	)

	if _, e := dt.Assemble(); e == nil || e.Error() != `field not found 'k.t'` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok", e)
	}
}

func TestValueTraitAssignment(t *testing.T) {
	dt := newTest(t.Name())
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

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readValues(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		`a:apple:a:"y"`,
		`a:boat:b:"z"`,
		`a:toy_boat:a:"w"`,
		`a:pear:a:"x"`,
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

func TestValuePaths(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// declare the existence of records
		&eph.Kinds{Kind: kindsOf.Record.String()},
		// a record with some fields
		&eph.Kinds{Kind: "inner", Ancestor: kindsOf.Record.String(), Contain: []eph.Params{
			{Name: "num", Affinity: affine.Number},
			{Name: "text", Affinity: affine.Text},
		}},
		// a record holding that record
		&eph.Kinds{Kind: "outer", Ancestor: kindsOf.Record.String(), Contain: []eph.Params{
			// we use the shortcut: a field named _ of type record will (attempt) to be a kind of that record.
			{Name: "inner", Affinity: affine.Record},
		}},
		//  a proper kind holding the record of records
		&eph.Kinds{Kind: "k", Contain: []eph.Params{
			{Name: "outer", Affinity: affine.Record},
		}},
		// a noun of that kind, with the record of records.
		&eph.Nouns{Noun: "test", Kind: "k"},
		// values targeting a field inside the record
		&eph.Values{Noun: "test", Field: "text", Value: T("some text"), Path: []string{
			"outer", "inner",
		}},
	)

	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if out, e := dt.readValues(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{
		// `a:test:outer:{"Fields:":[{"Field field:value:":["inner",{"Fields:":[{"Field field:value:":["text","some text"]}]}]}]}`,
		`a:test:outer:{"inner":{"text":"some text"}}`,
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
