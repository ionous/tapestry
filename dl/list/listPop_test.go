package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

func TestPopping(t *testing.T) {
	// pop from the front of a list
	// front := popTest(true, 5, "Orange", "Lemon", "Mango")
	// if d := pretty.Diff(front, []string{"Orange", "Lemon", "Mango", "x", "x"}); len(d) > 0 {
	// 	t.Fatal("pop front", front)
	// }
	// pop from the back of a list
	back := popTest(false, 5, "Orange", "Lemon", "Mango")
	if d := pretty.Diff(back, []string{"Mango", "Lemon", "Orange", "x", "x"}); len(d) > 0 {
		t.Fatal("pop back", back)
	}
}

func popTest(front bool, amt int, src ...string) []string {
	var out []string
	var start int
	if front {
		start = 1
	} else {
		start = -1
	}
	// this will be run "amt" times
	pop := &list.Erasing{
		Count:   literal.I(1),
		AtIndex: literal.I(start),
		Target:  object.Variable("source"),
		As:      "text",
		Exe: []rt.Execute{
			&logic.ChooseBranch{
				Condition: &math.CompareNum{
					A: &list.ListLen{
						List: &assign.FromTextList{Value: object.Variable("text")},
					},
					Is: math.C_Comparison_EqualTo,
					B:  literal.I(0),
				},
				Exe: []rt.Execute{
					&Write{&out, literal.T("x")},
				},
				Else: &logic.ChooseNothingElse{
					Exe: []rt.Execute{
						&Write{&out, object.Variable("text", 1)},
					},
				},
			}},
	}
	if run, e := newListTime(src, nil); e != nil {
		panic(e)
	} else {
		for i := 0; i < amt; i++ {
			if e := pop.Execute(run); e != nil {
				panic(e)
			}
		}
	}
	return out
}
