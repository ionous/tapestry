package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

func TestMapStrings(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name string
	}
	type Locals struct {
		Fruits, Results []string
	}
	type Reverse struct {
		In, Out string
	}
	kinds.AddKinds((*Fruit)(nil), (*Locals)(nil), (*Reverse)(nil))
	locals := kinds.NewRecord("locals") // a record.
	lt := testpat.Runtime{
		Map: testpat.Map{
			"reverse": &reverseText,
		},
		Runtime: testutil.Runtime{
			Kinds: &kinds,
		},
	}
	lt.Chain = scope.MakeChain(scope.FromRecord(&kinds, locals))
	if e := locals.SetNamedField("fruits", rt.StringsOf([]string{"Orange", "Lemon", "Mango", "Banana", "Lime"})); e != nil {
		t.Fatal(e)
	} else if e := remapStrings.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if results, e := locals.GetNamedField("results"); e != nil {
		t.Fatal(e)
	} else {
		res := results.Strings()
		if diff := pretty.Diff(res, []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}); len(diff) > 0 {
			t.Fatal(res)
		} else {
			t.Log("ok", res)
		}
	}
}

func TestMapRecords(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name string
	}
	type Locals struct {
		Fruits  []Fruit
		Results []Fruit
	}
	type Reverse struct {
		In, Out Fruit
	}
	kinds.AddKinds((*Fruit)(nil), (*Locals)(nil), (*Reverse)(nil))
	locals := kinds.NewRecord("locals")
	if k, e := kinds.GetKindByName("fruit"); e != nil {
		t.Fatal(e)
	} else {
		var fruits []*rt.Record
		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
			one := rt.NewRecord(k)
			if e := one.SetNamedField("name", rt.StringOf(f)); e != nil {
				t.Fatal(e)
			}
			fruits = append(fruits, one)
		}
		if e := locals.SetNamedField("fruits", rt.RecordsFrom(fruits, k.Name())); e != nil {
			t.Fatal(e)
		}
	}
	//
	lt := testpat.Runtime{
		Map: testpat.Map{
			"reverse": &reverseField,
		},
		Runtime: testutil.Runtime{
			Kinds: &kinds,
		},
	}
	lt.Chain = scope.MakeChain(scope.FromRecord(&kinds, locals))
	if e := remapRecords.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if val, e := locals.GetNamedField("results"); e != nil {
		t.Fatal(e)
	} else if res := val.Records(); len(res) != 5 {
		t.Fatal("expected 5 results, got:", len(res))
	} else {
		expect := []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}
		var got []string
		for _, el := range res {
			if v, e := el.GetNamedField("name"); e != nil {
				t.Fatal(e)
			} else {
				got = append(got, v.String())
			}
		}
		if diff := pretty.Diff(expect, got); len(diff) > 0 {
			t.Fatal("error", got)
		}
	}
}

var remapStrings = list.ListMap{
	Target:      object.Variable("results"),
	List:        &call.FromTextList{Value: object.Variable("fruits")},
	PatternName: ("reverse"),
}

var remapRecords = list.ListMap{
	Target:      object.Variable("results"),
	List:        &call.FromRecordList{Value: object.Variable("fruits")},
	PatternName: ("reverse"),
}

// a pattern which takes a string from "in" and returns the reverse of it via "out"
var reverseText = testpat.Pattern{
	Name:   "reverse",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Exe: []rt.Execute{
			&object.SetValue{
				Target: object.Variable("out"),
				Value:  &call.FromText{Value: &text.MakeReversed{Text: object.Variable("in")}}},
		},
	}},
}

// a pattern which takes a string from "in.name" and returns the reverse of it via "out.name"
var reverseField = testpat.Pattern{
	Name:   "reverse",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Exe: []rt.Execute{
			&object.SetValue{
				Target: object.Variable("out", "name"),
				Value:  &call.FromText{Value: &text.MakeReversed{Text: object.Variable("in", "name")}}},
		},
	}},
}
