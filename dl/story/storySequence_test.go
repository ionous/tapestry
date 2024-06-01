package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/printer"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// test that importing cycling text transforms to the proper runtime command
func TestImportSequence(t *testing.T) {
	db := tables.CreateTest(t.Name(), true)
	defer db.Close()
	cat := weave.NewCatalog(db)

	// create a statement that uses CycleText
	printText := &printer.PrintText{Text: &story.CycleText{Parts: []rt.TextEval{
		literal.T("a"),
		literal.T("b"),
		literal.T("c"),
	}}}
	// import that statement
	if e := story.Weave(cat, []story.StoryStatement{
		&story.DefineTest{
			TestName: t.Name(),
			Exe: []rt.Execute{
				printText,
			},
		},
	}); e != nil {
		t.Fatal(e)
	} else {
		// validate that it was transformed into this:
		expect := printer.CallCycle{
			Name: "seq-1",
			Parts: []rt.TextEval{
				T("a"),
				T("b"),
				T("c"),
			},
		}
		if diff := pretty.Diff(printText.Text, &expect); len(diff) > 0 {
			t.Fatal("want:", expect, "have:", pretty.Sprint(printText.Text))
		}
	}
}
