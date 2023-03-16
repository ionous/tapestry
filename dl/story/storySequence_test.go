package story_test

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// test that importing cycling text transforms to the proper runtime command
func TestImportSequence(t *testing.T) {
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els))
	// create a statement that uses CycleText
	printText := &core.PrintText{Text: &story.CycleText{Parts: []rt.TextEval{
		core.T("a"),
		core.T("b"),
		core.T("c"),
	}}}
	// import that statement
	if e := story.ImportStory(k, t.Name(), &story.StoryFile{
		StoryStatements: []story.StoryStatement{
			&debug.Test{
				Do: []rt.Execute{
					printText,
				},
			},
		},
	}); e != nil {
		t.Fatal(e)
	} else {
		// validate that it was transformed into this:
		expect := core.CallCycle{
			Name: "seq_1",
			Parts: []rt.TextEval{
				T("a"),
				T("b"),
				T("c"),
			},
		}
		if diff := pretty.Diff(printText.Text, &expect); len(diff) > 0 {
			t.Fatal(pretty.Print("want:", expect, "have:", printText.Text))
		}
	}
}
