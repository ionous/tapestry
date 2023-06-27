package jsn_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"

	"git.sr.ht/~ionous/tapestry/jsn/dout"
)

// test that the detailed format can be used to write out, and read back in the same data
// we dont much care what it looks like anymore.
func TestDetailsEncodeDecode(t *testing.T) {
	if d, e := dout.Encode(debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else if in, e := story.DetailedDecode(b); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(debug.FactorialStory, &in); len(diff) != 0 {
			pretty.Print(in)
			t.Fatal(diff)
		}
	}
}

// test that the compact encoding matches a particular "golden image"
func TestCompactEncoder(t *testing.T) {
	if str, e := cout.Marshal(debug.FactorialStory, story.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if str != jsnTestIf {
		t.Fatal(str)
	}
}

// test the compact decoder can read from the "golden image" and get the hardwired factorial story.
func TestCompactDecode(t *testing.T) {
	errutil.Panic = true
	if file, e := story.CompactDecode([]byte(jsnTestIf)); e != nil {
		pretty.Println(file)
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(debug.FactorialStory, &file); len(diff) != 0 {
			pretty.Print(file)
			t.Fatal(diff)
		}
	}
}

//go:embed jsnTest.if
var jsnTestIf string

// TestAnonymousSwap - unit test for broken parsing case
// -- EventTarget no longer exists
// func TestAnonymousSwap(t *testing.T) {
// 	var jsnTestIf = `{"Listen kinds:handlers:": ["things",[]]}`
// 	want := story.EventBlock{
// 		Target: story.EventTarget{
// 			Value: &story.PluralKinds{
// 				Str: "things",
// 			},
// 			Choice: story.EventTarget_Kinds_Opt,
// 		},
// 		Handlers: make([]story.EventHandler, 0, 0),
// 	}
// 	//
// 	var have story.EventBlock
// 	if e := story.Decode(&have, []byte(jsnTestIf), story.AllSignatures); e != nil {
// 		pretty.Println(have)
// 		t.Fatal(e)
// 	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
// 		pretty.Println(have)
// 		t.Fatal(diff)
// 	}
// }

func TestMissingSlot(t *testing.T) {
	in := `{"Join parts:":["one","two","three"]}`
	want := core.Join{Parts: []rt.TextEval{
		core.T("one"), core.T("two"), core.T("three"),
	}}
	var have core.Join
	if e := story.Decode(&have, []byte(in), story.AllSignatures); e != nil {
		pretty.Println("got:", have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println("got:", have)
		t.Fatal(diff)
	}
}
