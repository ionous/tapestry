package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"

	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

// load and run some actions
func TestStoryActivity(t *testing.T) {
	var prog core.Activity
	if b, e := json.Marshal(_pattern_activity); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&prog, iffy.Registry(), b); e != nil {
		t.Fatal(e)
	} else {
		var run testRuntime
		out := print.NewLines()
		run.SetWriter(out)
		if e := safe.Run(&run, &prog); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out.Lines(), []string{"hello", "hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
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
