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
