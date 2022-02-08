package bgen_test

import (
  "testing"

  "git.sr.ht/~ionous/tapestry/dl/core"
  "git.sr.ht/~ionous/tapestry/dl/list"
  "git.sr.ht/~ionous/tapestry/dl/story"
)

// close strings use drop downs, and the field should be their $KEY
func TestStringChoice(t *testing.T) {
  if e := testBlocks(
    &story.TraitPhrase{
      AreEither: story.AreEither{
        Str: story.AreEither_Canbe,
      }}, `{
  "id": "test-1",
  "type": "trait_phrase",
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
    &story.NamedNoun{
      Determiner: story.Determiner{Str: story.Determiner_The},
      Name:       story.NounName{Str: "table"},
    }, `{
  "id": "test-1",
  "type": "named_noun",
  "extraState": {
    "DETERMINER": 1,
    "NAME": 1
  },
  "fields": {
    "DETERMINER": "the",
    "NAME": "table"
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// blocks with optional members should just skip happily to the next member
func TestSkippedSlot(t *testing.T) {
  if e := testBlocks(&list.ListEach{}, `{
  "id": "test-1",
  "type": "list_each",
  "extraState": {
    "DO": 1
  },
  "inputs": {
    "DO": {
      "block": {
        "id": "test-2",
        "type": "activity",
        "extraState": {}
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// empty slots shouldn't get extra closes
// previously this was getting 4 extra closes
func TestEndSlot(t *testing.T) {
  if e := testBlocks(&list.PutEdge{
    From: &core.GetVar{},
  }, `{
  "id": "test-1",
  "type": "put_edge",
  "extraState": {
    "FROM": 1
  },
  "inputs": {
    "FROM": {
      "block": {
        "id": "test-2",
        "type": "get_var",
        "extraState": {
          "NAME": 1
        },
        "fields": {
          "NAME": ""
        }
      }
    }
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
  "id": "test-1",
  "type": "event_block",
  "extraState": {
    "TARGET": 1
  },
  "fields": {
    "TARGET": "$KINDS"
  },
  "inputs": {
    "TARGET": {
      "block": {
        "id": "test-2",
        "type": "plural_kinds",
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
  if e := testBlocks(&story.StoryLines{
    Lines: []story.StoryStatement{
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
  "id": "test-1",
  "type": "story_lines",
  "extraState": {
    "LINES": 2
  },
  "inputs": {
    "LINES": {
      "block": {
        "id": "test-2",
        "type": "_kinds_of_kind_stack",
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
            "id": "test-3",
            "type": "_kinds_of_kind_stack",
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
