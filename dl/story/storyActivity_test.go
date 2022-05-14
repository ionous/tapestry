package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry"

	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/kr/pretty"
)

// load and run some actions
func TestStoryActivity(t *testing.T) {
	var prog rt.Execute
	if b, e := json.Marshal(_say_exec); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(rt.Execute_Slot{&prog}, tapestry.Registry(), b); e != nil {
		t.Fatal(e)
	} else {
		var run testRuntime
		out := print.NewLines()
		run.SetWriter(out)
		if e := safe.Run(&run, prog); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out.Lines(), []string{"hello"}); len(diff) > 0 {
			t.Fatal(diff)
		}
	}
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
			"$VALUE": map[string]interface{}{
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
