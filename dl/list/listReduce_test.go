package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

func TestReduce(t *testing.T) {
	type Fruit struct {
		Name string
	}
	type Locals struct {
		Fruits  []Fruit
		Results string
	}
	type Reduce struct {
		In  Fruit
		Out string
	}
	var kinds testutil.Kinds
	kinds.AddKinds((*Fruit)(nil), (*Locals)(nil), (*Reduce)(nil))
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
			"reduce": &reduceRecords,
		},
		Runtime: testutil.Runtime{
			Kinds: &kinds,
		},
	}
	lt.Chain = scope.MakeChain(scope.FromRecord(&kinds, locals))
	if e := reduce.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if res, e := locals.GetNamedField("results"); e != nil {
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
	Target:      object.Variable("results"),
	List:        &call.FromRecordList{Value: object.Variable("fruits")},
	PatternName: ("reduce"),
}

// join each record in turn
var reduceRecords = testpat.Pattern{
	Name:   "reduce",
	Return: "out",
	Labels: []string{"in", "out"},
	Rules: []rt.Rule{{
		Exe: []rt.Execute{
			&object.SetValue{
				Target: object.Variable("out"),
				Value: &call.FromText{Value: &text.Join{
					Sep: literal.T(", "),
					Parts: []rt.TextEval{
						object.Variable("out"),
						object.Variable("in", "name"),
					}}}},
		},
	}},
}
