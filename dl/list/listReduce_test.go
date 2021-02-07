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

func TestReduce(t *testing.T) {
	type Fruit struct {
		Name string
	}
	type Values struct {
		Fruits  []Fruit
		Results string
	}

	var kinds testutil.Kinds
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
			"reduce": &reduceRecords,
		},
		Stack: []rt.Scope{
			g.RecordOf(values),
		},
	}
	if e := reduce.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if res, e := values.GetNamedField("Results"); e != nil {
		t.Fatal(e)
	} else {
		out := res.String()
		expected := "Orange, Lemon, Mango, Banana, Lime"
		if expected != out {
			t.Fatal(out)
		} else {
			pretty.Println("ok", out) // emiL
		}
	}
}

var reduce = list.Reduce{
	FromList:     &core.FromRecords{&core.Var{"Fruits"}},
	IntoValue:    "Results",
	UsingPattern: "reduce",
}

// join each record in turn
var reduceRecords = pattern.Pattern{
	Name:   "reduce",
	Return: "out",
	Labels: []string{"in", "out"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Record, Type: "Fruit"},
		{Name: "out", Affinity: affine.Text},
	},
	Rules: []*pattern.Rule{
		&pattern.Rule{
			Execute: &core.Assign{
				Var: N("out"),
				From: &core.FromText{&core.Join{Sep: T(", "), Parts: []rt.TextEval{
					V("out"),
					&core.GetAtField{Field: "Name", From: &core.FromVar{N("in")}},
				}}},
			},
		},
	},
}
