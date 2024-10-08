package test

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/testdl"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// contents of block files and their matching representation in go.
// ( used by un/block to test read/writing of block files )
var Pairs = []struct {
	Name string
	Test typeinfo.Instance
	Json string
}{{
	// repeats of a specific flow
	/*test*/ "Slice",
	&literal.FieldList{
		Fields: []literal.FieldValue{{
			FieldName: "first",
			Value: &literal.NumValue{
				Value: 5,
			}}, {
			FieldName: "second",
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
          "FIELD_NAME": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD_NAME": "first"
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
          "FIELD_NAME": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD_NAME": "second"
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
		FieldName: "test",
		Value: &literal.NumValue{
			Value: 5,
		}}, `{
  "type": "field_value",
  "id": "test-1",
  "extraState": {
    "FIELD_NAME": 1,
    "VALUE": 1
  },
  "fields": {
    "FIELD_NAME": "test"
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
	&literal.TextList{
		Values: []string{"a", "b", "c"},
	}, `{
  "type": "text_list",
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
	// fix? these used to be empty story breaks;
	// now they are notes -- which can have text
	// so .. these aren't "content free" anymore.
	/*test*/ "Stack",
	&story.StoryFile{
		Statements: []story.StoryStatement{
			&story.StoryNote{},
			&story.StoryNote{},
			&story.StoryNote{},
		}}, `{
  "type": "story_file",
  "id": "test-1",
  "extraState": {
    "STATEMENTS": 3
  },
  "inputs": {
    "STATEMENTS": {
      "block": {
        "type": "_story_note_stack",
        "id": "test-2",
        "extraState": {
          "TEXT": 1
        },
        "fields": {
          "TEXT": ""
        },
        "next": {
          "block": {
            "type": "_story_note_stack",
            "id": "test-3",
            "extraState": {
              "TEXT": 1
            },
            "fields": {
              "TEXT": ""
            },
            "next": {
              "block": {
                "type": "_story_note_stack",
                "id": "test-4",
                "extraState": {
                  "TEXT": 1
                },
                "fields": {
                  "TEXT": ""
                }
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
