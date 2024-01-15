package jsn_test

import (
	_ "embed"
	"testing"
)

// test that the compact encoding matches a particular "golden image"
func TestCompactEncoder(t *testing.T) {
	t.Fatal("does this need porting?")
	// if str, e := cout.Marshal(&debug.FactorialStory, story.CompactEncoder); e != nil {
	// 	t.Fatal(e)
	// } else if str != debug.FactorialJs {
	// 	t.Fatal(str)
	// }
}

// cinStates: xDecoder.readFlow() reads the Text:initially signature
// but previously only used it to separate out the parameter names;
// it didnt validate it against the flow, and handed the json map to Arg_Marshal
// which read any matching fields, ex. the the _ field holding "description".
// now: it tests that the lede of the command matches the first part of the signature
func TestExpectedFailure(t *testing.T) {
	t.Fatal("does this need porting?")
	// var dst assign.Arg_Slice
	// var msg map[string]any
	// if e := json.Unmarshal([]byte(failure), &msg); e != nil {
	// 	t.Fatal(e)
	// } else if e := cin.NewDecoder(cin.Signatures(tapestry.AllSignatures)).
	// 	SetSlotDecoder(core.CompactSlotDecoder).
	// 	Decode(&dst, msg); e == nil {
	// 	t.Fatal("expected error")
	// } else {
	// 	t.Log("ok:", e)
	// }
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
