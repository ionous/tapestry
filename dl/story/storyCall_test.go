package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/story"

	"git.sr.ht/~ionous/iffy/jsn/din"
	"github.com/kr/pretty"
)

// test assembling a pattern call.
// not a huge point to this test: just verifies that it generates a pattern reference.
func TestDetermineNum(t *testing.T) {
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	//
	var rule story.Determine
	if b, e := json.Marshal(factorialDetermine); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&rule, iffy.Registry(), b); e != nil {
		t.Fatal(e)
	} else if ptr, e := rule.ImportStub(k); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(ptr, &core.CallPattern{
		Pattern: core.PatternName{Str: "factorial"},
		Arguments: core.CallArgs{
			Args: []core.CallArg{{
				Name: "num",
				From: &core.FromNum{
					&literal.NumValue{3},
				}}}}}); len(diff) > 0 {
		t.Fatal(diff)
	} else if diff := pretty.Diff(els, []eph.Ephemera{
		&eph.EphRefs{
			Kinds: "factorial",
			From:  "pattern",
			ReferTo: []eph.EphParams{{
				Affinity: eph.Affinity{eph.Affinity_Number},
				Name:     "num",
			},
			},
		},
	}); len(diff) > 0 {
		t.Fatal(diff)
	}
}

// determine num of factorial where num = 3
var factorialDetermine = map[string]interface{}{
	"type": "determine",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$ARGUMENTS": map[string]interface{}{
			"type": "arguments",
			"value": map[string]interface{}{
				"$ARGS": []interface{}{
					map[string]interface{}{
						"type": "argument",
						"value": map[string]interface{}{
							"$FROM": map[string]interface{}{
								"type": "assignment",
								"value": map[string]interface{}{
									"type": "from_num",
									"value": map[string]interface{}{
										"$VAL": map[string]interface{}{
											"type": "number_eval",
											"value": map[string]interface{}{
												"type": "num_value",
												"value": map[string]interface{}{
													"$NUM": map[string]interface{}{
														"type":  "number",
														"value": 3.0,
													}}}}}}},
							"$NAME": map[string]interface{}{
								"type":  "variable_name",
								"value": "num",
							}}}}}}},
}
