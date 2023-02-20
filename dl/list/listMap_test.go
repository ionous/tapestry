package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
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
	type Remap struct {
		In, Out string
	}
	kinds.AddKinds((*Fruit)(nil), (*Locals)(nil), (*Remap)(nil))
	locals := kinds.NewRecord("locals") // a record.
	lt := testpat.Runtime{
		testpat.Map{
			"remap": &reverseStrings,
		},
		testutil.Runtime{
			Stack: []rt.Scope{
				g.RecordOf(locals),
			},
			Kinds: &kinds,
		},
	}
	if e := locals.SetNamedField("fruits", g.StringsOf([]string{"Orange", "Lemon", "Mango", "Banana", "Lime"})); e != nil {
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
	type Remap struct {
		In, Out Fruit
	}
	kinds.AddKinds((*Fruit)(nil), (*Locals)(nil), (*Remap)(nil))
	locals := kinds.NewRecord("locals")
	if k, e := kinds.GetKindByName("fruit"); e != nil {
		t.Fatal(e)
	} else {
		var fruits []*g.Record
		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
			one := k.NewRecord()
			if e := one.SetNamedField("name", g.StringOf(f)); e != nil {
				t.Fatal(e)
			}
			fruits = append(fruits, one)
		}
		if e := locals.SetNamedField("fruits", g.RecordsFrom(fruits, k.Name())); e != nil {
			t.Fatal(e)
		}
	}
	//
	lt := testpat.Runtime{
		testpat.Map{
			"remap": &reverseRecords,
		},
		testutil.Runtime{
			Kinds: &kinds,
			Stack: []rt.Scope{
				g.RecordOf(locals),
			},
		},
	}
	if e := remapRecords.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if val, e := locals.GetNamedField("results"); e != nil {
		t.Fatal(e)
	} else if res := val.Records(); len(res) != 5 {
		t.Fatal("missing results")
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
	Target:       core.Variable("results"),
	List:         &assign.FromTextList{Value: core.GetVariable("fruits")},
	UsingPattern: W("remap"),
}

var remapRecords = list.ListMap{
	Target:       core.Variable("results"),
	List:         &assign.FromRecordList{Value: core.GetVariable("fruits")},
	UsingPattern: W("remap"),
}

var reverseStrings = testpat.Pattern{
	Name:   "remap",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Execute: core.MakeActivity(
			&core.SetValue{
				Target: core.Variable("out"),
				Value:  &assign.FromText{Value: &core.MakeReversed{Text: GetVariable("in")}}},
		),
	}},
}

var reverseRecords = testpat.Pattern{
	Name:   "remap",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Execute: core.MakeActivity(
			&core.SetValue{
				Target: core.Variable("out", "name"),
				Value:  &assign.FromText{Value: &core.MakeReversed{Text: GetVariable("in", "name")}}},
		),
	}},
}
