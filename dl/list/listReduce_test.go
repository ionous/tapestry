package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
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
	type Reduce struct {
		In  Fruit
		Out string
	}
	var kinds testutil.Kinds
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil), (*Reduce)(nil))
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
		if e := values.SetNamedField("fruits", g.RecordsFrom(fruits, k.Name())); e != nil {
			t.Fatal(e)
		}
	}
	//
	lt := testpat.Runtime{
		testpat.Map{
			"reduce": &reduceRecords,
		},
		testutil.Runtime{
			Kinds: &kinds,
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
		},
	}
	if e := reduce.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if res, e := values.GetNamedField("results"); e != nil {
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

var reduce = list.ListReduce{
	FromList:     &core.FromRecords{Vals: V("fruits")},
	IntoValue:    W("results"),
	UsingPattern: W("reduce"),
}

// join each record in turn
var reduceRecords = testpat.Pattern{
	Name:   "reduce",
	Return: "out",
	Labels: []string{"in", "out"},
	Rules: []rt.Rule{{
		Execute: core.MakeActivity(&core.Assign{
			Var: N("out"),
			From: &core.FromText{Val: &core.Join{Sep: T(", "), Parts: []rt.TextEval{
				V("out"),
				&core.GetAtField{Field: W("name"), From: &core.FromVar{Var: N("in")}},
			}}},
		}),
	}},
}
