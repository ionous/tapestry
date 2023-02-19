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

// closed strings use drop downs, and the field should be their $KEY
func TestStringChoice(t *testing.T) {
  if e := testBlocks(
    &story.TraitPhrase{
      AreEither: story.AreEither{
        Str: story.AreEither_Canbe,
      }}, `{
  "type": "trait_phrase",
  "id": "test-1",
  "extraState": {
    "ARE_EITHER": 1
  },
  "fields": {
    "ARE_EITHER": "$CANBE"
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// until there are variables ( or something ) for the hints
// their text should hold normal text not $KEY values
func TestStringHints(t *testing.T) {
  if e := testBlocks(
    &story.CommonNoun{
      Determiner: story.Determiner{Str: story.Determiner_The},
      Noun:       story.NounNamed{Name: story.NounName{Str: "table"}},
    }, `{
  "type": "common_noun",
  "id": "test-1",
  "extraState": {
    "DETERMINER": 1,
    "NOUN": 1
  },
  "fields": {
    "DETERMINER": "the"
  },
  "inputs": {
    "NOUN": {
      "block": {
        "type": "noun_named",
        "id": "test-2",
        "extraState": {
          "NAME": 1
        },
        "fields": {
          "NAME": "table"
        }
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// blocks with optional members should just skip happily to the next member
func TestSkippedSlot(t *testing.T) {
  if e := testBlocks(&list.ListEach{}, `{
  "type": "list_each",
  "id": "test-1",
  "extraState": {}
}`); e != nil {
    t.Fatal(e)
  }
}

// empty slots shouldn't get extra closes
// previously this was getting 4 extra closes
// func TestEndSlot(t *testing.T) {
//   if e := testBlocks(&list.ListPush{
//     From: &core.GetVar{},
//   }, `{
//   "type": "put_edge",
//   "id": "test-1",
//   "extraState": {
//     "FROM": 1
//   },
//   "inputs": {
//     "FROM": {
//       "block": {
//         "type": "get_var",
//         "id": "test-2",
//         "extraState": {
//           "NAME": 1
//         },
//         "fields": {
//           "NAME": ""
//         }
//       }
//     }
//   }
// }`); e != nil {
//     t.Fatal(e)
//   }
// }

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
      &story.KindsOfKind{
        PluralKinds:  story.PluralKinds{"cats"},
        SingularKind: story.SingularKind{"animal"},
      },
      &story.KindsOfKind{
        PluralKinds:  story.PluralKinds{"cats"},
        SingularKind: story.SingularKind{"animal"},
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
        "type": "_kinds_of_kind_stack",
        "id": "test-2",
        "extraState": {
          "PLURAL_KINDS": 1,
          "SINGULAR_KIND": 1
        },
        "fields": {
          "PLURAL_KINDS": "cats",
          "SINGULAR_KIND": "animal"
        },
        "next": {
          "block": {
            "type": "_kinds_of_kind_stack",
            "id": "test-3",
            "extraState": {
              "PLURAL_KINDS": 1,
              "SINGULAR_KIND": 1
            },
            "fields": {
              "PLURAL_KINDS": "cats",
              "SINGULAR_KIND": "animal"
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
