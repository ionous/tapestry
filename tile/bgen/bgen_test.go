package bgen

import (
	"bytes"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// test writing a block with a field:value pair (use some literal text)
func TestFields(t *testing.T) {
	el := &literal.TextValue{Text: "hello world"}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else {
		var indent bytes.Buffer
		if e := json.Indent(&indent, []byte(out.String()), "", "  "); e != nil {
			t.Log(out.String())
			t.Fatal(e)
		} else {
			t.Log("ok", indent.String())
		}
	}
}

// test writing some blockly nested next blocks
func TestStack(t *testing.T) {
	// errutil.Panic = true
	el := &story.StoryLines{Lines: []story.StoryStatement{
		&story.StoryBreak{},
		&story.StoryBreak{},
		&story.StoryBreak{},
	}}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else {
		var indent bytes.Buffer
		if e := json.Indent(&indent, []byte(out.String()), "", "  "); e != nil {
			t.Fatal(e)
		} else {
			t.Log("ok", indent.String())
		}
	}
}
