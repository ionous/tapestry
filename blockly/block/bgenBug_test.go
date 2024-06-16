package block_test

import (
	"testing"

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
	if e := testBlocks(&list.ListRepeat{}, `{
  "type": "list_repeat",
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
// -- EventTarget no longer exists.
// func TestExcessState(t *testing.T) {
// 	if e := testBlocks(&story.EventBlock{
// 		Target: story.EventTarget{
// 			Choice: story.EventTarget_Kinds_Opt,
// 			Value: &story.PluralKinds{
// 				Str: "x",
// 			},
// 		},
// 	}, `{
//   "type": "event_block",
//   "id": "test-1",
//   "extraState": {
//     "TARGET": 1
//   },
//   "fields": {
//     "TARGET": "$KINDS"
//   },
//   "inputs": {
//     "TARGET": {
//       "block": {
//         "type": "plural_kinds",
//         "id": "test-2",
//         "fields": {
//           "PLURAL_KINDS": "x"
//         }
//       }
//     }
//   }
// }`); e != nil {
// 		t.Fatal(e)
// 	}
// }

// story lines should be a block with no output, and one stacking input
// the stacks should all use the "_define_kind_stack" type
func TestStoryLines(t *testing.T) {
	if e := testBlocks(&story.StoryFile{
		Statements: []story.StoryStatement{
			&story.DefineKind{
				KindName:         literal.T("cats"),
				AncestorKindName: literal.T("animal"),
			},
			&story.DefineKind{
				KindName:         literal.T("cats"),
				AncestorKindName: literal.T("animal"),
			},
		},
	}, `{
  "type": "story_file",
  "id": "test-1",
  "extraState": {
    "STATEMENTS": 2
  },
  "inputs": {
    "STATEMENTS": {
      "block": {
        "type": "_define_kind_stack",
        "id": "test-2",
        "extraState": {
          "KIND_NAME": 1,
          "ANCESTOR_KIND_NAME": 1
        },
        "inputs": {
          "KIND_NAME": {
            "block": {
              "type": "text_value",
              "id": "test-3",
              "extraState": {
                "VALUE": 1
              },
              "fields": {
                "VALUE": "cats"
              }
            }
          },
          "ANCESTOR_KIND_NAME": {
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
            "type": "_define_kind_stack",
            "id": "test-5",
            "extraState": {
              "KIND_NAME": 1,
              "ANCESTOR_KIND_NAME": 1
            },
            "inputs": {
              "KIND_NAME": {
                "block": {
                  "type": "text_value",
                  "id": "test-6",
                  "extraState": {
                    "VALUE": 1
                  },
                  "fields": {
                    "VALUE": "cats"
                  }
                }
              },
              "ANCESTOR_KIND_NAME": {
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
