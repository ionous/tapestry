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
  "Swap",
  &testdl.TestFlow{
    Swap: testdl.TestSwap{
      Choice: testdl.TestSwap_C_Opt,
      Value: &testdl.TestTxt{
        Str: "something",
      },
    },
  }, `{
  "id": "test-1",
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
        "id": "test-2",
        "type": "test_txt",
        "fields": {
          "TEST_TXT": "something"
        }
      }
    }
  }
}`}, {
  // repeats of a specific flow
  "Slice",
  &literal.FieldValues{
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
  }, `{
  "id": "test-1",
  "type": "field_values",
  "extraState": {
    "CONTAINS": 2
  },
  "inputs": {
    "CONTAINS0": {
      "block": {
        "id": "test-2",
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
              "id": "test-3",
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
        "id": "test-4",
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
              "id": "test-5",
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
}`}, {
  // ------------------------------------------------------------
  // a flow within a flow
  "Embed",
  &testdl.TestEmbed{
    TestFlow: testdl.TestFlow{},
  }, `{
  "id": "test-1",
  "type": "test_embed",
  "extraState": {
    "TEST_FLOW": 1
  },
  "inputs": {
    "TEST_FLOW": {
      "block": {
        "id": "test-2",
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`}, {
  // ------------------------------------------------------------
  // a flow with a single slot.
  "Slot",
  &literal.FieldValue{
    Field: "test",
    Value: &literal.NumValue{
      Num: 5,
    }}, `{
  "id": "test-1",
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
        "id": "test-2",
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
}`}, {
  // ------------------------------------------------------------
  // repeats of a non-stacking slot.
  "Series",
  &testdl.TestFlow{
    Slots: []testdl.TestSlot{
      &testdl.TestFlow{},
      &testdl.TestFlow{},
    }}, `{
  "id": "test-1",
  "type": "test_flow",
  "extraState": {
    "SLOTS": 2
  },
  "inputs": {
    "SLOTS0": {
      "block": {
        "id": "test-2",
        "type": "test_flow",
        "extraState": {}
      }
    },
    "SLOTS1": {
      "block": {
        "id": "test-3",
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`}, {
  // ------------------------------------------------------------
  // test a block with a field:value pair (use some literal text)
  "Field",
  &literal.TextValue{Text: "hello world"}, `{
  "id": "test-1",
  "type": "text_value",
  "extraState": {
    "TEXT": 1
  },
  "fields": {
    "TEXT": "hello world"
  }
}`}, {
  // ------------------------------------------------------------
  // an array of primitives is a list of dummy inputs.
  // ( noting that blockly ignores dummies when saving,
  // so they get saved in the "fields" section )
  "List",
  &literal.TextValues{
    Values: []string{"a", "b", "c"},
  }, `{
  "id": "test-1",
  "type": "text_values",
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
  // test nested content-free blocks
  "Stack",
  &story.StoryFile{
    StoryLines: []story.StoryStatement{
      &story.StoryBreak{},
      &story.StoryBreak{},
      &story.StoryBreak{},
    }}, `{
  "id": "test-1",
  "type": "story_file",
  "extraState": {
    "STORY_LINES": 3
  },
  "inputs": {
    "STORY_LINES": {
      "block": {
        "id": "test-2",
        "type": "_story_break_stack",
        "extraState": {},
        "next": {
          "block": {
            "id": "test-3",
            "type": "_story_break_stack",
            "extraState": {},
            "next": {
              "block": {
                "id": "test-4",
                "type": "_story_break_stack",
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
  "EmptySeries",
  &testdl.TestFlow{
    Slots: []testdl.TestSlot{},
  }, `{
  "id": "test-1",
  "type": "test_flow",
  "extraState": {}
}`,
}}
