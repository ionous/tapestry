package jsn_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// previously being unable to parse the contents of a pattern call would fail to throw an error
func TestBug_PatternOfPattern(t *testing.T) {
	var dst rt.Execute
	var msg map[string]any
	if e := json.Unmarshal([]byte(`{"PrintVantage where:":{"ParentOf obj:":"@actor"}}`), &msg); e != nil {
		t.Fatal(e)
	} else if e := story.Decode(rt.Execute_Slot{&dst}, msg, nil); e == nil {
		pretty.Println(dst)
		t.Fatal("expected failure")
	}
}
