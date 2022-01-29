package bgen

import (
	"bytes"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/kr/pretty"
)

// test writing a block with a field:value pair (use some literal text)
func TestFields(t *testing.T) {
	el := &literal.TextValue{Text: "hello world"}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "text_value",
  "extraState": {
    "TEXT": 1
  },
  "fields": {
    "TEXT": "hello world"
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// test writing some blockly nested next blocks
func TestStack(t *testing.T) {
	el := &story.StoryLines{Lines: []story.StoryStatement{
		&story.StoryBreak{},
		&story.StoryBreak{},
		&story.StoryBreak{},
	}}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "story_lines",
  "extraState": {
    "LINES": 3
  },
  "next": {
    "block": {
      "type": "story_break",
      "extraState": {},
      "next": {
        "block": {
          "type": "story_break",
          "extraState": {},
          "next": {
            "block": {
              "type": "story_break",
              "extraState": {}
            }
          }
        }
      }
    }
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// a primitive list is a list of dummy inputs
// noting that blockly ignores dummies when saving --
// so get saved in the "fields" section
func TestSeries(t *testing.T) {
	el := &literal.TextValues{
		Values: []string{"a", "b", "c"},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "text_values",
  "extraState": {
    "VALUES": 3
  },
  "fields": {
    "VALUES0": "a",
    "VALUES1": "b",
    "VALUES2": "c"
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// test a slot member of the flow
func TestSlot(t *testing.T) {
	el := &literal.FieldValue{
		Field: "test",
		Value: &literal.NumValue{
			Num: 5,
		}}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "field_value",
  "extraState": {
    "FIELD": 1,
    "VALUE": 1
  },
  "fields": {
    "FIELD": "test"
  },
  "inputs": {
    "VALUE": {
      "block": {
        "type": "num_value",
        "extraState": {
          "NUM": 1
        },
        "fields": {
          "NUM": 5
        }
      }
    }
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// func TestEmbed(t *testing.T) {
// 	el := &literal.FieldValues{
// 		Contains: []literal.FieldValue{
// 			Field: "test",
// 			Value: &literal.NumValue{
// 				Num: 5,
// 			}}}
// 	var out js.Builder
// 	enc := chart.MakeEncoder()
// 	if e := enc.Marshal(el, newTopBlock(&enc, &out)); e != nil {
// 		t.Fatal(e)
// 	} else if str, e := indent(out.String()); e != nil {
// 		t.Fatal(e, str)
// 	} else if diff := pretty.Diff(str, ``); len(diff) > 0 {
// 		t.Log(str)
// 		t.Fatal("ng", diff)
// 	}
// }

func indent(str string) (ret string, err error) {
	var indent bytes.Buffer
	if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
		ret = str
	} else {
		ret = indent.String()
	}
	return
}

// // test a flow containing flows
// func TestEmbed(t *testing.T) {
// 		el := &literal.FieldValues{
// 			Contains: []literal.FieldValue{
// 		Field       string       `if:"label=field,type=text"`
// 	Value       LiteralValue `if:"label=value"`
// }}
// 	}}

// }
//    "field_values": {
//                 "uses": "flow",
//                 "spec": "fields {_%contains+field_value}",
//                 "group": "internal",
//                 "slot": ["literal_value"],
//                 "desc": [
//                     "A series of values all for the same record.",
//                     "While it can be specified wherever a literal value can, it only has meaning when the record type is known."
//                 ]
//             }
