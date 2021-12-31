package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
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
	pop := &list.Erasing{
		Count:   I(1),
		AtIndex: I(start),
		From:    &list.FromTxtList{N("source")},
		As:      W("text"),
		Do: core.MakeActivity(&core.ChooseAction{
			If: &core.CompareNum{&list.ListLen{V("text")}, &core.Equal{}, I(0)},
			Do: core.MakeActivity(&Write{&out, T("x")}),
			Else: &core.ChooseNothingElse{
				Do: core.MakeActivity(&Write{&out, &list.ListAt{V("text"), I(1)}}),
			},
		}),
	}
	if run, _, e := newListTime(src, nil); e != nil {
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
