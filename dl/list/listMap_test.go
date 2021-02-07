package list_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
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
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil))
	values := kinds.New("Values") // a record.
	lt := testutil.Runtime{
		PatternMap: testutil.PatternMap{
			"remap": &reverseStrings,
		},
		Stack: []rt.Scope{
			g.RecordOf(values),
		},
		Kinds: &kinds,
	}
	if e := values.SetNamedField("Fruits", g.StringsOf([]string{"Orange", "Lemon", "Mango", "Banana", "Lime"})); e != nil {
		t.Fatal(e)
	} else if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if results, e := values.GetNamedField("Results"); e != nil {
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
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil))
	values := kinds.New("Values")
	if k, e := kinds.GetKindByName("Fruit"); e != nil {
		t.Fatal(e)
	} else {
		var fruits []*g.Record
		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
			one := k.NewRecord()
			if e := one.SetNamedField("Name", g.StringOf(f)); e != nil {
				t.Fatal(e)
			}
			fruits = append(fruits, one)
		}
		if e := values.SetNamedField("Fruits", g.RecordsOf(k.Name(), fruits)); e != nil {
			t.Fatal(e)
		}
	}
	//
	lt := testutil.Runtime{
		Kinds: &kinds,
		PatternMap: testutil.PatternMap{
			"remap": &reverseRecords,
		},
		Stack: []rt.Scope{
			g.RecordOf(values),
		},
	}
	if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if val, e := values.GetNamedField("Results"); e != nil {
		t.Fatal(e)
	} else if res := val.Records(); len(res) != 5 {
		t.Fatal("missing results")
	} else {
		expect := []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}
		var got []string
		for _, el := range res {
			if v, e := el.GetNamedField("Name"); e != nil {
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

var remap = list.Map{FromList: &core.Var{Name: "Fruits"}, ToList: "Results", UsingPattern: "remap"}

var reverseRecords = pattern.Pattern{
	Name:   "remap",
	Return: "out",
	Labels: []string{"in"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Record, Type: "Fruit"},
		{Name: "out", Affinity: affine.Record, Type: "Fruit"},
	},
	Rules: []*pattern.Rule{
		&pattern.Rule{
			Execute: &core.PutAtField{
				Into:    &core.IntoVar{N("out")},
				AtField: "Name",
				From: &core.FromText{
					&core.MakeReversed{
						&core.GetAtField{
							Field: "Name",
							From:  &core.FromVar{N("in")},
						},
					},
				},
			},
		},
	},
}

var reverseStrings = pattern.Pattern{
	Name:   "remap",
	Labels: []string{"in"},
	Return: "out",
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Text},
	},
	Rules: []*pattern.Rule{
		&pattern.Rule{
			Execute: &core.Assign{
				Var: core.Variable{Str: "out"},
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
