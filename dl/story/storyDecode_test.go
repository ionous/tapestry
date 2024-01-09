package story_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"
)

func TestDecodeStory(t *testing.T) {
	var m map[string]any
	if e := json.Unmarshal([]byte(debug.FactorialJs), &m); e != nil {
		t.Fatal(e)
	} else if file, e := story.Decode(m); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, file); len(diff) != 0 {
		pretty.Print(file)
		t.Fatal(diff)
	}
}

func TestMissingSlot(t *testing.T) {
	in := `{"Join parts:":["one","two","three"]}`
	want := core.Join{Parts: []rt.TextEval{
		core.T("one"), core.T("two"), core.T("three"),
	}}
	var have core.Join
	var msg map[string]any
	if e := json.Unmarshal([]byte(in), &msg); e != nil {
		t.Fatal(e)
	} else if e := story.DecodeMessage(&have, msg); e != nil {
		pretty.Println("got:", have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println("got:", have)
		t.Fatal(diff)
	}
}
