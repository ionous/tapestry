package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"

	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestPatternActivity(t *testing.T) {
	var prog core.Activity
	if b, e := json.Marshal(_pattern_activity); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&prog, iffy.Registry(), b); e != nil {
		t.Fatal(e)
	} else {
		var run testRuntime
		out := print.NewLines()
		run.SetWriter(out)
		// should this call/test buildRule
		if e := safe.Run(&run, &prog); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out.Lines(), []string{"hello", "hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
}

func TestPatternActions(t *testing.T) {
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)

	var prog story.PatternActions
	if b, e := json.Marshal(_pattern_actions); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&prog, iffy.Registry(), b); e != nil {
		t.Fatal(e)
	} else if e := prog.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		t.Fatal("implemented")
		// var buf strings.Builder
		// // -- eph?rule
		// // rule
		// tables.WriteCsv(db, &buf, "select progType from eph_prog", 1)
		// // example, pattern_name
		// tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		// // 1 pattern handler reference
		// tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		// // 1, 1 - the first name, the first program are used to make the rule
		// tables.WriteCsv(db, &buf, "select idNamedPattern, idProg from eph_rule", 2)
		// if have, want := buf.String(), lines(
		// 	"rule",
		// 	"example,pattern",
		// 	"0", // eph_pattern -- rules are recorded via eph_prog,
		// 	"2,1",
		// ); have != want {
		// 	t.Fatal(have)
		// }
	}
}

var _pattern_actions = map[string]interface{}{
	"type": "pattern_actions",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "example",
		},
		"$PATTERN_RULES": map[string]interface{}{
			"type": "pattern_rules",
			"value": map[string]interface{}{
				"$PATTERN_RULE": []interface{}{
					map[string]interface{}{
						"type": "pattern_rule",
						"value": map[string]interface{}{
							"$GUARD": map[string]interface{}{
								"type": "bool_eval",
								"value": map[string]interface{}{
									"type":  "always",
									"value": map[string]interface{}{},
								}},
							"$HOOK": map[string]interface{}{
								"type": "program_hook",
								"value": map[string]interface{}{
									"$ACTIVITY": _pattern_activity,
								}}}}}}}},
}

var _pattern_activity = map[string]interface{}{
	"type": "activity",
	"value": map[string]interface{}{
		"$EXE": []interface{}{
			_say_exec,
			_say_exec,
		},
	},
}

var _say_exec = map[string]interface{}{
	"type": "execute",
	"value": map[string]interface{}{
		"type": "say_text",
		"value": map[string]interface{}{
			"$TEXT": _text_eval,
		},
	},
}

var _text_eval = map[string]interface{}{
	"type": "text_eval",
	"value": map[string]interface{}{
		"type": "text_value",
		"value": map[string]interface{}{
			"$TEXT": map[string]interface{}{
				"type":  "text",
				"value": "hello",
			},
		},
	},
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type testRuntime struct {
	baseRuntime
	writer.Sink
}
