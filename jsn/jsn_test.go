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
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"
)

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
	} else if file, e := story.Decode(msg); e != nil {
		pretty.Println(file)
		t.Fatal(e)
	} else {
		if diff := pretty.Diff(debug.FactorialStory, file); len(diff) != 0 {
			pretty.Print(file)
			t.Fatal(diff)
		}
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
