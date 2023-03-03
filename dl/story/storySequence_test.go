package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/support/asm"

	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// test that importing cycling text transforms to the proper runtime command
func TestImportSequence(t *testing.T) {
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els), storyMarshaller)

	var cmd story.CycleText
	if b, e := json.Marshal(_cycle_text); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&cmd, tapestry.Registry(), b); e != nil {
		t.Fatal(e)
	} else {
		p := assign.FromText{Value: &cmd} // wrap the cycle text in a slot since that's the level PreImport operates on
		if asm.ImportStory(k, t.Name(), &p); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(p.Value, &_import_target); len(diff) > 0 {
			t.Fatal(pretty.Print("want:", _import_target, "have:", p.Value))
		}
	}
}

// the cycle text data should look like this after import
var _import_target = core.CallCycle{
	Name: "seq_1",
	Parts: []rt.TextEval{
		T("a"),
		T("b"),
		T("c"),
	},
}

var _cycle_text = map[string]interface{}{
	"type": "cycle_text",
	"value": map[string]interface{}{
		"$PARTS": []interface{}{
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$VALUE": map[string]interface{}{
							"type":  "text",
							"value": "a",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$VALUE": map[string]interface{}{
							"type":  "text",
							"value": "b",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$VALUE": map[string]interface{}{
							"type":  "text",
							"value": "c",
						}}}}}},
}
