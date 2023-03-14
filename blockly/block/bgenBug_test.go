package block_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
)

// bool values should ( for now ) be $KEY
func TestBoolChoice(t *testing.T) {
	if e := testBlocks(
		&literal.BoolValue{
			Value: true,
		}, `{
  "type": "bool_value",
  "id": "test-1",
  "extraState": {
    "VALUE": 1
  },
  "fields": {
    "VALUE": "$TRUE"
  }
}`); e != nil {
		t.Fatal(e)
	}
}

// blocks with optional members should just skip happily to the next member
// fix? empty strings render extraState -- but they probably dont need to.
func TestSkippedSlot(t *testing.T) {
	if e := testBlocks(&list.ListEach{}, `{
  "type": "list_each",
  "id": "test-1",
  "extraState": {
    "AS": 1
  },
  "fields": {
    "AS": ""
  }
}`); e != nil {
		t.Fatal(e)
	}
}

// blocks without mutations shouldnt get extra data
// ( or blockly exceptions and gets very unhappy )
func TestExcessState(t *testing.T) {
	if e := testBlocks(&story.EventBlock{
		Target: story.EventTarget{
			Choice: story.EventTarget_Kinds_Opt,
			Value: &story.PluralKinds{
				Str: "x",
			},
		},
	}, `{
  "type": "event_block",
  "id": "test-1",
  "extraState": {
    "TARGET": 1
  },
  "fields": {
    "TARGET": "$KINDS"
  },
  "inputs": {
    "TARGET": {
      "block": {
        "type": "plural_kinds",
        "id": "test-2",
        "fields": {
          "PLURAL_KINDS": "x"
        }
      }
    }
  }
}`); e != nil {
		t.Fatal(e)
	}
}

// story lines should be a block with no output, and one stacking input
// the stacks should all use the "stacked_kinds_of_kind" type
func TestStoryLines(t *testing.T) {
	if e := testBlocks(&story.StoryFile{
		StoryLines: []story.StoryStatement{
			&story.DefineKinds{
				Kinds:    assign.Ts("cats"),
				Ancestor: assign.T("animal"),
			},
			&story.DefineKinds{
				Kinds:    assign.Ts("cats"),
				Ancestor: assign.T("animal"),
			},
		},
	}, `{
  "type": "story_file",
  "id": "test-1",
  "extraState": {
    "STORY_LINES": 2
  },
  "inputs": {
    "STORY_LINES": {
      "block": {
        "type": "_define_kinds_stack",
        "id": "test-2",
        "extraState": {
          "KINDS": 1,
          "ANCESTOR": 1
        },
        "inputs": {
          "KINDS": {
            "block": {
              "type": "text_values",
              "id": "test-3",
              "extraState": {
                "VALUES": 1
              },
              "fields": {
                "VALUES0": "cats"
              }
            }
          },
          "ANCESTOR": {
            "block": {
              "type": "text_value",
              "id": "test-4",
              "extraState": {
                "VALUE": 1
              },
              "fields": {
                "VALUE": "animal"
              }
            }
          }
        },
        "next": {
          "block": {
            "type": "_define_kinds_stack",
            "id": "test-5",
            "extraState": {
              "KINDS": 1,
              "ANCESTOR": 1
            },
            "inputs": {
              "KINDS": {
                "block": {
                  "type": "text_values",
                  "id": "test-6",
                  "extraState": {
                    "VALUES": 1
                  },
                  "fields": {
                    "VALUES0": "cats"
                  }
                }
              },
              "ANCESTOR": {
                "block": {
                  "type": "text_value",
                  "id": "test-7",
                  "extraState": {
                    "VALUE": 1
                  },
                  "fields": {
                    "VALUE": "animal"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`); e != nil {
		t.Fatal(e)
	}
}
