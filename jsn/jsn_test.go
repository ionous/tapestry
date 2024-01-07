package jsn_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"

	"git.sr.ht/~ionous/tapestry/jsn/dout"
)

// test that the detailed format can be used to write out, and read back in the same data
// we dont much care what it looks like anymore.
func TestDetailsEncodeDecode(t *testing.T) {
	if d, e := dout.Encode(&debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else if in, e := story.DetailedDecode(b); e != nil {
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(debug.FactorialStory, in); len(diff) != 0 {
			pretty.Print(in)
			t.Fatal(diff)
		}
	}
}

// test that the compact encoding matches a particular "golden image"
func TestCompactEncoder(t *testing.T) {
	if str, e := cout.Marshal(&debug.FactorialStory, story.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if str != debug.FactorialJs {
		t.Fatal(str)
	}
}

// test the compact decoder can read from the "golden image" and get the hardwired factorial story.
func TestCompactDecode(t *testing.T) {
	var msg map[string]any
	if e := json.Unmarshal([]byte(debug.FactorialJs), &msg); e != nil {
		t.Fatal(e)
	} else if file, e := story.CompactDecode(msg); e != nil {
		pretty.Println(file)
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(debug.FactorialStory, file); len(diff) != 0 {
			pretty.Print(file)
			t.Fatal(diff)
		}
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
	} else if e := story.Decode(&have, msg, story.AllSignatures); e != nil {
		pretty.Println("got:", have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println("got:", have)
		t.Fatal(diff)
	}
}

// cinStates: xDecoder.readFlow() reads the Text:initially signature
// but previously only used it to separate out the parameter names;
// it didnt validate it against the flow, and handed the json map to Arg_Marshal
// which read any matching fields, ex. the the _ field holding "description".
// now: it tests that the lede of the command matches the first part of the signature
func TestExpectedFailure(t *testing.T) {
	var dst assign.Arg_Slice
	var msg map[string]any
	if e := json.Unmarshal([]byte(failure), &msg); e != nil {
		t.Fatal(e)
	} else if e := cin.NewDecoder(cin.Signatures(tapestry.AllSignatures)).
		SetSlotDecoder(core.CompactSlotDecoder).
		Decode(&dst, msg); e == nil {
		t.Fatal("expected error")
	} else {
		t.Log("ok:", e)
	}
}

var failure = `{
  "Text:initially:": [
    "description",
    {
      "Object:field:": [
        "@obj",
        "description"
      ]
    }
  ]
}`
