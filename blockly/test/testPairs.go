package test

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/testdl"
	"git.sr.ht/~ionous/tapestry/jsn"
)

// contents of block files and their matching representation in go.
// ( used by un/block to test read/writing of block files )
var Pairs = []struct {
	Name string
	Test jsn.Marshalee
	Json string
}{{
	// swap member of the flow
	/*test*/ "Swap",
	&testdl.TestFlow{
		Swap: testdl.TestSwap{
			Choice: testdl.TestSwap_C_Opt,
			Value: &testdl.TestTxt{
				Str: "something",
			},
		},
	}, `{
  "type": "test_flow",
  "id": "test-1",
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
        "id": "test-2",
        "fields": {
          "TEST_TXT": "something"
        }
      }
    }
  }
}`}, {
	// repeats of a specific flow
	/*test*/ "Slice",
	&literal.FieldList{
		Fields: []literal.FieldValue{{
			Field: "first",
			Value: &literal.NumValue{
				Value: 5,
			}}, {
			Field: "second",
			Value: &literal.TextValue{
				Value: "five",
			}},
		},
	}, `{
  "type": "field_list",
  "id": "test-1",
  "extraState": {
    "FIELDS": 2
  },
  "inputs": {
    "FIELDS0": {
      "block": {
        "type": "field_value",
        "id": "test-2",
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
              "id": "test-3",
              "extraState": {
                "VALUE": 1
              },
              "fields": {
                "VALUE": 5
              }
            }
          }
        }
      }
    },
    "FIELDS1": {
      "block": {
        "type": "field_value",
        "id": "test-4",
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
              "id": "test-5",
              "extraState": {
                "VALUE": 1
              },
              "fields": {
                "VALUE": "five"
              }
            }
          }
        }
      }
    }
  }
}`}, {
	// ------------------------------------------------------------
	// a flow within a flow
	/*test*/ "Embed",
	&testdl.TestEmbed{
		TestFlow: testdl.TestFlow{},
	}, `{
  "type": "test_embed",
  "id": "test-1",
  "extraState": {
    "TEST_FLOW": 1
  },
  "inputs": {
    "TEST_FLOW": {
      "block": {
        "type": "test_flow",
        "id": "test-2",
        "extraState": {}
      }
    }
  }
}`}, {
	// ------------------------------------------------------------
	// a flow with a single slot.
	/*test*/ "Slot",
	&literal.FieldValue{
		Field: "test",
		Value: &literal.NumValue{
			Value: 5,
		}}, `{
  "type": "field_value",
  "id": "test-1",
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
        "id": "test-2",
        "extraState": {
          "VALUE": 1
        },
        "fields": {
          "VALUE": 5
        }
      }
    }
  }
}`}, {
	// ------------------------------------------------------------
	// repeats of a non-stacking slot.
	/*test*/ "Series",
	&testdl.TestFlow{
		Slots: []testdl.TestSlot{
			&testdl.TestFlow{},
			&testdl.TestFlow{},
		}}, `{
  "type": "test_flow",
  "id": "test-1",
  "extraState": {
    "SLOTS": 2
  },
  "inputs": {
    "SLOTS0": {
      "block": {
        "type": "test_flow",
        "id": "test-2",
        "extraState": {}
      }
    },
    "SLOTS1": {
      "block": {
        "type": "test_flow",
        "id": "test-3",
        "extraState": {}
      }
    }
  }
}`}, {
	// ------------------------------------------------------------
	// test a block with a field:value pair (use some literal text)
	/*test*/ "FieldDefinition",
	&literal.TextValue{Value: "hello world"}, `{
  "type": "text_value",
  "id": "test-1",
  "extraState": {
    "VALUE": 1
  },
  "fields": {
    "VALUE": "hello world"
  }
}`}, {
	// ------------------------------------------------------------
	// an array of primitives is a list of dummy inputs.
	// ( noting that blockly ignores dummies when saving,
	// so they get saved in the "fields" section )
	/*test*/ "List",
	&literal.TextValues{
		Values: []string{"a", "b", "c"},
	}, `{
  "type": "text_values",
  "id": "test-1",
  "extraState": {
    "VALUES": 3
  },
  "fields": {
    "VALUES0": "a",
    "VALUES1": "b",
    "VALUES2": "c"
  }
}`}, {
	// ------------------------------------------------------------
	// nested content-free blocks
	/*test*/ "Stack",
	&story.StoryFile{
		StoryLines: []story.StoryStatement{
			&story.StoryBreak{},
			&story.StoryBreak{},
			&story.StoryBreak{},
		}}, `{
  "type": "story_file",
  "id": "test-1",
  "extraState": {
    "STORY_LINES": 3
  },
  "inputs": {
    "STORY_LINES": {
      "block": {
        "type": "_story_break_stack",
        "id": "test-2",
        "extraState": {},
        "next": {
          "block": {
            "type": "_story_break_stack",
            "id": "test-3",
            "extraState": {},
            "next": {
              "block": {
                "type": "_story_break_stack",
                "id": "test-4",
                "extraState": {}
              }
            }
          }
        }
      }
    }
  }
}`}, {
	// ------------------------------------------------------------
	// repeats of an empty series
	/*test*/ "EmptySeries",
	&testdl.TestFlow{
		Slots: []testdl.TestSlot{},
	}, `{
  "type": "test_flow",
  "id": "test-1",
  "extraState": {}
}`,
}}
