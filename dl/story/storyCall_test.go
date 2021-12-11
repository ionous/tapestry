package story_test

import (
	"encoding/json"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/story"

	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
	"github.com/kr/pretty"
)

// test calling a pattern
// note: the pattern is undefined.
func TestDetermineNum(t *testing.T) {
	expect := core.CallPattern{
		Pattern: value.PatternName{Str: "factorial"},
		Arguments: core.CallArgs{
			Args: []core.CallArg{{
				Name: "num",
				From: &core.FromNum{
					&literal.NumValue{3},
				}}}}}
	k, db := newImporter(t, testdb.Memory)
	defer db.Close()
	//
	var rule story.Determine
	if b, e := json.Marshal(factorialDetermine); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&rule, iffy.Registry(), b); e != nil {
		t.Fatal(e)
	} else if ptr, e := rule.ImportStub(k); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(ptr, &expect); len(diff) != 0 {
		t.Fatal(diff)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
		tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		if have, want := buf.String(), lines(
			// eph_prog count
			// no programs b/c no container for the call into determine.
			"0",
			// eph_rule count
			// no rules b/c the pattern is called but not implemented.
			"0",
			// eph_pattern
			"2,3,4,-1", // from NewPatternRef -> "determine num" takes a parameter that is from a number eval
			"2,2,5,-1", // from NewPatternRef -> "determine num" indicates factorial returns a number eval
			//
			"factorial,pattern", // 1.
			"num,argument",      // 2.
			"number_eval,type",  // 3.
			"patterns,type",     // 4.
		); have != want {
			t.Fatal(have)
		}
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
