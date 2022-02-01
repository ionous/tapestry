package bgen_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/bgen"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/testdl"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/kr/pretty"
)

// test writing a block with a field:value pair (use some literal text)
func TestFields(t *testing.T) {
	el := &literal.TextValue{Text: "hello world"}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
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

// test a flow within a flow
func TestEmbeds(t *testing.T) {
	el := &testdl.TestEmbed{
		TestFlow: testdl.TestFlow{},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "test_embed",
  "extraState": {
    "TEST_FLOW": 1
  },
  "inputs": {
    "TEST_FLOW": {
      "block": {
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// test a swap member of the flow
func TestSwap(t *testing.T) {
	el := &testdl.TestFlow{
		Swap: testdl.TestSwap{
			Choice: testdl.TestSwap_C_Opt,
			Value: &testdl.TestTxt{
				Str: "something",
			},
		},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "test_flow",
  "extraState": {
    "SWAP": 1
  },
  "fields": {
    "SWAP": "$C"
  },
  "inputs": {
    "SWAP": {
      "block": {
        "type": "test_txt",
        "extraState": {
          "TEST_TXT": 1
        },
        "fields": {
          "TEST_TXT": "something"
        }
      }
    }
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
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
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

// test writing some blockly nested next blocks
func TestStack(t *testing.T) {
	el := &story.StoryLines{Lines: []story.StoryStatement{
		&story.StoryBreak{},
		&story.StoryBreak{},
		&story.StoryBreak{},
	}}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
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
func TestList(t *testing.T) {
	el := &literal.TextValues{
		Values: []string{"a", "b", "c"},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
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

// repeats of a specific flow
func TestSlice(t *testing.T) {
	el := &literal.FieldValues{
		Contains: []literal.FieldValue{{
			Field: "first",
			Value: &literal.NumValue{
				Num: 5,
			}}, {
			Field: "second",
			Value: &literal.TextValue{
				Text: "five",
			}},
		},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "field_values",
  "extraState": {
    "CONTAINS": 2
  },
  "inputs": {
    "CONTAINS0": {
      "block": {
        "type": "field_value",
        "extraState": {
          "FIELD": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD": "first"
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
      }
    },
    "CONTAINS1": {
      "block": {
        "type": "field_value",
        "extraState": {
          "FIELD": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD": "second"
        },
        "inputs": {
          "VALUE": {
            "block": {
              "type": "text_value",
              "extraState": {
                "TEXT": 1
              },
              "fields": {
                "TEXT": "five"
              }
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

// repeats of a non-stacking slot.
func TestSeries(t *testing.T) {
	// FIX FIX FIX -- the series doesnt generate remotely correct.
	el := &testdl.TestFlow{
		Slots: []testdl.TestSlot{
			&testdl.TestFlow{},
			&testdl.TestFlow{},
		},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "test_flow",
  "extraState": {
    "SLOTS": 2
  },
  "inputs": {
    "SLOTS0": {
      "block": {
        "type": "test_flow",
        "extraState": {}
      }
    },
    "SLOTS1": {
      "block": {
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// repeats of an empty series
func TestEmptySeries(t *testing.T) {
	el := &testdl.TestFlow{
		Slots: []testdl.TestSlot{},
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if e := enc.Marshal(el, bgen.NewTopBlock(&enc, &out)); e != nil {
		t.Fatal(e)
	} else if str, e := indent(out.String()); e != nil {
		t.Fatal(e, str)
	} else if diff := pretty.Diff(str, `{
  "type": "test_flow",
  "extraState": {}
}`); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// other ideas: might verify an empty repeat
// a flow member....

func indent(str string) (ret string, err error) {
	var indent bytes.Buffer
	if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
		ret = str
	} else {
		ret = indent.String()
	}
	return
}
