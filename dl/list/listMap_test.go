package list_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testpat"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestMapStrings(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name string
	}
	type Values struct {
		Fruits, Results []string
	}
	type Remap struct {
		In, Out string
	}
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil), (*Remap)(nil))
	values := kinds.NewRecord("values") // a record.
	lt := testpat.Runtime{
		testpat.Map{
			"remap": &reverseStrings,
		},
		testutil.Runtime{
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
			Kinds: &kinds,
		},
	}
	if e := values.SetNamedField("fruits", g.StringsOf([]string{"Orange", "Lemon", "Mango", "Banana", "Lime"})); e != nil {
		t.Fatal(e)
	} else if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if results, e := values.GetNamedField("results"); e != nil {
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
	type Values struct {
		Fruits  []Fruit
		Results []Fruit
	}
	type Remap struct {
		In, Out Fruit
	}
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil), (*Remap)(nil))
	values := kinds.NewRecord("values")
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
		if e := values.SetNamedField("fruits", g.RecordsOf(k.Name(), fruits)); e != nil {
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
				g.RecordOf(values),
			},
		},
	}
	if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if val, e := values.GetNamedField("results"); e != nil {
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

var remap = list.Map{
	FromList:     &core.Var{Name: "fruits"},
	ToList:       "results",
	UsingPattern: "remap",
}

var reverseRecords = testpat.Pattern{
	Name:   "remap",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Execute: &core.PutAtField{
			Into:    &core.IntoVar{"out"},
			AtField: "name",
			From: &core.FromText{
				&core.MakeReversed{
					&core.GetAtField{
						Field: "name",
						From:  &core.FromVar{"in"},
					},
				},
			},
		},
	},
	},
}

var reverseStrings = testpat.Pattern{
	Name:   "remap",
	Labels: []string{"in"},
	Return: "out",
	Rules: []rt.Rule{{
		Execute: &core.Assign{
			Var: "out",
			From: &core.FromText{
				&core.MakeReversed{
					&core.Var{
						Name: "in",
					},
				},
			},
		},
	},
	},
}
